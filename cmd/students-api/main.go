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
)

func main() {
	// our tasks:

	cf := config.MustLoad() // load config
	// 2. database setup
	// 3. setup http server
	router := http.NewServeMux()
	// route 1: health check
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Students API"))
	})

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
