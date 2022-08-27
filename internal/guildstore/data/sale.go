package data

import "time"

type Sale struct {
	SellerId   int
	ItemLinkId int
	Kiosk      bool
	GuildId    int
	BuyerId    int
	Id         string
	Quantity   uint
	Timestamp  time.Time
	Price      int
}
