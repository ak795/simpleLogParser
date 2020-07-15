package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"logParser/parser"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	msg := "hello from the parser service !!!"
	_, err := w.Write([]byte(msg))
	if err != nil {
		fmt.Print("some error occurred: ", err.Error())
	}
}

func main() {
	// Router configuration
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheck)
	r.HandleFunc("/upload", parser.LogUploadHandler)

	// http server configuration
	server := &http.Server{
		Addr:         ":7212",
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      r,
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
		fmt.Print("Server started....")
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	fmt.Print("Shutting down server")
	os.Exit(0)
}
