package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	httpDelivery "github.com/nouval/library-app/delivery/http"
	"github.com/nouval/library-app/repository/postgres"
	"github.com/nouval/library-app/usecase"
)

func main() {
	if err := godotenv.Overload(); err != nil {
		log.Println("file .env tidak ditemukan")
	}

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASS", "postgres")
	dbName := getEnv("DB_NAME", "backend-gdgoc")

	var dsn string
	if dbPass != "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	} else {
		dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("gagal terhubung dengan database", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("gagal melakukan ping ke database", err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	jwtSecret := getEnv("JWT_SECRET", "super-secret-key")
	timeoutContext := time.Duration(2) * time.Second

	userRepo := postgres.NewPostgresUserRepository(dbConn)
	bookRepo := postgres.NewPostgresBookRepository(dbConn)
	borrowingRepo := postgres.NewPostgresBorrowingRepository(dbConn)

	userUsecase := usecase.NewUserUsecase(userRepo, timeoutContext, jwtSecret)
	bookUsecase := usecase.NewBookUsecase(bookRepo, timeoutContext)
	borrowingUsecase := usecase.NewBorrowingUsecase(borrowingRepo, bookRepo, timeoutContext)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	httpDelivery.NewUserHandler(app, userUsecase)
	httpDelivery.NewBookHandler(app, bookUsecase, jwtSecret)
	httpDelivery.NewBorrowingHandler(app, borrowingUsecase, jwtSecret)

	log.Println("Server running on port 8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
