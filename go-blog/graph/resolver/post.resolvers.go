package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"blog/graph/auth"
	"blog/graph/model"
	"context"
	"errors"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input *model.PostInput) (*model.Post, error) {
	// Validate author
	author := auth.ForContext(ctx)
	if author == nil || author.Email == "" {
		return &model.Post{}, errors.New("Access Denied")
	}

	post := &model.Post{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: input.AuthorID,
	}

	err := r.DB.Create(post).Error

	if err != nil {
		return nil, errors.New("Unable to create post at this time")
	}

	post.Author = author

	return post, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, input *model.PostInput, uuid string) (*model.Post, error) {
	if author := auth.ForContext(ctx); author == nil || author.Email == "" {
		return &model.Post{}, errors.New("Access Denied")
	}

	post := &model.Post{
		UUID:     "",
		AuthorID: "",
		Title:    input.Title,
		Content:  input.Content,
	}

	result := r.DB.Debug().Where("uuid = ?", uuid).Omit("uuid", "author_id").Updates(post)

	r.DB.Where("uuid = ?", uuid).Find(post)

	if result.Error != nil {
		return nil, result.Error
	}
	author := &model.Author{}
	r.DB.Debug().Where("id = ?", post.AuthorID).Find(author)
	post.Author = author
	return post, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, uuid string) (bool, error) {
	if author := auth.ForContext(ctx); author == nil || author.Email == "" {
		return false, errors.New("Access Denied")
	}
	err := r.DB.Where("uuid = ?", uuid).Delete(&model.Post{}).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) GetMany(ctx context.Context, limit int, page int) ([]*model.Post, error) {
	var posts []*model.Post

	offset := (limit * page) - limit

	r.DB.Limit(limit).Offset(offset).Find(&posts)

	if len(posts) > 0 {
		for _, post := range posts {
			author := &model.Author{}
			r.DB.Where("id = ?", post.AuthorID).Find(author)
			post.Author = author
		}
	}

	return posts, nil
}

func (r *queryResolver) GetPost(ctx context.Context, uuid string) (*model.Post, error) {
	post := &model.Post{}
	err := r.DB.Where("uuid = ?", uuid).Take(post).Error

	if err != nil {
		return nil, err
	}
	author := &model.Author{}
	r.DB.Debug().Where("id = ?", post.AuthorID).Find(author)

	post.Author = author
	return post, nil
}
