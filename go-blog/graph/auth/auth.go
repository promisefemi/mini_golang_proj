package auth

import (
	"blog/graph/model"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type contextKey struct {
	name string
}

var authorContext = &contextKey{"author"}

var jwtAppToken = []byte("Secret_key_for_blog_app")

type jwtS struct {
	AuthorID string `json:"author_id"`
	jwt.StandardClaims
}

// GenerateJWT is a function for generating the AuthPayload before loggin a user in.
func GenerateJWT(author *model.Author) (model.AuthPayload, error) {

	expiration := time.Now().Add(time.Minute * 5)

	claims := &jwtS{
		AuthorID: author.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tknKey, err := token.SignedString(jwtAppToken)

	if err != nil {
		return model.AuthPayload{}, err
	}

	return model.AuthPayload{
		Token:  tknKey,
		Author: author,
	}, nil

}

func validateJWT(tokenKey string) (string, error) {
	claims := &jwtS{}
	token, err := jwt.ParseWithClaims(tokenKey, claims, func(key *jwt.Token) (interface{}, error) {
		return jwtAppToken, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("Invalid Token")
	}
	return claims.AuthorID, nil

}

// Middleware to validate the JWT in the Authorization header create a new context and continue
func Middleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			key := r.Header.Get("Authorization")

			if strings.TrimSpace(key) == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenKey := strings.Split(key, " ")[1]
			userid, err := validateJWT(tokenKey)

			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			author := &model.Author{}

			db.Find(author, "id = ?", userid)

			ctx := context.WithValue(r.Context(), authorContext, author)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})

	}
}

// ForContext gets the author from the context
func ForContext(ctx context.Context) *model.Author {
	raw, _ := ctx.Value(authorContext).(*model.Author)
	return raw
}
