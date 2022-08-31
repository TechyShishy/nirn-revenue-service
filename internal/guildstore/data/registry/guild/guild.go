package guild

import (
	"fmt"
)

type GuildRegistry struct {
	guilds map[uint]*Guild
}

func New() *GuildRegistry {
	return &GuildRegistry{guilds: make(map[uint]*Guild)}
}

func (g *GuildRegistry) Add(id uint, name string) (*Guild, error) {
	if g.guilds[id] != nil {
		if name != "" && g.guilds[id].Name != name {
			if g.guilds[id].Name != "" {
				return nil, fmt.Errorf("cowardly refusing to update conflicting guild ids")
			}
			g.guilds[id].Name = name
		}
		return g.guilds[id], nil
	}
	guild := &Guild{Id: id, Name: name}
	g.guilds[id] = guild
	return guild, nil
}

type Guild struct {
	Id   uint
	Name string
}
