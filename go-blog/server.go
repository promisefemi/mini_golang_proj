package main

import (
	"blog/graph/generated"
	"blog/graph/model"
	"blog/graph/resolver"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

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
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		DB: db}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
