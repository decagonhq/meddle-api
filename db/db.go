package db

import (
	"fmt"
	"log"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	DB *gorm.DB
}

func GetDB(c *config.Config) *GormDB {
	gormDB := &GormDB{}
	gormDB.Init(c)
	return gormDB
}

func (g *GormDB) Init(c *config.Config) {
	g.DB = getPostgresDB(c)

	if err := migrate(g.DB); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}
}

func getPostgresDB(c *config.Config) *gorm.DB {
	log.Printf("Connecting to postgres: %+v", c)


	//log.Println("Connecting to postgres: am here", c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB)
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Africa/Lagos",
		c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)
	postgresDB, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{
		Logger: logger.Default,
	})
	if err != nil {
		log.Fatalf("Error connecting to postgres: %s ", err)
	}
	return postgresDB
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.BlackList{})
	if err != nil {
		return fmt.Errorf("migrations error: %v", err)
	}

	return nil
}
