package data

import (
	"time"

	accountregistry "github.com/techyshishy/nirn-revenue-service/internal/guildstore/data/registry/account"
	guildregistry "github.com/techyshishy/nirn-revenue-service/internal/guildstore/data/registry/guild"
	itemlinkregistry "github.com/techyshishy/nirn-revenue-service/internal/guildstore/data/registry/itemlink"

	pbs "github.com/techyshishy/nirn-revenue-service/gen/api/proto/sale/v1"
	"google.golang.org/protobuf/proto"
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

func (s *Sale) Proto() *pbs.Sale {
	return &pbs.Sale{
		Seller:    *proto.String(s.Seller.Name),
		ItemLink:  *proto.String(s.ItemLink.Link),
		Kiosk:     *proto.Bool(s.Kiosk),
		Guild:     *proto.String(s.Guild.Name),
		Buyer:     *proto.String(s.Buyer.Name),
		Id:        *proto.String(s.Id),
		Quantity:  *proto.Uint64(uint64(s.Quantity)),
		Timestamp: *proto.Int64(s.Timestamp.Unix()),
		Price:     *proto.Uint64(uint64(s.Price)),
	}
}
