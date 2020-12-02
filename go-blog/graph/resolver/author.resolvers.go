package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"blog/graph/auth"
	"blog/graph/generated"
	"blog/graph/model"
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (r *mutationResolver) Signup(ctx context.Context, input model.AuthorInput) (*bool, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	status := true

	if err != nil {
		status = false
	}

	author := &model.AuthorInput{
		Username: input.Username,
		Email:    input.Email,
		Name:     input.Name,
		Password: string(password),
	}

	err = r.DB.Table("authors").Create(author).Error

	if err != nil {
		status = false
	}
	return &status, err
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.AuthPayload, error) {
	author := &model.Author{}

	r.DB.Table("authors").Where("email = ?", email).First(author)

	if author.Email != "" {
		err := bcrypt.CompareHashAndPassword([]byte(author.Password), []byte(password))
		if err != nil {
			return &model.AuthPayload{}, err
		}

		payload, err := auth.GenerateJWT(author)

		if err != nil {
			fmt.Println(err)
			return &model.AuthPayload{}, err
		}

		return &payload, nil

	}
	// err := bcrypt.CompareHashAndPassword()
	return &model.AuthPayload{}, nil
}

func (r *queryResolver) CheckUsername(ctx context.Context, username string) (bool, error) {
	author := &model.Author{}
	r.DB.Debug().Where("username = ?", username).First(author)
	if author.Email == "" {
		return true, nil
	}
	return false, nil
}

func (r *queryResolver) CheckEmail(ctx context.Context, email string) (bool, error) {
	author := &model.Author{}
	r.DB.Debug().Where("email = ?", email).First(author)
	if author.Email == "" {
		return true, nil
	}
	return false, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
