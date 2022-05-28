package service

import (
	"context"
	"database/sql"

	"github.com/clubo-app/party-service/dto"
	"github.com/clubo-app/party-service/repository"
	"github.com/segmentio/ksuid"
)

type PartyService interface {
	Create(ctx context.Context, p dto.Party) (repository.Party, error)
	Update(ctx context.Context, p dto.Party) (repository.Party, error)
	Delete(ctx context.Context, uId, pId string) (repository.Party, error)
	Get(ctx context.Context, pId string) (repository.Party, error)
	GetMany(ctx context.Context, ids []string) ([]repository.Party, error)
	GetByUser(ctx context.Context, uId string, limit, offset int32) ([]repository.Party, error)
}

type partyService struct {
	q *repository.Queries
}

func NewPartyService(q *repository.Queries) PartyService {
	return &partyService{q: q}
}

func (s partyService) Create(ctx context.Context, p dto.Party) (res repository.Party, err error) {
	res, err = s.q.CreateParty(ctx, repository.CreatePartyParams{
		ID:            ksuid.New().String(),
		UserID:        p.UserId,
		Title:         p.Title,
		IsPublic:      p.IsPublic,
		Location:      p.Location,
		StreetAddress: sql.NullString{String: p.StreetAddress, Valid: p.StreetAddress != ""},
		PostalCode:    sql.NullString{String: p.PostalCode, Valid: p.PostalCode != ""},
		State:         sql.NullString{String: p.State, Valid: p.State != ""},
		Country:       sql.NullString{String: p.Country, Valid: p.Country != ""},
		StartDate:     sql.NullTime{Time: p.StartDate, Valid: !p.StartDate.IsZero()},
		EndDate:       sql.NullTime{Time: p.EndDate, Valid: !p.EndDate.IsZero()},
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) Update(ctx context.Context, p dto.Party) (res repository.Party, err error) {
	res, err = s.q.UpdateParty(ctx, repository.UpdatePartyParams{
		ID:            p.ID,
		Title:         sql.NullString{String: p.Title, Valid: p.Title != ""},
		StreetAddress: sql.NullString{String: p.StreetAddress, Valid: p.StreetAddress != ""},
		PostalCode:    sql.NullString{String: p.PostalCode, Valid: p.PostalCode != ""},
		State:         sql.NullString{String: p.State, Valid: p.State != ""},
		Country:       sql.NullString{String: p.Country, Valid: p.Country != ""},
		StartDate:     sql.NullTime{Time: p.StartDate, Valid: !p.StartDate.IsZero()},
		EndDate:       sql.NullTime{Time: p.EndDate, Valid: !p.EndDate.IsZero()},
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) Delete(ctx context.Context, uId, pId string) (res repository.Party, err error) {
	res, err = s.q.DeleteParty(ctx, repository.DeletePartyParams{
		ID:     pId,
		UserID: uId,
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) Get(ctx context.Context, pId string) (res repository.Party, err error) {
	res, err = s.q.GetParty(ctx, pId)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) GetMany(ctx context.Context, ids []string) (res []repository.Party, err error) {
	res, err = s.q.GetManyParties(ctx, repository.GetManyPartiesParams{
		Ids:   ids,
		Limit: int32(len(ids)),
	})
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s partyService) GetByUser(ctx context.Context, uId string, limit, offset int32) (res []repository.Party, err error) {
	res, err = s.q.GetPartiesByUser(ctx, repository.GetPartiesByUserParams{
		UserID: uId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return res, err
	}

	return res, nil
}
