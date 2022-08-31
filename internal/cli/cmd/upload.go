package cmd

import (
	"log"

	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore/parser"
	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore/region"
	"github.com/spf13/cobra"
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
		p := parser.Parser{GSDataFileGlob: *fileGlob}

		regionsData, err := p.ParseGlob()
		if err != nil {
			log.Print(err)
			return
		}

		log.Printf("Found %v records", len(regionsData[region.NA].ItemVariants))
		sale := regionsData[region.NA].ItemVariants[12].Sales[0]
		log.Printf(
			"ItemVariant 12 Sale 0: Seller (%#v) Buyer (%#v) Guild(%#v) Link (%#v)",
			sale.Seller.Name,
			sale.Buyer.Name,
			sale.Guild.Name,
			sale.ItemLink.Link,
		)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	fileGlob = uploadCmd.Flags().
		StringP("glob", "g", parser.DefaultGSDataFileGlob, "glob path that matches GSXXData files to upload")
}
