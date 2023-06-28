package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Zukoonfire/booking-service/database"
)

type CreateBookingRequest struct {
	Seats []string `json:"seats"`
	Name  string   `json:"name"`
	Email string   `json:"email"`
	Phone string   `json:"phone"`
}

type CreateBookingResponse struct {
	BookingID int     `json:"booking_id"`
	Amount    float64 `json:"amount"`
}

func CreateBooking(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ggggggggg")
	var req CreateBookingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if len(req.Seats) == 0 || req.Name == "" || req.Email == "" || req.Phone == "" {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	var totalAmount float64
	for _, seatID := range req.Seats {
		pricing, err := database.GetSeatPricing(seatID)
		if err != nil {
			log.Println("Error retrieving seat pricing:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if pricing == nil {
			http.NotFound(w, r)
			return
		}
		percentage := calculateBookedPercentage(pricing.SeatClass)
		price := determinePrice(pricing.MinPrice, pricing.NormalPrice, pricing.MaxPrice, percentage)
		totalAmount += price.Float64
	}

	bookingID, err := database.CreateBooking(req.Seats, req.Name, req.Email, req.Phone, totalAmount)
	if err != nil {
		log.Println("Error creating booking:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := CreateBookingResponse{
		BookingID: bookingID,
		Amount:    totalAmount,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error converting response to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}
