package itemlink

import (
	"fmt"
)

type ItemLinkRegistry struct {
	links map[uint]*ItemLink
}

func New() *ItemLinkRegistry {
	return &ItemLinkRegistry{links: make(map[uint]*ItemLink)}
}

func (i *ItemLinkRegistry) Add(id uint, link string) (*ItemLink, error) {
	if i.links[id] != nil {
		if link != "" && i.links[id].Link != link {
			if i.links[id].Link != "" {
				return nil, fmt.Errorf("cowardly refusing to update conflicting link ids")
			}
			i.links[id].Link = link
		}
		return i.links[id], nil
	}
	itemLink := &ItemLink{Id: id, Link: link}
	i.links[id] = itemLink
	return itemLink, nil
}

type ItemLink struct {
	Id   uint
	Link string
}
