package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Seat struct {
	Id             int    `json:"id"`
	SeatIdentifier string `json:"seat_identifier"`
	SeatClass      string `json:"seat_class"`
	SeatStatus     string `json:"seat_status"`
}

func GetAllSeats() ([]Seat, error) {

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/booking_service")
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM seat")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()
	seats := []Seat{}
	for rows.Next() {
		seat := Seat{}
		err := rows.Scan(&seat.Id, &seat.SeatIdentifier, &seat.SeatClass, &seat.SeatStatus)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		seats = append(seats, seat)
		fmt.Println("fromm append")
	}
	return seats, nil
}

type SeatPricing struct {
	Id          int             `json:"id"`
	SeatClass   string          `json:"seat_class"`
	MinPrice    sql.NullFloat64 `json:"min_price"`
	NormalPrice sql.NullFloat64 `json:"normal_price"`
	MaxPrice    sql.NullFloat64 `json:"max_price"`
}

func GetSeatPricing(seatID string) (*SeatPricing, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/booking_service")
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow(`
	SELECT seat_pricing.id, seat_pricing.seat_class, seat_pricing.min_price, seat_pricing.normal_price, seat_pricing.max_price
	FROM seat_pricing
	JOIN seat ON seat_pricing.seat_class = seat.seat_class
	WHERE seat.seat_identifier = ?
	`, seatID)

	pricing := SeatPricing{}
	err = row.Scan(&pricing.Id, &pricing.SeatClass, &pricing.MinPrice, &pricing.NormalPrice, &pricing.MaxPrice)
	if err != nil {
		log.Println("Error scanning row:", err)
		return nil, err
	}

	return &pricing, nil
}

type Booking struct {
	Id     int
	Seats  string
	Name   string
	Email  string
	Phone  string
	Amount float64
}

func CreateBooking(seats []string, name, email, phone string, amount float64) (int, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/booking_service")
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return 0, err
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO booking (seats, name, email, phone, amount) VALUES (?, ?, ?, ?, ?)",
		convertSeatsToString(seats), name, email, phone, amount)
	if err != nil {
		log.Println("Error inserting booking details:", err)
		return 0, err
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		log.Println("Error retrieving booking ID:", err)
		return 0, err
	}

	return int(bookingID), nil
}

func convertSeatsToString(seats []string) string {
	seatsString := ""
	for i, seat := range seats {
		if i > 0 {
			seatsString += ","
		}
		seatsString += seat
	}
	return seatsString
}

func GetBookingsByUserIdentifier(userIdentifier string) ([]Booking, error) {

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/booking_service")
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT booking_id, seats, name, email, phone, amount FROM booking WHERE email = ? OR phone = ?", userIdentifier, userIdentifier)
	if err != nil {
		log.Println("Error retrieving bookings:", err)
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking

	for rows.Next() {
		var booking Booking
		err := rows.Scan(&booking.Id, &booking.Seats, &booking.Name, &booking.Email, &booking.Phone, &booking.Amount)
		if err != nil {
			log.Println("Error scanning booking row:", err)
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		log.Println("Error iterating through bookings rows:", err)
		return nil, err
	}

	return bookings, nil
}
