package driver

import (
	"database/sql"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	SQL *gorm.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL creates database pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	sqlDB, err := d.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(maxOpenDbConn)
	sqlDB.SetMaxIdleConns(maxIdleDbConn)
	sqlDB.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = d

	err = testDB(sqlDB)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*gorm.DB, error) {
	gc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db, err := gc.DB()
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return gc, nil
}

// testDB tries to ping the database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
