package main

import (
	"fmt"
	"net/http"

	"github.com/Aytaditya/students-api-golang/internal/config"
)

func main() {
	fmt.Println("Students API Service")

	// our task
	// 1. load config
	cf := config.MustLoad()
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

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
