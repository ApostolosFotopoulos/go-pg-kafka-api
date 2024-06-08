package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var users []UserPG = []UserPG{{Id:1,Name: "John", Age:18, Money:100},{Id:2,Name:"Jack", Age:35, Money:50000}, {Id:3,Name:"Ann", Age:22, Money:1000}}

func getUsers(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, users)
}

func getUserById(c *gin.Context) {
    id := c.Param("id")

    for _, u := range users {
        if idInt, err := strconv.Atoi(id); err == nil && u.Id == idInt {
            c.IndentedJSON(http.StatusOK, u)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func postUser(c *gin.Context) {
    var newUser User
    err := c.BindJSON(&newUser)

    if (handleError(c, err)) {
        return
    }

    newId := users[len(users)-1].Id + 1
    var userToInsert UserPG = UserPG{Id:newId,Name: newUser.Name, Age:newUser.Age, Money:newUser.Money}

    users = append(users, userToInsert)
    c.IndentedJSON(http.StatusCreated, users)
}

func updateUserById(c *gin.Context) {
    id := c.Param("id")

	var newUser User
    err := c.BindJSON(&newUser)
    if (handleError(c, err)) {
        return
    }

    for i, u := range users {
        if idInt, err := strconv.Atoi(id); err == nil && u.Id == idInt {
            users[i].Name = newUser.Name
            users[i].Age = newUser.Age
            users[i].Money = newUser.Money
			c.IndentedJSON(http.StatusOK, users)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func deleteUserById(c *gin.Context) {
    id := c.Param("id")

    for i, u := range users {
        if idInt, err := strconv.Atoi(id); err == nil && u.Id == idInt {
            users[i] = users[len(users)-1]
			users = users[:len(users)-1]
			c.IndentedJSON(http.StatusOK, users)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
}