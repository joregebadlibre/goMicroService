package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	// Llamar al servicio persona

	personaData := map[string]interface{}{
		"nombre": "Luis",
		"edad":   30,
	}
	personaJSON, _ := json.Marshal(personaData)
	_, err := http.Post("http://localhost:8081/persona", "application/json", bytes.NewBuffer(personaJSON))
	if err != nil {
		log.Fatalf("Error llamando al servicio persona: %s", err)
	}

	/*
		// Llamar al servicio cuenta
		resp, err := http.Get("http://localhost:8082/cuenta")
		if err != nil {
			log.Fatalf("Error llamando al servicio cuenta: %s", err)
		}
		defer resp.Body.Close()

		// Procesar la respuesta del servicio cuenta
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		log.Printf("Respuesta del servicio cuenta: %v", result)
	*/
}
