package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "0.0.1"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func (app application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	cfg := app.config
	fmt.Fprintln(w, "status: OK")
	fmt.Fprintf(w, "env: %s\n", cfg.env)
	fmt.Fprintf(w, "version: %s\n", version)
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Application environment (dev|stg|prod)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{config: cfg, logger: logger}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, server.Addr)

	err := server.ListenAndServe()
	logger.Fatal(err)
}
