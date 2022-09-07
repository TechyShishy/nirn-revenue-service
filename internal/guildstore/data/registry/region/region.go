package region

import (
	accountregistry "github.com/techyshishy/nirn-revenue-service/internal/guildstore/data/registry/account"
	guildregistry "github.com/techyshishy/nirn-revenue-service/internal/guildstore/data/registry/guild"
	itemlinkregistry "github.com/techyshishy/nirn-revenue-service/internal/guildstore/data/registry/itemlink"
	"github.com/techyshishy/nirn-revenue-service/internal/guildstore/region"
)

type RegionRegistry struct {
	ItemLinkRegistry *itemlinkregistry.ItemLinkRegistry
	AccountRegistry *accountregistry.AccountRegistry
	GuildRegistry *guildregistry.GuildRegistry
	regionMap map[region.Name]*region.Region
}

func New() *RegionRegistry {
	return &RegionRegistry{
		ItemLinkRegistry: itemlinkregistry.New(),
		AccountRegistry: accountregistry.New(),
		GuildRegistry: guildregistry.New(),
		regionMap: make(map[region.Name]*region.Region),
	}
}

func (m *RegionRegistry) Region(name region.Name) *region.Region {
	if m.regionMap[name] == nil {
		m.regionMap[name] = &region.Region{
			Name: name,
			ItemLinkRegistry: m.ItemLinkRegistry,
			AccountRegistry:  m.AccountRegistry,
			GuildRegistry:    m.GuildRegistry,
		}
	}
	return m.regionMap[name]
}