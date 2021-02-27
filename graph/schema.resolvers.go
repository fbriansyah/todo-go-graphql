package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fbriansyah/todo-go-graphql/graph/generated"
	"github.com/fbriansyah/todo-go-graphql/graph/model"
	"github.com/fbriansyah/todo-go-graphql/internal/auth"
	"github.com/fbriansyah/todo-go-graphql/internal/pkg/jwt"
	"github.com/fbriansyah/todo-go-graphql/internal/todos"
	"github.com/fbriansyah/todo-go-graphql/internal/users"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Todo{}, fmt.Errorf("access denied")
	}
	var todo todos.Todo
	todo.Text = input.Text
	todo.Done = false
	todoID := todo.Save()
	grahpqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}

	return &model.Todo{ID: strconv.FormatInt(todoID, 10), Text: todo.Text, Done: todo.Done, User: grahpqlUser}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		// 1
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil

}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	var resultTodos []*model.Todo
	var dbTodos []todos.Todo
	dbTodos = todos.GetAll()
	for _, todo := range dbTodos {
		grahpqlUser := &model.User{
			ID:   todo.User.ID,
			Name: todo.User.Username,
		}

		resultTodos = append(resultTodos, &model.Todo{
			ID:   todo.ID,
			Text: todo.Text,
			Done: todo.Done,
			User: grahpqlUser,
		})
	}
	return resultTodos, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
