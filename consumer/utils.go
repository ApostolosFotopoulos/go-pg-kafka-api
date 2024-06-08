package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

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

    query = `CREATE TABLE IF NOT EXISTS "userKafka" (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        age INT NOT NULL,
        money NUMERIC(10, 2) NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    )`

    _, err = db.Exec(query)
    if (err != nil) {
        log.Fatal(err.Error())
    }
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

