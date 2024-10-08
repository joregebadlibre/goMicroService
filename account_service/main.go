// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

type Cuenta struct {
	ID        string  `json:"id"`
	PersonaID string  `json:"persona_id"`
	Balance   float64 `json:"balance"`
}

func main() {
	http.HandleFunc("/cuenta", func(w http.ResponseWriter, r *http.Request) {
		cuenta := Cuenta{ID: "1", PersonaID: "1", Balance: 1000.0}
		json.NewEncoder(w).Encode(cuenta)
	})

	log.Println("Cuenta service running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func receiveMessages() {
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

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
}

func main() {
	go receiveMessages()
	// Resto del código...
}
