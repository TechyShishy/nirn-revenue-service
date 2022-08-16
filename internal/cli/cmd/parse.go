package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/sys/windows"

	"github.com/TechyShishy/nirn-revenue-service/internal/parser"
)

var parseFilePath *string

const GS00DataPath string = "Elder Scrolls Online/live/SavedVariables/GS00Data.lua"

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a single GSXXData file, checking for problems",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		globalVar := strings.TrimSuffix(
			filepath.Base(*parseFilePath),
			filepath.Ext(*parseFilePath),
		) + "SavedVariables"
		regionsData, err := parser.Parse(*parseFilePath, globalVar)
		if err != nil {
			log.Print(err)
			return
		}
		log.Printf("%#v", regionsData)
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	documentsPath, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		log.Print(fmt.Errorf("couldn't get Documents path: %w", err))
		documentsPath = "%USERPROFILE%/Documents"
	}
	parseFilePath = parseCmd.Flags().
		StringP("path", "p", filepath.Join(documentsPath, GS00DataPath), "Path to GSXXData file to parse")
}
