package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

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
var kafkaReader *kafka.Reader
const kafkaTopic = "user"
const kafkaUrl = "localhost:9092"
const kafkaGroupId = "g1"

func main() {

	setUpPG()
	defer db.Close()

	kafkaReader = getKafkaReader(kafkaUrl, kafkaTopic, kafkaGroupId)

	defer kafkaReader.Close()

	fmt.Println("start consuming ... !!")
	for {
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var insertedUser = &UserPG{}

   		err = json.Unmarshal(m.Value, insertedUser)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(insertedUser.Name)
	}
}
