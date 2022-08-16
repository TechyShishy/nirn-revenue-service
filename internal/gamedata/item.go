package gamedata

import "time"

type ItemVariant struct {
	Id            int
	Variant       string
	Description   string
	OldestTime    time.Time
	NewestTime    time.Time
	TotalCount    uint
	Altered       bool
	ItemAdderText string // TODO(TechyShishy): fix type name
	Icon          string
	Sales         []Sale
}
