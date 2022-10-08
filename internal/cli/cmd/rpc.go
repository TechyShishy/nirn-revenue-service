package cmd

import (
	"log"
	"net"
	"net/rpc"

	"github.com/spf13/cobra"
	pbss "github.com/techyshishy/nirn-revenue-service/gen/api/proto/sync/service/v1"
	"google.golang.org/protobuf/proto"
)

// parseCmd represents the parse command
var rpcCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Send an rpc command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := rpc.DialHTTP("tcp6", net.JoinHostPort("::1", "8080"))
		if err != nil {
			log.Print(err)
		}
		defer c.Close()
		req := &pbss.UpdateRegionsRequest{}

		req.Foo = *proto.String("test1")
		resp := &pbss.UpdateRegionsResponse{}
		c.Call("SyncService.UpdateRegions", req, resp)

		log.Print(resp.String())
	},
}

func init() {
	rootCmd.AddCommand(rpcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// parseCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
