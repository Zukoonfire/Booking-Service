package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/gorilla/mux"

	"github.com/Zukoonfire/booking-service/database"
)

type GetBookingsResponse struct {
	Bookings []database.Booking `json:"bookings"`
}

func GetBookings(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user identifier from the query parameters
	userIdentifier := r.URL.Query().Get("userIdentifier")
	// vars := mux.Vars(r)
	// userIdentifier := vars["userIdentifier"]

	// Validate the user identifier
	if userIdentifier == "" {
		fmt.Println("testinggggg")
		http.Error(w, "User identifier is required", http.StatusBadRequest)
		return
	}

	// Retrieve the bookings from the database
	bookings, err := database.GetBookingsByUserIdentifier(userIdentifier)
	if err != nil {
		log.Println("Error retrieving bookings:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := GetBookingsResponse{
		Bookings: bookings,
	}

	// Convert the response struct to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error converting response to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
