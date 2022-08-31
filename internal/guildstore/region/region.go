package region

import (
	"fmt"
	"log"
	"time"

	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data"
	accountregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/account"
	guildregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/guild"
	itemlinkregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/itemlink"
	luaconv "github.com/TechyShishy/nirn-revenue-service/internal/lua/conv"
	lua "github.com/yuin/gopher-lua"
)

type Name string

const (
	UNKNOWN Name = "Unknown"
	NA      Name = "NA"
	EU      Name = "EU"
)

type Region struct {
	ItemVariants     []data.ItemVariant
	ItemLinkRegistry *itemlinkregistry.ItemLinkRegistry
	GuildRegistry    *guildregistry.GuildRegistry
	AccountRegistry  *accountregistry.AccountRegistry
}

func NewFromLT(regionLT *lua.LTable) *Region {
	r := &Region{
		ItemLinkRegistry: itemlinkregistry.New(),
	}
	r.ItemVariants = r.parseRegion(regionLT)
	return r
}

func (r *Region) AddVariantsFromLT(regionLT *lua.LTable) *Region {
	return r.AddVariants(r.parseRegion(regionLT))
}

func (r *Region) AddVariants(variants []data.ItemVariant) *Region {
	r.ItemVariants = append(r.ItemVariants, variants...)
	return r
}

func (r *Region) parseRegion(regionLT *lua.LTable) []data.ItemVariant {
	region := []data.ItemVariant{}
	err := regionLT.ForEachWithError(func(idLV, variantLV lua.LValue) error {
		id, err := luaconv.Int(idLV)
		if err != nil {
			return err
		}
		variantLT, err := luaconv.Table(variantLV)
		if err != nil {
			return err
		}
		err = variantLT.ForEachWithError(func(vIdLV, listingLV lua.LValue) error {
			vId, err := luaconv.String(vIdLV)
			if err != nil {
				return err
			}
			listingLT, err := luaconv.Table(listingLV)
			if err != nil {
				return err
			}
			listing, err := r.parseListing(id, vId, listingLT)
			if err != nil {
				return err
			}
			region = append(region, listing)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Print(err)
		return nil
	}
	return region
}

func (r *Region) parseListing(id int, vId string, listingLT *lua.LTable) (data.ItemVariant, error) {
	listing := data.ItemVariant{Id: id, Variant: vId}
	err := listingLT.ForEachWithError(func(propertyLV, valueLV lua.LValue) error {
		property, err := luaconv.String(propertyLV)
		if err != nil {
			return fmt.Errorf("could not parse ItemVariant table key: %w", err)
		}
		switch property {
		case "itemAdderText":
			itemAdderText, err := luaconv.String(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.ItemAddrText: %w", err)
			}
			listing.ItemAdderText = itemAdderText
		case "totalCount":
			totalCount, err := luaconv.Int(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.TotalCount: %w", err)
			}
			listing.TotalCount = uint(totalCount)
		case "itemIcon":
			icon, err := luaconv.String(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.Icon: %w", err)
			}
			listing.Icon = icon
		case "newestTime":
			newestTime, err := luaconv.Int(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.NewestTime: %w", err)
			}
			listing.NewestTime = time.Unix(int64(newestTime), 0)
		case "oldestTime":
			oldestTime, err := luaconv.Int(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.OldestTime: %w", err)
			}
			listing.OldestTime = time.Unix(int64(oldestTime), 0)
		case "wasAltered":
			altered, err := luaconv.Bool(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.Altered: %w", err)
			}
			listing.Altered = altered
		case "sales":
			salesLT, err := luaconv.Table(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.Sales: %w", err)
			}
			sales, err := r.parseSales(salesLT)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.Sales: %w", err)
			}
			listing.Sales = sales
		case "itemDesc":
			description, err := luaconv.String(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse ItemVariant.Description: %w", err)
			}
			listing.Description = description
		}
		return nil
	})
	if err != nil {
		return data.ItemVariant{}, err // Passthrough parse error without comment
	}
	return listing, nil
}

func (r *Region) parseSales(salesLT *lua.LTable) ([]data.Sale, error) {
	sales := []data.Sale{}
	err := salesLT.ForEachWithError(func(_, saleLV lua.LValue) error {
		saleLT, err := luaconv.Table(saleLV)
		if err != nil {
			return fmt.Errorf("could not parse sale table: %w", err)
		}
		sale, err := r.parseSale(saleLT)
		if err != nil {
			return err // Passthrough parse error without comment
		}
		sales = append(sales, sale)
		return nil
	})
	if err != nil {
		return nil, err // Passthrough array error without comment
	}
	return sales, nil
}

func (r *Region) parseSale(saleLT *lua.LTable) (data.Sale, error) {
	sale := data.Sale{}
	err := saleLT.ForEachWithError(func(propertyLV, valueLV lua.LValue) error {
		property, err := luaconv.String(propertyLV)
		if err != nil {
			return fmt.Errorf("could not parse sale table key: %w", err)
		}
		switch property {
		case "wasKiosk":
			kiosk, err := luaconv.Bool(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse sale.Kiosk: %w", err)
			}
			sale.Kiosk = kiosk
		case "buyer":
			buyerId, err := luaconv.Uint(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse sale.BuyerId: %w", err)
			}
			sale.BuyerId = buyerId
			buyer, err := r.AccountRegistry.Add(buyerId, "")
			if err != nil {
				return fmt.Errorf(
					"could not add sale.BuyerId to region.AccountRegistry: %w",
					err,
				)
			}
			sale.Buyer = buyer
		case "price":
			price, err := luaconv.Int(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse sale.Price: %w", err)
			}
			sale.Price = price
		case "itemLink":
			itemLinkId, err := luaconv.Uint(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse sale.ItemLinkId: %w", err)
			}
			sale.ItemLinkId = itemLinkId
			itemLink, err := r.ItemLinkRegistry.Add(itemLinkId, "")
			if err != nil {
				return fmt.Errorf(
					"could not add sale.ItemLinkId to region.ItemLinkRegistry: %w",
					err,
				)
			}
			sale.ItemLink = itemLink
		case "seller":
			sellerId, err := luaconv.Uint(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse sale.SellerId: %w", err)
			}
			sale.SellerId = sellerId
			seller, err := r.AccountRegistry.Add(sellerId, "")
			if err != nil {
				return fmt.Errorf(
					"could not add sale.SellerId to region.AccountRegistry: %w",
					err,
				)
			}
			sale.Seller = seller
		case "timestamp":
			timestamp, err := luaconv.Int(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse sale.Timestamp: %w", err)
			}
			sale.Timestamp = time.Unix(int64(timestamp), 0)
		case "quant":
			quantity, err := luaconv.Int(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse sale.Quantity: %w", err)
			}
			sale.Quantity = uint(quantity)
		case "guild":
			guildId, err := luaconv.Uint(valueLV)
			if err != nil {
				return fmt.Errorf("could not parse sale.GuildId: %w", err)
			}
			sale.GuildId = guildId
			guild, err := r.GuildRegistry.Add(guildId, "")
			if err != nil {
				return fmt.Errorf(
					"could not add sale.GuildId to region.GuildRegistry: %w",
					err,
				)
			}
			sale.Guild = guild
		case "id":
			id, err := luaconv.String(valueLV)
			if err != nil {
				return fmt.Errorf("could not get sale.Id: %w", err)
			}
			sale.Id = id
		}
		return nil
	})
	if err != nil {
		return data.Sale{}, fmt.Errorf("could not parse sale: %w", err)
	}
	return sale, nil
}
