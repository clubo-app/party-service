package service

import (
	"context"

	"github.com/clubo-app/party-service/dto"
	"github.com/clubo-app/party-service/repository"
	"github.com/segmentio/ksuid"
)

type PartyService interface {
	Create(ctx context.Context, p dto.Party) (repository.Party, error)
	Update(ctx context.Context, p dto.Party) (repository.Party, error)
	Delete(ctx context.Context, uId, pId string) error
	Get(ctx context.Context, pId string) (repository.Party, error)
	GetMany(ctx context.Context, ids []string) ([]repository.Party, error)
	GetByUser(ctx context.Context, uId string, limit, offset uint64) ([]repository.Party, error)
}

type partyService struct {
	r *repository.PartyRepository
}

func NewPartyService(r *repository.PartyRepository) PartyService {
	return &partyService{r: r}
}

func (s partyService) Create(ctx context.Context, p dto.Party) (res repository.Party, err error) {
	res, err = s.r.CreateParty(ctx, repository.CreatePartyParams{
		ID:            ksuid.New().String(),
		UserID:        p.UserId,
		Title:         p.Title,
		IsPublic:      p.IsPublic,
		Lat:           p.Lat,
		Long:          p.Long,
		StreetAddress: p.StreetAddress,
		PostalCode:    p.PostalCode,
		State:         p.State,
		Country:       p.Country,
		StartDate:     p.StartDate,
		EndDate:       p.EndDate,
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) Update(ctx context.Context, p dto.Party) (res repository.Party, err error) {
	res, err = s.r.UpdateParty(ctx, repository.UpdatePartyParams{
		ID:            p.ID,
		Title:         p.Title,
		Lat:           p.Lat,
		Long:          p.Long,
		StreetAddress: p.StreetAddress,
		PostalCode:    p.PostalCode,
		State:         p.State,
		Country:       p.Country,
		StartDate:     p.StartDate,
		EndDate:       p.EndDate,
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) Delete(ctx context.Context, uId, pId string) error {
	err := s.r.DeleteParty(ctx, repository.DeletePartyParams{
		ID:     pId,
		UserID: uId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s partyService) Get(ctx context.Context, pId string) (res repository.Party, err error) {
	res, err = s.r.GetParty(ctx, pId)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) GetMany(ctx context.Context, ids []string) (res []repository.Party, err error) {
	res, err = s.r.GetManyParties(ctx, repository.GetManyPartiesParams{
		Ids:   ids,
		Limit: uint64(len(ids)),
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) GetByUser(ctx context.Context, uId string, limit, offset uint64) (res []repository.Party, err error) {
	res, err = s.r.GetPartiesByUser(ctx, repository.GetPartiesByUserParams{
		UserID: uId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return res, err
	}

	return res, nil
}
