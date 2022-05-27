package rpc

import (
	"context"

	"github.com/clubo-app/packages/utils"
	cg "github.com/clubo-app/protobuf/common"
	pg "github.com/clubo-app/protobuf/party"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s partyServer) DeleteParty(c context.Context, req *pg.DeletePartyRequest) (*cg.SuccessIndicator, error) {
	_, err := ksuid.Parse(req.PartyId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Party id")
	}

	_, err = ksuid.Parse(req.RequesterId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Requester id")
	}

	_, err = s.ps.Delete(c, req.RequesterId, req.PartyId)
	if err != nil {
		return nil, utils.HandleError(err)
	}

	return &cg.SuccessIndicator{Sucess: true}, nil
}
