package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// Object pada Redis
type User struct {
	RealName string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"` // password dalam bentuk hash SHA1
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Global redis client
var rdb *redis.Client
var ctx = context.Background()

func loginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	key := fmt.Sprintf("login_%s", req.Username)
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	var user User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse user data",
		})
	}

	h := sha1.New()
	h.Write([]byte(req.Password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))

	if hashedPassword != user.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Login success",
		"realname": user.RealName,
		"email":    user.Email,
	})
}

func main() {
	// Init redis
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // sesuaikan
		DB:   0,
	})

	app := fiber.New()

	// Endpoint login
	app.Post("/login", loginHandler)

	log.Fatal(app.Listen(":3000"))
}
