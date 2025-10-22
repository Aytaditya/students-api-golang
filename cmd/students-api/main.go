package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aytaditya/students-api-golang/internal/config"
	"github.com/Aytaditya/students-api-golang/internal/http/handlers/student"
	"github.com/Aytaditya/students-api-golang/internal/storage/sqlite"
)

func main() {
	// our tasks:

	cf := config.MustLoad() // load config (this return is pointer to config struct)
	// 2. database setup
	storage, er := sqlite.New(cf)
	if er != nil {
		slog.Error("Failed to connect to database", "error", er.Error())
		return
	}
	slog.Info("Connected to database successfully", slog.String("storage_path", cf.StoragePath))
	// 3. setup http server
	router := http.NewServeMux() // create a new http serve mux (router)

	// route 1: health check
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	// route 2: student creation menthod
	router.HandleFunc("POST /api/students", student.New(storage))

	// route 2 (get student by id)
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	server := http.Server{
		Addr:    cf.HttpServer.Addr,
		Handler: router,
	}

	fmt.Println("Server Running at:", cf.HttpServer.Addr)

	done := make(chan os.Signal, 1) //buffered channel

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // notify channel when interrupt signal is received

	// start server in a goroutine and blocking main goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	// basically what is happening when we interupt (ctrl+c) the program gets stopped immediatly but we want gracefull shutdown for that we use goroutine and channels
	// we are running server in go routine and we are blocking it using done channel. it will get unblocked and below code will run when we get interrupt signal

	<-done // block main goroutine until we get signal

	slog.Info("Server Stopped")

	// 5 second timer based shutdown (new request wont be accepted and existing request will be given 5 seconds to complete)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shutdown server", "error", err.Error())
		return
	}

	slog.Info("Server Exited Gracefully")

}
