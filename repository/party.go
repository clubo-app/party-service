package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/clubo-app/party-service/datastruct"
	"github.com/go-playground/validator/v10"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
	"github.com/scylladb/gocqlx/v2/table"
)

const (
	TABLE_NAME string = "parties"
)

var partyMetadata = table.Metadata{
	Name:    TABLE_NAME,
	Columns: []string{"id", "user_id", "title", "is_public", "lat", "long", "street_address", "postal_code", "state", "country", "start_date", "end_date"},
	PartKey: []string{"user_id"},
	SortKey: []string{"is_public", "start_date"},
}

type PartyRepository interface {
	Create(ctx context.Context, p datastruct.Party) (datastruct.Party, error)
	Update(ctx context.Context, p datastruct.Party) error
	Delete(ctx context.Context, uId, pId string) error
	Get(ctx context.Context, pId string) (datastruct.Party, error)
	GetMany(ctx context.Context, ids []string) ([]datastruct.Party, error)
	GetByUser(ctx context.Context, uId string, page []byte, limit uint32) ([]datastruct.Party, []byte, error)
}

type partyRepository struct {
	sess *gocqlx.Session
}

func (r *partyRepository) Create(ctx context.Context, p datastruct.Party) (datastruct.Party, error) {
	v := validator.New()
	err := v.Struct(p)
	if err != nil {
		return datastruct.Party{}, err
	}

	stmt, names := qb.
		Insert(TABLE_NAME).
		Columns(partyMetadata.Columns...).
		ToCql()

	log.Println(stmt)

	err = r.sess.
		Query(stmt, names).
		BindStruct(p).
		ExecRelease()
	if err != nil {
		return datastruct.Party{}, err
	}

	return p, nil
}

func (r *partyRepository) Get(ctx context.Context, pId string) (res datastruct.Party, err error) {
	stmt, names := qb.
		Select(TABLE_NAME).
		Columns(partyMetadata.Columns...).
		Where(qb.Eq("id")).
		ToCql()

	err = r.sess.
		Query(stmt, names).
		BindMap((qb.M{"id": pId})).
		GetRelease(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *partyRepository) Update(ctx context.Context, p datastruct.Party) error {
	b := qb.
		Update(TABLE_NAME).
		Where(qb.Eq("id"))

	if p.Title != "" {
		b.Set("title")
	}

	if p.Lat != 0 {
		b.Set("lat")
	}

	if p.Long != 0 {
		b.Set("long")
	}

	if !p.StartDate.IsZero() {
		b.Set("start_date")
	}

	if !p.EndDate.IsZero() {
		b.Set("end_date")
	}

	if p.StreetAddress != "" {
		b.Set("street_address")
	}

	if p.PostalCode != "" {
		b.Set("postal_code")
	}

	if p.State != "" {
		b.Set("state")
	}

	if p.Country != "" {
		b.Set("country")
	}

	b.If(qb.Eq("user_id"))
	stmt, names := b.ToCql()

	err := r.sess.Query(stmt, names).
		BindMap((qb.M{
			"id":             p.Id,
			"user_id":        p.UserId,
			"title":          p.Title,
			"lat":            p.Lat,
			"long":           p.Long,
			"street_address": p.StreetAddress,
			"postal_code":    p.PostalCode,
			"state":          p.State,
			"country":        p.Country,
			"start_date":     p.StartDate,
			"end_date":       p.EndDate,
		})).
		ExecRelease()
	if err != nil {
		return err
	}

	return nil
}

func (r *partyRepository) Delete(ctx context.Context, uId, pId string) error {
	stmt, names := qb.
		Delete(TABLE_NAME).
		Where(qb.Eq("id")).
		If(qb.Eq("user_id")).
		ToCql()

	err := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"id": pId, "user_id": uId})).
		ExecRelease()
	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			return errors.New("you can only Delete your own Parties")
		}
		return err
	}
	return nil
}

func (r *partyRepository) GetMany(ctx context.Context, ids []string) (res []datastruct.Party, err error) {
	stmt, names := qb.
		Select(TABLE_NAME).
		Where(qb.In("id")).
		ToCql()

	err = r.sess.Query(stmt, names).
		BindMap((qb.M{
			"id": ids,
		})).
		GetRelease(&res)
	if err != nil {
		return res, err
	}

	return res, nil

}

func (r *partyRepository) GetByUser(ctx context.Context, uId string, page []byte, limit uint32) (result []datastruct.Party, nextPage []byte, err error) {
	stmt, names := qb.
		Select(TABLE_NAME).
		Where(qb.Eq("user_id")).
		ToCql()

	q := r.sess.
		Query(stmt, names).
		BindMap((qb.M{"user_id": uId}))
	defer q.Release()

	q.PageState(page)
	if limit == 0 {
		q.PageSize(10)
	} else {
		q.PageSize(int(limit))
	}

	iter := q.Iter()
	err = iter.Select(&result)
	if err != nil {
		return []datastruct.Party{}, nil, errors.New("no parties found")
	}

	return result, iter.PageState(), nil
}
