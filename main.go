package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	client, err := NewSMSClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/handoff", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "invalid order", http.StatusBadRequest)
			return
		}
		id, err := NotifyHandoff(client, order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"message_id": id, "order_id": order.ID})
	})
	log.Println("marketplace SMS service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
