package resolver

import (
	"context"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/jackc/pgx/v5/pgxpool"
)

type todo struct {
	ID    graphql.ID
	Title string
}

type todoInput struct {
	Title string
}

type Resolver struct {
	DB *pgxpool.Pool
	*TodoResolver
}

func (r *Resolver) Query() *Resolver {
	return &Resolver{
		DB: r.DB,
	}
}

func (r *Resolver) Mutation() *Resolver {
	return &Resolver{
		DB: r.DB,
	}
}

type TodoResolver struct {
	t *todo
}

func (r *TodoResolver) ID() graphql.ID { return r.t.ID }
func (r *TodoResolver) Title() string  { return r.t.Title }

func (r *Resolver) Todo(args struct{ ID graphql.ID }) *TodoResolver {
	const stmt = `SELECT title FROM todo WHERE id = $1`

	row := r.DB.QueryRow(context.Background(), stmt, args.ID)

	t := todo{}
	if err := row.Scan(&t.Title); err != nil {
		return nil
	}

	return &TodoResolver{t: &t}
}

func (r *Resolver) Todos() []*TodoResolver {
	const stmt = `SELECT id, title FROM todo`

	rows, err := r.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil
	}

	var todos []*TodoResolver
	for rows.Next() {
		t := todo{}
		if err := rows.Scan(&t.ID, &t.Title); err != nil {
			return nil
		}

		todos = append(todos, &TodoResolver{t: &t})
	}

	return todos
}

func (r *Resolver) CreateTodo(args struct{ Todo todoInput }) *TodoResolver {
	const stmt = `INSERT INTO todo (title)
		VALUES ($1)
		RETURNING id`

	row := r.DB.QueryRow(context.Background(), stmt, args.Todo.Title)

	var id string
	if err := row.Scan(&id); err != nil {
		return nil
	}

	return &TodoResolver{
		t: &todo{
			ID:    graphql.ID(id),
			Title: args.Todo.Title,
		},
	}
}

func (r *Resolver) DeleteTodo(args struct{ ID graphql.ID }) *int32 {
	const stmt = `DELETE FROM todo WHERE id = $1`

	result, err := r.DB.Exec(context.Background(), stmt, args.ID)
	if err != nil {
		return nil
	}

	return toPtr(int32(result.RowsAffected()))
}
