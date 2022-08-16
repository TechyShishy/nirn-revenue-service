package parser

import (
	"fmt"
	"log"
	"time"

	"github.com/TechyShishy/nirn-revenue-service/internal/gamedata"
	"github.com/TechyShishy/nirn-revenue-service/internal/gamedata/region"
	lua "github.com/yuin/gopher-lua"
)

func Parse(luaFile string, globalVar string) (map[region.Region][]gamedata.ItemVariant, error) {
	regionsData := make(map[region.Region][]gamedata.ItemVariant)
	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer L.Close()
	if err := L.DoFile(luaFile); err != nil {
		return nil, err
	}

	SavedVariablesLValue, ok := (L.GetGlobal(globalVar)).(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("parsed file did not result in valid data")
	}
	SavedVariablesLValue.ForEach(func(key, value lua.LValue) {
		keyString, err := luaString(key)
		if err != nil {
			log.Print(err)
			return
		}
		valueTable, err := luaTable(value)
		if err != nil {
			log.Print(err)
			return
		}

		switch keyString {
		case "dataeu":
			regionsData[region.EU] = parseRegion(valueTable)
		case "datana":
			regionsData[region.NA] = parseRegion(valueTable)
		}
	})

	return regionsData, nil
}

func parseRegion(dataTable *lua.LTable) []gamedata.ItemVariant {
	regionData := []gamedata.ItemVariant{}
	dataTable.ForEach(func(id, variant lua.LValue) {
		idInt, err := luaInt(id)
		if err != nil {
			log.Print(err)
			return
		}
		variantTable, err := luaTable(variant)
		if err != nil {
			log.Print(err)
			return
		}
		variantTable.ForEach(func(vId, listing lua.LValue) {
			vIdString, err := luaString(vId)
			if err != nil {
				log.Print(err)
				return
			}
			listingTable, err := luaTable(listing)
			if err != nil {
				log.Print(err)
				return
			}

			i := gamedata.ItemVariant{Id: idInt, Variant: vIdString}
			listingTable.ForEach(func(key, value lua.LValue) {
				keyString, err := luaString(key)
				if err != nil {
					log.Print(err)
					return
				}
				switch keyString {
				case "itemAdderText":
					valueString, err := luaString(value)
					if err != nil {
						log.Print(err)
						return
					}
					i.ItemAdderText = valueString
				case "totalCount":
					valueInt, err := luaInt(value)
					if err != nil {
						log.Print(err)
						return
					}
					i.TotalCount = uint(valueInt)
				case "itemIcon":
					valueString, err := luaString(value)
					if err != nil {
						log.Print(err)
						return
					}
					i.Icon = valueString
				case "newestTime":
					valueInt, err := luaInt(value)
					if err != nil {
						log.Print(err)
						return
					}
					i.NewestTime = time.Unix(int64(valueInt), 0)
				case "oldestTime":
					valueInt, err := luaInt(value)
					if err != nil {
						log.Print(err)
						return
					}
					i.OldestTime = time.Unix(int64(valueInt), 0)
				case "wasAltered":
					valueBool, err := luaBool(value)
					if err != nil {
						log.Print(err)
						return
					}
					i.Altered = valueBool
				case "sales":
					valueTable, err := luaTable(value)
					if err != nil {
						log.Print(err)
						return
					}
					sales, err := parseSales(valueTable)
					if err != nil {
						log.Print(err)
					}
					i.Sales = sales
				case "itemDesc":
					valueString, err := luaString(value)
					if err != nil {
						log.Print(err)
						return
					}
					i.Description = valueString
				}
			})
			regionData = append(regionData, i)
		})
	})
	return regionData
}

func parseSales(t *lua.LTable) ([]gamedata.Sale, error) {
	sales := []gamedata.Sale{}
	t.ForEach(func(i, sale lua.LValue) {
		saleTable, err := luaTable(sale)
		if err != nil {
			log.Print(err)
			return
		}
		s := gamedata.Sale{}
		saleTable.ForEach(func(key, value lua.LValue) {
			keyString, err := luaString(key)
			if err != nil {
				log.Print(err)
				return
			}
			switch keyString {
			case "wasKiosk":
				valueBool, err := luaBool(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.Kiosk = valueBool
			case "buyer":
				valueInt, err := luaInt(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.Buyer = valueInt
			case "price":
				valueInt, err := luaInt(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.Price = valueInt
			case "itemLink":
				valueInt, err := luaInt(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.ItemLink = valueInt
			case "seller":
				valueInt, err := luaInt(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.Seller = valueInt
			case "timestamp":
				valueInt, err := luaInt(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.Timestamp = time.Unix(int64(valueInt), 0)
			case "quant":
				valueInt, err := luaInt(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.Quantity = uint(valueInt)
			case "guild":
				valueInt, err := luaInt(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.Guild = valueInt
			case "id":
				valueString, err := luaString(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.Id = valueString
			}
		})
		sales = append(sales, s)
	})
	return sales, nil
}

func luaString(l lua.LValue) (string, error) {
	lString, ok := l.(lua.LString)
	if !ok {
		return "", fmt.Errorf("wanted string, got %v from lua: %#v", l.Type(), l)
	}
	return string(lString), nil
}

func luaInt(l lua.LValue) (int, error) {
	lNumber, ok := l.(lua.LNumber)
	if !ok {
		return 0, fmt.Errorf("wanted number, got %v from lua: %#v", l.Type(), l)
	}
	lFloat := float64(lNumber)
	lInt := int(lFloat)
	if lFloat != float64(lInt) {
		log.Print(fmt.Errorf("wanted int, got float from lua: %#v", lFloat))
	}
	return lInt, nil
}

func luaBool(l lua.LValue) (bool, error) {
	lBool, ok := l.(lua.LBool)
	if !ok {
		return false, fmt.Errorf("wanted bool, got %v from lua: %#v", l.Type(), l)
	}
	return bool(lBool), nil
}

func luaTable(l lua.LValue) (*lua.LTable, error) {
	lTable, ok := l.(*lua.LTable)
	if !ok {
		return &lua.LTable{}, fmt.Errorf("wanted lua.LTable, got %v from lua: %#v", l.Type(), l)
	}
	return lTable, nil
}
