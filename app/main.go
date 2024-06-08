package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type User struct {
	Name string `json:"name"`
	Age uint8 `json:"age"`
	Money float64 `json:"money"`
}

type UserPG struct {
    Id int `json:"id"`
	Name string `json:"name"`
	Age uint8 `json:"age"`
	Money float64 `json:"money"`
    CreatedAt string `json:"created_at"` 
}

var db *sql.DB
var kafkaWriter *kafka.Writer
const kafkaTopic = "user"
const kafkaUrl = "localhost:9092"
const version = 2

func main() {
	router := gin.Default()
    if (version == 1) {
        router.GET("/users", getUsers)
        router.GET("/users/:id", getUserById)
        router.PUT("/users/:id", updateUserById)
        router.DELETE("/users/:id", deleteUserById)
        router.POST("/users", postUser)
    } else if (version == 2) {
        setUpPG()
        defer db.Close()

        setUpKafka()
        defer kafkaWriter.Close()

        router.GET("/users", getUsersPG)
        router.GET("/users/:id", getUserByIdPG)
        router.PUT("/users/:id", updateUserByIdPG)
        router.DELETE("/users/:id", deleteUserByIdPG)
        router.POST("/users", postUserPG)
    }
    
    router.Run("localhost:8080")
}
