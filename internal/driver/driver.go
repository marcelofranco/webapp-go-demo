package driver

import (
	"database/sql"
	"errors"
	"time"

	"github.com/marcelofranco/webapp-go-demo/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	SQL *sql.DB
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

	err = runMigrations(d)
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

	dbConn.SQL = sqlDB

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

func runMigrations(d *gorm.DB) error {
	err := d.AutoMigrate(&models.User{},
		&models.Reservation{},
		&models.Restriction{},
		&models.Room{},
		&models.RoomRestriction{})

	if err != nil {
		return err
	}

	if d.Migrator().HasTable(&models.Room{}) {
		if err = d.First(&models.Room{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			var rooms = []models.Room{
				{
					RoomName: "General's Quarters",
				},
				{
					RoomName: "Major's Suite",
				},
			}
			for _, r := range rooms {
				if err = d.Create(&r).Error; err != nil {
					return err
				}
			}
		}
	}

	if d.Migrator().HasTable(&models.Restriction{}) {
		if err = d.First(&models.Restriction{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			var restrictions = []models.Restriction{
				{
					RestrictionName: "Reservation",
				},
				{
					RestrictionName: "Owner Block",
				},
			}
			for _, rr := range restrictions {
				if err = d.Create(&rr).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}
