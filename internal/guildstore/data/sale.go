package data

import (
	"time"

	accountregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/account"
	guildregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/guild"
	itemlinkregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/itemlink"
)

type Sale struct {
	SellerId   uint
	Seller     *accountregistry.Account
	ItemLinkId uint
	ItemLink   *itemlinkregistry.ItemLink
	Kiosk      bool
	GuildId    uint
	Guild      *guildregistry.Guild
	BuyerId    uint
	Buyer      *accountregistry.Account
	Id         string
	Quantity   uint
	Timestamp  time.Time
	Price      int
}
