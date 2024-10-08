package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

type Persona struct {
	Nombre string `json:"nombre"`
	Edad   int    `json:"edad"`
}

func main() {

	http.HandleFunc("/cuenta", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Microservicio Cuenta")
	})

	go receiveMessages()

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func receiveMessages() {

	hostRabbitMQ := os.Getenv("SPRING_RABBITMQ_URL")
	fmt.Println("host rabbitMQ:", hostRabbitMQ)
	conn, err := amqp.Dial(hostRabbitMQ)
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

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var persona Persona
			err := json.Unmarshal(d.Body, &persona)
			if err != nil {
				log.Printf("Error deserializando el mensaje: %s", err)
				continue
			}
			log.Printf("Recibido: Nombre=%s, Edad=%d", persona.Nombre, persona.Edad)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}
