// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

type Persona struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/persona", func(w http.ResponseWriter, r *http.Request) {
		persona := Persona{ID: "1", Name: "Juan Perez"}
		json.NewEncoder(w).Encode(persona)
	})

	log.Println("Persona service running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendMessage(message string) {
	conn, err := amqp.Dial("RABBITMQ_URL")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"persona_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
