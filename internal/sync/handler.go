package sync

import (
	"log"
	"net/rpc"

	pb "github.com/techyshishy/nirn-revenue-service/gen/api/proto/sync/service/v1"
)

type Handler struct{}

func (*Handler) UpdateRegions(req *pb.UpdateRegionsRequest, resp *pb.UpdateRegionsResponse) error {
	log.Print(req.Foo)
	resp.Bar = "test"
	// return fmt.Errorf("foobar")
	return nil
}

func Register() *rpc.Server {
	handler := &Handler{}
	server := rpc.NewServer()
	err := server.RegisterName("SyncService", handler)
	if err != nil {
		log.Print(err)
		return nil
	}
	return server
}
