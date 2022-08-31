package parser

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore"
	accountregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/account"
	guildregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/guild"
	itemlinkregistry "github.com/TechyShishy/nirn-revenue-service/internal/guildstore/data/registry/itemlink"
	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore/region"
	luaconv "github.com/TechyShishy/nirn-revenue-service/internal/lua/conv"
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

func (p *Parser) ParseGlob() (map[region.Name]*region.Region, error) {
	gsDataFiles, err := p.globFiles(p.GSDataFileGlob)
	if err != nil {
		return nil, err
	}
	p.GSDataFiles = gsDataFiles
	return p.ParseAll()
}

func (p *Parser) ParseAll() (map[region.Name]*region.Region, error) {
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

func parseGlobals(globals []lua.LTable) (map[region.Name]*region.Region, error) {
	regionsData := make(map[region.Name]*region.Region)
	itemLinks := itemlinkregistry.New()
	accounts := accountregistry.New()
	guilds := guildregistry.New()
	for _, global := range globals {
		err := global.ForEachWithError(func(keyLV, sectionLV lua.LValue) error {
			sectionKey, err := luaconv.String(keyLV)
			if err != nil {
				return err
			}
			sectionLT, err := luaconv.Table(sectionLV)
			if err != nil {
				return err
			}

			switch sectionKey {
			case "dataeu":

				regionsData[region.EU] = parseRegion(
					regionsData[region.EU],
					itemLinks,
					accounts,
					guilds,
					sectionLT,
				)
			case "datana":
				regionsData[region.NA] = parseRegion(
					regionsData[region.NA],
					itemLinks,
					accounts,
					guilds,
					sectionLT,
				)
			case "itemLink":
				itemLinks = parseItemLinks(itemLinks, sectionLT)
			case "accountNames":
				accounts = parseAccounts(accounts, sectionLT)
			case "guildNames":
				guilds = parseGuilds(guilds, sectionLT)
			default:
				return nil // Not one of the data sections we care about right now
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return regionsData, nil
}

func parseRegion(
	regionData *region.Region,
	itemlinks *itemlinkregistry.ItemLinkRegistry,
	accounts *accountregistry.AccountRegistry,
	guilds *guildregistry.GuildRegistry,
	regionLT *lua.LTable,
) *region.Region {
	if regionData == nil {
		regionData = &region.Region{
			ItemLinkRegistry: itemlinks,
			AccountRegistry:  accounts,
			GuildRegistry:    guilds,
		}
	}
	return regionData.AddVariantsFromLT(regionLT)
}

func parseItemLinks(
	itemLinks *itemlinkregistry.ItemLinkRegistry,
	itemLinksLT *lua.LTable,
) *itemlinkregistry.ItemLinkRegistry {
	err := itemLinksLT.ForEachWithError(func(linkLV, idLV lua.LValue) error {
		link, err := luaconv.String(linkLV)
		if err != nil {
			return err
		}
		id, err := luaconv.Uint(idLV)
		if err != nil {
			return err
		}
		itemLinks.Add(id, link)
		return nil
	})
	if err != nil {
		log.Print(err)
		return nil
	}
	return itemLinks
}

func parseAccounts(
	accounts *accountregistry.AccountRegistry,
	accountsLT *lua.LTable,
) *accountregistry.AccountRegistry {
	err := accountsLT.ForEachWithError(func(nameLV, idLV lua.LValue) error {
		name, err := luaconv.String(nameLV)
		if err != nil {
			return err
		}
		id, err := luaconv.Uint(idLV)
		if err != nil {
			return err
		}
		accounts.Add(id, name)
		return nil
	})
	if err != nil {
		log.Print(err)
		return nil
	}
	return accounts
}

func parseGuilds(
	guilds *guildregistry.GuildRegistry,
	guildsLT *lua.LTable,
) *guildregistry.GuildRegistry {
	err := guildsLT.ForEachWithError(func(nameLV, idLV lua.LValue) error {
		name, err := luaconv.String(nameLV)
		if err != nil {
			return err
		}
		id, err := luaconv.Uint(idLV)
		if err != nil {
			return err
		}
		guilds.Add(id, name)
		return nil
	})
	if err != nil {
		log.Print(err)
		return nil
	}
	return guilds
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
