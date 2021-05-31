package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rcpb "github.com/brotherlogic/recordcollection/proto"
	pb "github.com/brotherlogic/stobridge/proto"
	stopb "github.com/brotherlogic/straightenthemout-library/proto"
	stologicpb "github.com/brotherlogic/straightenthemout-logic/proto"
)

func (s *Server) setMetadata(iid int32, meta *stopb.Metadata) error {
	req := &stologicpb.SetMetadataRequest{
		Stoid:    s.key,
		Metadata: meta,
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://straightenthemout-qo2wxnmyfq-uw.a.run.app/straightenthemout.STOService/SetMetadata", "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

//ClientUpdate runs an update
func (s *Server) ClientUpdate(ctx context.Context, in *rcpb.ClientUpdateRequest) (*rcpb.ClientUpdateResponse, error) {
	version := int32(in.ProtoReflect().Descriptor().Fields().Len())
	config, err := s.read(ctx)
	if err != nil {
		if status.Convert(err).Code() != codes.InvalidArgument {
			return nil, err
		}
		config = &pb.Config{Tracked: make(map[int32]int32)}
	}

	if config.Tracked[in.GetInstanceId()] != version {
		record, err := s.getRecord(ctx, in.GetInstanceId())
		if err != nil {
			return nil, err
		}
		err = s.setMetadata(in.GetInstanceId(), &stopb.Metadata{
			InstanceId: in.GetInstanceId(),
			Width:      record.GetMetadata().GetRecordWidth(),
		})
		if err != nil {
			return nil, err
		}

		config.Tracked[in.GetInstanceId()] = version
	}

	s.Log(fmt.Sprintf("Completed Update %v with %v and then %v", in.GetInstanceId(), version, config.Tracked[in.GetInstanceId()]))
	return &rcpb.ClientUpdateResponse{}, s.save(ctx, config)
}
