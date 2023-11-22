package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GingFreecss2/rss-go-server/api"
	"github.com/GingFreecss2/rss-go-server/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {

	portStr, exists := os.LookupEnv("PORT")
	if !exists {
		log.Fatal("PORT not found in .env file")
	}

	dbURL, exists := os.LookupEnv("DB_URL")
	if !exists {
		log.Fatal("dbURL not found in .env file")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database ", err)
	}

	db := database.New(conn)
	apiCfg := api.ApiConfig{
		DB: db,
	}

	go api.StartScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/ready", api.HandlerReadiness)
	v1Router.Get("/err", api.HandlerErr)

	v1Router.Post("/users", apiCfg.HandlerCreateUser)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUser))

	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeed))
	v1Router.Get("/feeds", (apiCfg.HandlerGetFeeds))

	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteFeedFollow))

	v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.HandlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}

	log.Printf("Server starting on port %v", portStr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port: ", portStr)

}
