package main

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/stobridge/proto"
)

//ClientUpdate runs an update
func (s *Server) ClientUpdate(ctx context.Context, in *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	version := (in.ProtoReflect().Descriptor().Fields().Len())
	config, err := s.read(ctx)
	if err != nil {
		if status.Convert(err).Code() != codes.InvalidArgument {
			return nil, err
		}
		config = &pb.Config{Tracked: make(map[int32]int32)}
	}
	s.Log(fmt.Sprintf("Updating %v with %v and then %v", in.GetInstanceId(), version, config.Tracked[in.GetInstanceId()]))
	return &rcpb.ClientUpdateResponse{}, s.save(ctx, config)
}
