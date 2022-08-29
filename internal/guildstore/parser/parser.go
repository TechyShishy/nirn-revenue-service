package parser

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore"
	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data"
	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore/region"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/sys/windows"
)

const (
	defaultSavedVariablesPathBase string = "Elder Scrolls Online/live/SavedVariables"
	defaultGSDataFileGlobBase     string = "GS[01][0-9]Data.lua"
)

var DefaultGSDataFileGlob string

type Parser struct {
	GSDataFileGlob string
	GSDataFiles    []guildstore.GSDataFile
}

func New() Parser {
	return Parser{
		GSDataFileGlob: DefaultGSDataFileGlob,
	}
}

func (p *Parser) ParseGlob() (map[region.Region][]data.ItemVariant, error) {
	gsDataFiles, err := p.globFiles(p.GSDataFileGlob)
	if err != nil {
		return nil, err
	}
	p.GSDataFiles = gsDataFiles
	return p.ParseAll()
}

func (p *Parser) ParseAll() (map[region.Region][]data.ItemVariant, error) {
	L := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer L.Close()

	globalLValues, err := readFiles(L, p.GSDataFiles)
	if err != nil {
		return nil, err
	}

	regionsData, err := parseGlobals(globalLValues)
	if err != nil {
		return nil, err
	}

	return regionsData, nil
}

func readFiles(l *lua.LState, files []guildstore.GSDataFile) (r []lua.LTable, err error) {
	for _, file := range files {
		if err := l.DoFile(file.Path); err != nil {
			return nil, err
		}
		lv, ok := (l.GetGlobal(file.GlobalVar)).(*lua.LTable)
		if !ok {
			return nil, fmt.Errorf("parsed file did not result in valid data")
		}
		r = append(r, *lv)
	}

	return r, nil
}

func (*Parser) globFiles(glob string) ([]guildstore.GSDataFile, error) {
	paths, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}
	var files []guildstore.GSDataFile
	for _, path := range paths {
		files = append(files, guildstore.GSDataFile{
			Path:      path,
			GlobalVar: GlobalVarFromPath(path),
		})
	}
	return files, nil
}

func GlobalVarFromPath(path string) string {
	return strings.TrimSuffix(
		filepath.Base(path),
		filepath.Ext(path),
	) + "SavedVariables"
}

func parseGlobals(globals []lua.LTable) (map[region.Region][]data.ItemVariant, error) {
	regionsData := make(map[region.Region][]data.ItemVariant)
	for _, global := range globals {
		r := make(map[region.Region][]data.ItemVariant)
		global.ForEach(func(key, value lua.LValue) {
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
				r[region.EU] = parseRegion(valueTable)
			case "datana":
				r[region.NA] = parseRegion(valueTable)
			}
		})
		regionsData = mergeRegions(regionsData, r)
	}
	return regionsData, nil
}

func mergeRegions(
	r ...map[region.Region][]data.ItemVariant,
) map[region.Region][]data.ItemVariant {
	allFiles := make(map[region.Region][]data.ItemVariant)
	for _, r2 := range r {
		for k, v := range r2 {
			allFiles[k] = append(allFiles[k], v...)
		}
	}
	return allFiles
}

func parseRegion(dataTable *lua.LTable) []data.ItemVariant {
	regionData := []data.ItemVariant{}
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

			i := data.ItemVariant{Id: idInt, Variant: vIdString}
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

func parseSales(t *lua.LTable) ([]data.Sale, error) {
	sales := []data.Sale{}
	t.ForEach(func(i, sale lua.LValue) {
		saleTable, err := luaTable(sale)
		if err != nil {
			log.Print(err)
			return
		}
		s := data.Sale{}
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
				s.BuyerId = valueInt
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
				s.ItemLinkId = valueInt
			case "seller":
				valueInt, err := luaInt(value)
				if err != nil {
					log.Print(err)
					return
				}
				s.SellerId = valueInt
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
				s.GuildId = valueInt
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

func init() {
	documentsPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		log.Print(err)
		return
	}
	DefaultGSDataFileGlob = filepath.Join(
		documentsPath,
		defaultSavedVariablesPathBase,
		defaultGSDataFileGlobBase,
	)
}
