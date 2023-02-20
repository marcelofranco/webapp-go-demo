package dbrepo

import (
	"github.com/marcelofranco/webapp-go-demo/internal/config"
	"github.com/marcelofranco/webapp-go-demo/internal/repository"
	"gorm.io/gorm"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *gorm.DB
}

func NewPostgresRepo(a *config.AppConfig, conn *gorm.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
