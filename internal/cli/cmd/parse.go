package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore"
	"github.com/TechyShishy/nirn-revenue-service/internal/guildstore/parser"
)

var parseFilePath *string

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a single GSXXData file, checking for problems",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		p := parser.Parser{
			GSDataFiles: []guildstore.GSDataFile{
				{
					Path:      *parseFilePath,
					GlobalVar: parser.GlobalVarFromPath(*parseFilePath),
				},
			},
		}
		regionsData, err := p.ParseAll()
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
	parseFilePath = parseCmd.Flags().
		StringP("path", "p", parser.DefaultGSDataFileGlob, "Path to GSXXData file to parse")
}
