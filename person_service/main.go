package main

import (
	"encoding/json"
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

	http.HandleFunc("/persona", func(w http.ResponseWriter, r *http.Request) {

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

		persona.Nombre = "Modificado antes de enviar a cuenta"
		// Serializar el objeto Persona a JSON para la respuesta
		response, err := json.Marshal(persona)
		if err != nil {
			http.Error(w, "Error al serializar el objeto Persona", http.StatusInternalServerError)
			return
		}
		// Establecer el encabezado de contenido y devolver la respuesta
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)

		sendMessage(persona)
	})

	go receiveMessagesCuenta()

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func sendMessage(persona Persona) {
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

func receiveMessagesCuenta() {
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
