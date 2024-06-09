package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func handleError(c *gin.Context, err error) bool {
    if (err !=nil) {
        fmt.Println(err.Error())
        if (err == sql.ErrNoRows) {
            c.JSON(http.StatusInternalServerError, "No record found")
        } else {
            c.JSON(http.StatusInternalServerError, "Internal server error")
        }
        return true
    }

    return false
}

func setUpPG() {
    connStr := "postgres://postgres:postgres123@localhost:5432/goapipg?sslmode=disable"

    newDb, err := sql.Open("postgres", connStr)

    if (err != nil) {
        log.Fatal(err.Error())
    }

    db=newDb
    err = db.Ping()
    if (err != nil) {
        log.Fatal(err.Error())
    }

    createUserTable()
}

func createUserTable() {
    query := `CREATE TABLE IF NOT EXISTS "user" (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        age INT NOT NULL,
        money NUMERIC(10, 2) NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    )`

    _, err := db.Exec(query)
    if (err != nil) {
        log.Fatal(err.Error())
    }
}

func setUpKafka() {
    conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", kafkaTopic, 0)
    if err != nil {
        log.Fatal(err.Error())
    }
    defer conn.Close()

	kafkaWriter = newKafkaWriter(kafkaUrl, kafkaTopic)
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

