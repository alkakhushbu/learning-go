package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-mgmt-v2/db"
	"task-mgmt-v2/handlers"
	"time"
)

func main() {
	//set up database connection
	log.Println("Creating new connection to DB.....")
	conn, err := db.NewConn()
	if err != nil {
		log.Println("Error creating connection to DB:", err)
		panic(err)
	}
	//Ping database
	conn.Ping(context.Background())

	setupLog()

	//set up handler
	log.Println("Starting Task Management Service")

	// Channel to listen for OS signals (like SIGTERM, SIGINT) for graceful shutdown
	shutdown := make(chan os.Signal, 1)

	// Register the shutdown channel to receive specific system interrupt signals
	signal.Notify(shutdown, os.Kill, os.Interrupt, syscall.SIGTERM)

	// Channel to capture server errors during runtime, like port already being used
	serverError := make(chan error)

	api := http.Server{
		Addr:              ":8084",
		ReadHeaderTimeout: time.Second * 200,
		WriteTimeout:      time.Second * 200,
		IdleTimeout:       time.Second * 200,
		Handler:           handlers.SetupGINRoutes(conn),
	}

	// Goroutine to handle server startup and listen for incoming requests
	go func() {
		serverError <- api.ListenAndServe()
	}()

	select {
	case <-serverError:
		// Panic if the server fails to start
		panic(err)
	case <-shutdown:
		log.Println("Graceful Shutdown Server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		//Shutdown gracefully shuts down the server without interrupting any active connections.
		//Shutdown works by first closing all open listeners, then closing all idle connections,
		err := api.Shutdown(ctx)
		if err != nil {
			// force close
			err := api.Close()
			panic(err)
		}
	}

}
func setupLog() {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{
				AddSource: true, Level: slog.LevelDebug,
			}))
	slog.SetDefault(logger)
}
