package main

import (
	"fmt"

	"golang.org/x/net/context"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

//ClientUpdate forces a move
func (s *Server) ClientUpdate(ctx context.Context, in *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	version := (in.ProtoReflect().Descriptor().Fields().Len())
	config, err := s.read(ctx)
	if err != nil {
		return nil, err
	}
	s.Log(fmt.Sprintf("Updating %v with %v and %v", in.GetInstanceId(), version, config.Tracked[in.GetInstanceId()]))
	return &rcpb.ClientUpdateResponse{}, s.save(ctx, config)
}
