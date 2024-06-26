package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var db *sql.DB

type DatabaseConfig struct {
	DBDriver string
	DBUser   string
	DBPass   string
	DBName   string
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST",
	}))

	app.Post("/grade", Grade)

	if err := app.Listen(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func Database() (*sql.DB, error) {
	config := DatabaseConfig{
		DBDriver: "mysql",
		DBUser:   "root",
		DBPass:   "",
		DBName:   "shoplek",
	}

	db, err := sql.Open(config.DBDriver,
		fmt.Sprintf("%s:%s@/%s",
			config.DBUser,
			config.DBPass,
			config.DBName))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Score    int    `json:"score"`
	Grade    string `json:"grade"`
}

func Score(c *fiber.Ctx, user User) error {
	db, err := Database()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection error"})
	}
	defer db.Close()

	// Insert user data into the database
	query := "INSERT INTO score (username, name, score, grade) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, user.Username, user.Name, user.Score, user.Grade)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert data"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":   "OK",
		"username": user.Username,
		"name":     user.Name,
		"score":    user.Score,
		// "grade":    user.Grade,
	})
}

func Grade(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Calculate grade
	switch {
	case user.Score >= 90:
		user.Grade = "A"
	case user.Score >= 80:
		user.Grade = "B"
	case user.Score >= 70:
		user.Grade = "C"
	case user.Score >= 60:
		user.Grade = "D"
	default:
		user.Grade = "F"
	}

	// Insert the user data along with the grade into the database
	return Score(c, user)
}
