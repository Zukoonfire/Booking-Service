package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Zukoonfire/booking-service/database"
)

func GetAllSeats(w http.ResponseWriter, r *http.Request) {
	seats, err := database.GetAllSeats()
	
	if err != nil {
		log.Println("Error Retrieving Seats", err)
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(seats)
	if err != nil {
		log.Println("Error encoding seats to JSON", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
