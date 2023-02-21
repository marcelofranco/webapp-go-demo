package dbrepo

import (
	"database/sql"

	"github.com/marcelofranco/webapp-go-demo/internal/config"
	"github.com/marcelofranco/webapp-go-demo/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(a *config.AppConfig, conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
