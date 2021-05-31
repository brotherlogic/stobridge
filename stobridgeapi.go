package main

import (
	"fmt"

	"golang.org/x/net/context"

	rcpb "github.com/brotherlogic/recordcollection/proto"
)

//ClientUpdate forces a move
func (s *Server) ClientUpdate(ctx context.Context, in *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	s.Log(fmt.Sprintf("Updating %v", in.GetInstanceId()))
	return &rcpb.ClientUpdateResponse{}, nil
}
