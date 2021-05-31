package main

import (
	"fmt"

	"golang.org/x/net/context"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

//ClientUpdate forces a move
func (s *Server) ClientUpdate(ctx context.Context, in *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	version := (in.ProtoReflect().Descriptor().Fields().Len())
	s.Log(fmt.Sprintf("Updating %v with %v", in.GetInstanceId(), version))
	return &rcpb.ClientUpdateResponse{}, nil
}
