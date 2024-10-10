package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		var persona Persona

		// Leer el cuerpo de la solicitud
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "No se pudo leer el cuerpo de la solicitud", http.StatusBadRequest)
			return
		}
		// Deserializar el JSON en el objeto Persona
		err = json.Unmarshal(body, &persona)
		if err != nil {
			http.Error(w, "Formato de JSON inválido", http.StatusBadRequest)
			return
		}

		sendMessageCuenta(persona)

	})

	go receiveMessages()

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func receiveMessages() {

	hostRabbitMQ := os.Getenv("SPRING_RABBITMQ_URL")
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

////agregadp

func sendMessageCuenta(persona Persona) {
	hostRabbitMQ := os.Getenv("SPRING_RABBITMQ_URL")
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
		"cuenta_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	body, err := json.Marshal(persona)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Fatal(err)
	}
}
