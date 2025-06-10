package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Tamiris15/university-blog/pkg/application"
	"github.com/Tamiris15/university-blog/pkg/infrastructure/repository/sqldb"
	"github.com/Tamiris15/university-blog/pkg/infrastructure/routes"
	"github.com/gorilla/mux"
)

func main() {
	databaseConnection, err := sql.Open("sqlite3", "./university_blog.db")
	if err != nil {
		log.Fatalf("Критическая ошибка при открытии базы данных: %v", err)
	}
	defer databaseConnection.Close()

	if err := createTables(databaseConnection); err != nil {
		log.Fatalf("Критическая ошибка при создании таблиц базы данных: %v", err)
	}

	userRepository := sqldb.NewUserRepository(databaseConnection)
	postRepository := sqldb.NewPostRepository(databaseConnection)
	commentRepository := sqldb.NewCommentRepository(databaseConnection)

	userService := application.NewUserService(userRepository)
	postService := application.NewPostService(postRepository)
	commentService := application.NewCommentService(commentRepository)

	userRouteHandler := routes.NewUserHandler(userService)
	postRouteHandler := routes.NewPostHandler(postService)
	commentRouteHandler := routes.NewCommentHandler(commentService)

	apiRouter := mux.NewRouter()

	apiRouter.HandleFunc("/users", userRouteHandler.Create).Methods("POST")
	apiRouter.HandleFunc("/users", userRouteHandler.GetAll).Methods("GET")
	apiRouter.HandleFunc("/users/{id:[0-9]+}", userRouteHandler.GetByID).Methods("GET")
	apiRouter.HandleFunc("/users/{id:[0-9]+}", userRouteHandler.Update).Methods("PUT")
	apiRouter.HandleFunc("/users/{id:[0-9]+}", userRouteHandler.Delete).Methods("DELETE")

	apiRouter.HandleFunc("/posts", postRouteHandler.Create).Methods("POST")
	apiRouter.HandleFunc("/posts", postRouteHandler.GetAll).Methods("GET")
	apiRouter.HandleFunc("/posts/{id:[0-9]+}", postRouteHandler.GetByID).Methods("GET")
	apiRouter.HandleFunc("/posts/author/{author_id:[0-9]+}", postRouteHandler.GetByAuthorID).Methods("GET")
	apiRouter.HandleFunc("/posts/{id:[0-9]+}", postRouteHandler.Update).Methods("PUT")
	apiRouter.HandleFunc("/posts/{id:[0-9]+}", postRouteHandler.Delete).Methods("DELETE")

	apiRouter.HandleFunc("/posts/{post_id:[0-9]+}/comments", commentRouteHandler.Create).Methods("POST")
	apiRouter.HandleFunc("/comments/{id:[0-9]+}", commentRouteHandler.GetByID).Methods("GET")
	apiRouter.HandleFunc("/posts/{post_id:[0-9]+}/comments", commentRouteHandler.GetByPostID).Methods("GET")
	apiRouter.HandleFunc("/comments/{id:[0-9]+}", commentRouteHandler.Update).Methods("PUT")
	apiRouter.HandleFunc("/comments/{id:[0-9]+}", commentRouteHandler.Delete).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Веб-сервер успешно запущен и ожидает подключения на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, apiRouter))
}

func createTables(databaseConnection *sql.DB) error {
	_, err := databaseConnection.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = databaseConnection.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			author_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	_, err = databaseConnection.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			post_id INTEGER NOT NULL,
			author_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES users(id)
		)
	`)
	return err
}
