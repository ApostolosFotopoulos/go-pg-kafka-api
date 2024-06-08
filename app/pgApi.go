package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func getUsersPG(c *gin.Context) {
    query := `SELECT id, name, money, age, created_at FROM "user"`

    var res []UserPG = []UserPG{}
    rows, err := db.Query(query)

    if (handleError(c, err)) {
        return
    }

    defer rows.Close()

    for rows.Next() {
        var user UserPG
        err = rows.Scan(&user.Id,&user.Name,&user.Money,&user.Age,&user.CreatedAt)

        if (handleError(c, err)) {
            return
        }

        res = append(res, user)
    }

    c.IndentedJSON(http.StatusOK, res)
}

func getUserByIdPG(c *gin.Context) {
    id := c.Param("id")

    query := `SELECT id, name, money, age, created_at
        FROM "user" 
        WHERE id = $1`

    var res UserPG
    row := db.QueryRow(query,id)
    err := row.Scan(&res.Id,&res.Name,&res.Money,&res.Age,&res.CreatedAt)

    if (handleError(c, err)) {
        return
    }

    c.IndentedJSON(http.StatusOK, res)
}

func postUserPG(c *gin.Context) {
    var newUser User

    err := c.BindJSON(&newUser)
    if (handleError(c, err)) {
        return
    }
    
    query := `INSERT INTO "user" (name, age, money) VALUES ($1,$2,$3)
    RETURNING id, name, money, age, created_at`

    var res UserPG
    row := db.QueryRow(query,newUser.Name, newUser.Age, newUser.Money)
    err = row.Scan(&res.Id,&res.Name,&res.Money,&res.Age,&res.CreatedAt)
    
    if (handleError(c, err)) {
        return
    }

    var payload []byte
    payload, err = json.Marshal(res)
    if (handleError(c, err)) {
        return
    }

    msg := kafka.Message{
        Key:   []byte("key"+fmt.Sprint(uuid.New())),
        Value: payload,
    }
    err = kafkaWriter.WriteMessages(context.Background(), msg)
    if (handleError(c, err)) {
        return
    }
    
    c.IndentedJSON(http.StatusCreated, res)
}

func updateUserByIdPG(c *gin.Context) {
    id := c.Param("id")

	var newUser User
    err := c.BindJSON(&newUser)
    if (handleError(c, err)) {
        return
    }

    query := ` UPDATE "user"
    SET name = $2, age = $3, money = $4
    WHERE id = $1
    RETURNING id, name, money, age, created_at`

    var res UserPG
    row := db.QueryRow(query, id, newUser.Name, newUser.Age, newUser.Money)
    err = row.Scan(&res.Id,&res.Name,&res.Money,&res.Age,&res.CreatedAt)
    
    if (handleError(c, err)) {
        return
    }
    
    c.IndentedJSON(http.StatusCreated, res)
}

func deleteUserByIdPG(c *gin.Context) {
    id := c.Param("id")

    query := `DELETE FROM "user"
        WHERE id = $1
        RETURNING id, name, money, age, created_at`

    var res UserPG
    row := db.QueryRow(query,id)
    err := row.Scan(&res.Id,&res.Name,&res.Money,&res.Age,&res.CreatedAt)

    if (handleError(c, err)) {
        return
    }

    c.IndentedJSON(http.StatusOK, res)
}