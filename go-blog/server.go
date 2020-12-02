package main

import (
	"blog/graph/auth"
	"blog/graph/generated"
	"blog/graph/model"
	"blog/graph/resolver"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/go-chi/chi"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8080"

var db *gorm.DB

func InitDB() {
	var err error
	dbDetails := "user=promise password=promise port=5432 database=go_blog sslmode=false"
	db, err = gorm.Open(postgres.Open(dbDetails), &gorm.Config{})
	if err != nil {
		log.Fatalf("%s", err)
	}

	// db.Exec("CREATE DATABASE go_blog")
	db.Debug().AutoMigrate(&model.Author{}, &model.Post{})

}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	InitDB()

	router := chi.NewRouter()
	router.Use(auth.Middleware(db))

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		DB: db}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	fmt.Printf("\n %s \n", "===========-------------------------==============")
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
