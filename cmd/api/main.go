package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/zmwilliam/greenlight/internal/data"
	"github.com/zmwilliam/greenlight/internal/jsonlog"
	"github.com/zmwilliam/greenlight/internal/mailer"
)

const version = "0.0.1"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Application environment (dev|stg|prod)")

	flag.StringVar(
		&cfg.db.dsn,
		"db-dsn",
		getEnv("GREENLIGHT_DB_DSN", "postgres://usr:pwd@localhost:5432/db"),
		"PostgreSQL DSN",
	)
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(
		&cfg.db.maxIdleTime,
		"db-max-idle-time",
		"15m",
		"PostgreSQL max connection idle time",
	)

	flag.Float64Var(
		&cfg.limiter.rps,
		"limiter-rps",
		2,
		"Rate limiter maximum requests per second",
	)
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 2, "Rate limiter maximum requests per second")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(
		&cfg.smtp.host,
		"smtp-host",
		getEnv("GREENLIGHT_SMTP_HOST", "smtp.mailtrap.io"),
		"SMTP host",
	)

	flag.IntVar(&cfg.smtp.port, "smtp-port", getEnvInt("GREENLIGHT_SMTP_PORT", 25), "SMTP port")
	flag.StringVar(
		&cfg.smtp.username,
		"smtp-username",
		getEnv("GREENLIGHT_SMTP_USERNAME", "smtp-user"),
		"SMTP username",
	)
	flag.StringVar(
		&cfg.smtp.password,
		"smtp-password",
		getEnv("GREENLIGHT_SMTP_PASSWORD", "smtp-password"),
		"SMTP password",
	)
	flag.StringVar(
		&cfg.smtp.sender,
		"smtp-sender",
		getEnv("GREENLIGHT_SMTP_SENDER", "Greenlight <no-reply@greenlight.zmwilliam.com>"),
		"SMTP sender",
	)

	flag.Func("cors-trusted-origins", "Trusted CORS origins", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(
			cfg.smtp.host,
			cfg.smtp.port,
			cfg.smtp.username,
			cfg.smtp.password,
			cfg.smtp.sender,
		),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getEnvInt(env_name string, default_val int) int {
	v, err := strconv.Atoi(os.Getenv(env_name))
	if err != nil {
		return default_val
	}

	return v
}

func getEnv(env_name, default_val string) string {
	if v := os.Getenv(env_name); v != "" {
		return v
	}

	return default_val
}
