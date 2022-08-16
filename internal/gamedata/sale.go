package gamedata

import "time"

type Sale struct {
	Seller    int
	ItemLink  int
	Kiosk     bool
	Guild     int
	Buyer     int
	Id        string
	Quantity  uint
	Timestamp time.Time
	Price     int
}
