package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Zukoonfire/booking-service/database"
	"github.com/gorilla/mux"
)

func GetSeatPricing(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seatID := vars["id"]

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

	response := struct {
		SeatPricing *database.SeatPricing `json:"seat_pricing"`
		Percentage  int                   `json:"percentage"`
		Price       sql.NullFloat64       `json:"price"`
	}{
		SeatPricing: pricing,
		Percentage:  percentage,
		Price:       price,
	}

	// Convert the response struct to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Println("Error converting response to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// ========================================
func calculateBookedPercentage(seatClass string) int {

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/booking_service")
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return 0
	}
	defer db.Close()

	var bookedSeats int
	err = db.QueryRow("SELECT COUNT(*) FROM seat WHERE seat_class = ? AND seat_status = 'booked'", seatClass).Scan(&bookedSeats)
	if err != nil {
		log.Println("Error executing query:", err)
		return 0
	}

	var totalSeats int
	err = db.QueryRow("SELECT COUNT(*) FROM seat WHERE seat_class = ?", seatClass).Scan(&totalSeats)
	if err != nil {
		log.Println("Error executing query:", err)
		return 0
	}

	percentage := (bookedSeats * 100) / totalSeats
	return percentage
}

//========================================

func determinePrice(minPrice, normalPrice, maxPrice sql.NullFloat64, percentage int) sql.NullFloat64 {
	if percentage < 40 {
		if minPrice.Valid {
			return minPrice
		} else {
			if normalPrice.Valid {
				return normalPrice
			} else {
				return maxPrice
			}
		}

	} else if percentage >= 40 && percentage <= 60 {
		if normalPrice.Valid {
			return normalPrice
		} else {
			if maxPrice.Valid {
				return maxPrice
			} else {
				return minPrice
			}
		}

	} else {
		if maxPrice.Valid {
			return maxPrice
		} else {
			if normalPrice.Valid {
				return normalPrice
			} else {
				return minPrice
			}
		}
	}
}
