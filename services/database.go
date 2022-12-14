package services

import (
	"log"

	"github.com/kamalshkeir/korm"
	"github.com/kamalshkeir/muzzsol/settings"
	"github.com/kamalshkeir/muzzsol/types"
	"github.com/kamalshkeir/mysqldriver"
)

func init() {
	// this will import driver once 
	mysqldriver.Use()
}

// DatabaseService wrapper for database service
type DatabaseService struct{}

// NewDatabaseService return new database service
func NewDatabaseService() *DatabaseService {
	return &DatabaseService{}
}

// Init make a connection to db
func (dbs *DatabaseService) Init() *DatabaseService {
	err := korm.New(settings.Config.DB.Type, settings.Config.DB.Name, settings.Config.DB.Dsn)
	if err != nil {
		log.Fatal("NewDatabaseService error:",err)
	}
	// disable cache because data will change a lot (on every swipe)
	korm.DisableCache()
	return dbs
}

// Migrate auto migrate all models
func (dbs *DatabaseService) Migrate() error {
	err := korm.AutoMigrate[types.User]("users")
	if err != nil {
		return err
	}
	err = korm.AutoMigrate[types.Swipe]("swipes")
	if err != nil {
		return err
	}
	err = korm.AutoMigrate[types.Location]("locations")
	if err != nil {
		return err
	}
	return nil
}