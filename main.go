package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zukoonfire/booking-service/handlers"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Endpoint Hitt")
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello,Worlddd")
	})
	router.HandleFunc("/seats", handlers.GetAllSeats).Methods("GET")
	router.HandleFunc("/seats/{id}", handlers.GetSeatPricing).Methods("GET")
	router.HandleFunc("/booking", handlers.CreateBooking).Methods("POST")
	router.HandleFunc("/bookings", handlers.GetBookings).Methods("GET").Queries("userIdentifier", "{userIdentifier}")
	log.Fatal(http.ListenAndServe(":8080", router))

}
