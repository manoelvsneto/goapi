package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "goapi/docs" // Import the docs package
	"goapi/handler"
	"log"
	"os"
)

var db *sql.DB
var err error

// @title Go CRUD SQL Server Swagger API
// @version 1.0
// @description This is a sample server for a Go CRUD API using SQL Server and Swagger.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	// Load connection string from environment variable
	connectionString := os.Getenv("DATABASE_CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("DATABASE_CONNECTION_STRING is not set")
	}

	// Connect to SQL Server
	db, err = sql.Open("sqlserver", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	fmt.Println("Connected to SQL Server!")

	// Set up the router
	router := gin.Default()

	// Set database for handlers
	handler.SetDB(db)

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes
	router.GET("/people", handler.GetPeople)
	router.GET("/people/:id", handler.GetPerson)
	router.POST("/people", handler.CreatePerson)
	router.PUT("/people/:id", handler.UpdatePerson)
	router.DELETE("/people/:id", handler.DeletePerson)

	// Start the server
	router.Run(":8080")
}
