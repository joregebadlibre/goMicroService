// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
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
