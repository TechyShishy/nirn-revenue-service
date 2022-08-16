package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/TechyShishy/nirn-revenue-service/internal/gamedata"
	"github.com/TechyShishy/nirn-revenue-service/internal/gamedata/region"
	"github.com/TechyShishy/nirn-revenue-service/internal/parser"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows"
)

const (
	SavedVariablesPath string = "Elder Scrolls Online/live/SavedVariables"
	GSXXDataGlobBase   string = "GS??Data.lua"
)

var fileGlob *string

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := filepath.Glob(*fileGlob)
		if err != nil {
			log.Print(err)
			return
		}
		regionsData := make(map[region.Region][]gamedata.ItemVariant)
		for _, file := range files {
			globalVar := strings.TrimSuffix(
				filepath.Base(file),
				filepath.Ext(file),
			) + "SavedVariables"
			log.Printf("Processing %v...", filepath.Base(file))
			r, err := parser.Parse(file, globalVar)
			if err != nil {
				log.Print(err)
				return
			}
			regionsData = mergeRegions(regionsData, r)
		}

		log.Printf("Found %v records", len(regionsData[region.NA]))
	},
}

func mergeRegions(
	r ...map[region.Region][]gamedata.ItemVariant,
) map[region.Region][]gamedata.ItemVariant {
	allFiles := make(map[region.Region][]gamedata.ItemVariant)
	for _, r2 := range r {
		for k, v := range r2 {
			allFiles[k] = append(allFiles[k], v...)
		}
	}
	return allFiles
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	documentsPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		log.Print(fmt.Errorf("couldn't get Documents path: %w", err))
		documentsPath = "%USERPROFILE%/Documents"
	}
	fileGlob = uploadCmd.Flags().
		StringP("glob", "g", filepath.Join(documentsPath, SavedVariablesPath, GSXXDataGlobBase), "glob path that matches GSXXData files to upload")
}
