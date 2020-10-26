package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var (
	vehicleID = mustGetenv("VEHICLE_ID")
)
type status struct {
	VehicleID	string		`json:"vehicle_id"`
	Connected bool 			`json:"connected"`
	Timestamp  int64 		`json:"timestamp"`
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func main() {
	conn, err := amqp.Dial( mustGetenv("RABBITMQ_CON_STRING"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("entity_%s",vehicleID), // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")



	forever := make(chan bool)




	for true  {
		message_body,err  := json.Marshal(status{Connected: RandBool(),Timestamp: time.Now().Unix(),VehicleID: vehicleID})
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(string(message_body)),
			})
		log.Printf(" [x] Sent %s", message_body)
		time.Sleep(60000 * time.Millisecond)
	}
	<-forever
}