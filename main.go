package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	db, err := bolt.Open("./my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte("answer"), []byte("42"))
		return err
	})

	// app.Post("/user", func(c *fiber.Ctx) error {
	// })

	app.Get("/", func(c *fiber.Ctx) error {
		var result []byte
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			result = b.Get([]byte("answer"))
			return nil
		})
		return c.SendString(string(result))
	})

	app.Listen(":3000")

	// defer db.Close()
}
