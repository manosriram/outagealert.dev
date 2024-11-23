package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/manosriram/outagealert.io/pkg/l"
	"github.com/manosriram/outagealert.io/sqlc/db"
)

func TagDefaultResponseHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "max-age=2628288")

		if err := next(c); err != nil {
			return err
		}
		return nil
	}
}

func ToDashboardIfAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := session.Get("session", c)
		if err != nil {
			c.Error(err)
		}

		email := s.Values["email"]
		if email == nil {
			return next(c)
		}
		return c.Redirect(302, "/projects")
	}
}

func IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := session.Get("session", c)
		if err != nil {
			c.Error(err)
			return c.Redirect(302, "/signin")
		}

		email := s.Values["email"]
		if email == nil {
			return c.Redirect(302, "/signin")
		}
		return next(c)
	}
}

func initDB() *db.Queries {
	psqlUser := os.Getenv("POSTGRES_USER")
	psqlPassword := os.Getenv("POSTGRES_PASSWORD")
	psqlPort := os.Getenv("POSTGRES_PORT")
	psqlDatabase := os.Getenv("POSTGRES_DATABASE")
	psqlHost := os.Getenv("POSTGRES_HOST")

	psqlString := fmt.Sprintf("user=%s password=%s port=%s database=%s sslmode=disable host=%s", psqlUser, psqlPassword, psqlPort, psqlDatabase, psqlHost)
	l.Log.Info("psql string ", psqlString)
	config, err := pgxpool.ParseConfig(psqlString)
	if err != nil {
		l.Log.Infof("Unable to parse connection string: %v", err)
	}
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	dbconn := db.New(pool)
	return dbconn
}
