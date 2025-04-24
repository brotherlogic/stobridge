package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbg "github.com/brotherlogic/goserver/proto"
	kmpb "github.com/brotherlogic/keymapper/proto"
	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/stobridge/proto"
)

const (
	// KEY used to save scores
	KEY = "github.com/brotherlogic/stobridge/config"
)

// Server main server type
type Server struct {
	*goserver.GoServer
	key string
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
	}
	return s
}

func (s *Server) save(ctx context.Context, config *pb.Config) error {
	return s.KSclient.Save(ctx, KEY, config)
}

func (s *Server) read(ctx context.Context) (*pb.Config, error) {
	scores := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, scores)

	if err != nil {
		return nil, err
	}

	return data.(*pb.Config), nil
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	rcpb.RegisterClientUpdateServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{}
}

func (s *Server) getRecord(ctx context.Context, instanceID int32) (*rcpb.Record, error) {
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rcpb.NewRecordCollectionServiceClient(conn)
	resp, err := client.GetRecord(ctx, &rcpb.GetRecordRequest{InstanceId: instanceID})
	if err != nil {
		return nil, err
	}

	return resp.GetRecord(), nil
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer("stobridge")
	server.Register = server

	ctx, cancel := utils.ManualContext("ghc", time.Minute)
	conn, err := server.FDialServer(ctx, "keymapper")
	if err != nil {
		if status.Convert(err).Code() == codes.Unknown {
			server.CtxLog(ctx, fmt.Sprintf("Cannot reach keymapper: %v", err))
		}
		return
	}
	client := kmpb.NewKeymapperServiceClient(conn)
	resp, err := client.Get(ctx, &kmpb.GetRequest{Key: "stobridge_user_id"})
	if err != nil {
		if status.Convert(err).Code() == codes.Unknown || status.Convert(err).Code() == codes.InvalidArgument {
			server.CtxLog(ctx, fmt.Sprintf("Cannot read external: %v", err))
		}
		return
	}
	conn.Close()
	server.key = resp.GetKey().GetValue()
	cancel()

	err = server.RegisterServerV2(false)
	if err != nil {
		return
	}

	fmt.Printf("%v", server.Serve())
}
