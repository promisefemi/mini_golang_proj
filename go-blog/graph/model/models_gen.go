// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthPayload struct {
	Token  string  `json:"token"`
	Author *Author `json:"author"`
}

type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthorInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Post struct {
	UUID     string  `json:"uuid"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	AuthorID string  `json:"author_id"`
	Author   *Author `json:"author"`
}

type PostInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID string `json:"author_id"`
}
