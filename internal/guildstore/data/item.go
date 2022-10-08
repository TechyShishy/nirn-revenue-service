package data

import (
	"time"

	pbiv "github.com/techyshishy/nirn-revenue-service/gen/api/proto/itemvariant/v1"
	pbs "github.com/techyshishy/nirn-revenue-service/gen/api/proto/sale/v1"
	"google.golang.org/protobuf/proto"
)

type ItemVariant struct {
	Id            uint
	Variant       string
	Description   string
	OldestTime    time.Time
	NewestTime    time.Time
	TotalCount    uint
	Altered       bool
	ItemAdderText string // TODO(TechyShishy): fix type name
	Icon          string
	Sales         []Sale
	Link          string
}

func (i *ItemVariant) Proto() *pbiv.ItemVariant {
	p := &pbiv.ItemVariant{
		Id:            *proto.Uint64(uint64(i.Id)),
		Variant:       *proto.String(i.Variant),
		Description:   *proto.String(i.Description),
		OldestTime:    *proto.Int64(i.OldestTime.Unix()),
		NewestTime:    *proto.Int64(i.NewestTime.Unix()),
		TotalCount:    *proto.Uint64(uint64(i.TotalCount)),
		Altered:       *proto.Bool(i.Altered),
		ItemAdderText: *proto.String(i.ItemAdderText),
		Icon:          *proto.String(i.Icon),
		Sale:          []*pbs.Sale{},
		Link:          *proto.String(i.Link),
	}
	for _, sale := range i.Sales {
		p.Sale = append(p.Sale, sale.Proto())
	}

	return p
}
