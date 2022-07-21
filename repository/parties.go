package repository

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/leporo/sqlf"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/paulmach/orb/encoding/wkb"
)

const (
	TABLE_NAME = "parties"
)

type PartyRepository struct {
	pool    *pgxpool.Pool
	querier Querier
}

func NewPartyRepository(dbUser, dbPW, dbName, dbHost string, dbPort uint16) (*PartyRepository, error) {
	urlStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPW, dbHost, fmt.Sprint(dbPort), dbName)
	pgURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	connURL := *pgURL
	if connURL.Scheme == "cockroachdb" {
		connURL.Scheme = "postgres"
	}

	c, err := pgxpool.ParseConfig(connURL.String())
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), c)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}

	err = validateSchema(connURL)
	if err != nil {
		log.Printf("Schema validation error: %v", err)
	}

	return &PartyRepository{
		pool:    pool,
		querier: New(pool),
	}, nil
}

func (r PartyRepository) Close() {
	r.pool.Close()
}

type CreatePartyParams struct {
	ID            string
	UserID        string
	Title         string
	IsPublic      bool
	Location      orb.Point
	StreetAddress string
	PostalCode    string
	State         string
	Country       string
	StartDate     time.Time
	EndDate       time.Time
}

const selectStmt = "id, user_id, title, is_public, ST_AsBinary(location) AS location, street_address, postal_code, state, country, start_date, end_date"

func (r PartyRepository) CreateParty(ctx context.Context, arg CreatePartyParams) (Party, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	b := sqlf.
		InsertInto(TABLE_NAME).
		Set("id", arg.ID).
		Set("user_id", arg.UserID).
		Set("title", arg.Title).
		Set("is_public", arg.IsPublic).
		SetExpr("location", "ST_GeomFromEWKB(?)", ewkb.Value(arg.Location, 4326)).
		Set("street_address", arg.StreetAddress).
		Set("postal_code", arg.PostalCode).
		Set("state", arg.State).
		Set("country", arg.Country).
		Set("start_date", arg.StartDate).
		Set("end_date", arg.EndDate).
		Returning(selectStmt)

	row := r.pool.QueryRow(ctx, b.String(), b.Args()...)
	var i Party
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.IsPublic,
		wkb.Scanner(&i.Location),
		&i.StreetAddress,
		&i.PostalCode,
		&i.State,
		&i.Country,
		&i.StartDate,
		&i.EndDate,
	)
	return i, err
}

type UpdatePartyParams struct {
	ID            string
	Title         string
	Location      orb.Point
	StreetAddress string
	PostalCode    string
	State         string
	Country       string
	StartDate     time.Time
	EndDate       time.Time
}

func (r PartyRepository) UpdateParty(ctx context.Context, arg UpdatePartyParams) (Party, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	b := sqlf.Update(TABLE_NAME)

	if arg.Title != "" {
		b = b.Set("title", arg.Title)
	}
	if arg.Location.Lat() != 0 && arg.Location.Lon() != 0 {
		b = b.SetExpr("location", "ST_GeomFromEWKB(?)", ewkb.Value(arg.Location, 4326))
	}
	if arg.StreetAddress != "" {
		b = b.Set("street_address", arg.StreetAddress)
	}
	if arg.PostalCode != "" {
		b = b.Set("postal_code", arg.PostalCode)
	}
	if arg.State != "" {
		b = b.Set("state", arg.State)
	}
	if arg.Country != "" {
		b = b.Set("country", arg.Country)
	}
	startYear := arg.StartDate.Year()
	if !(startYear == 1970) {
		b = b.Set("start_date", arg.StartDate)
	}
	endYear := arg.StartDate.Year()
	if !(endYear == 1970) {
		b = b.Set("end_date", arg.EndDate)
	}

	b.
		Where("id = ?", arg.ID).
		Returning("id, user_id, title, is_public, ST_AsBinary(location) AS location, street_address, postal_code, state, country, start_date, end_date")

	row := r.pool.QueryRow(ctx, b.String(), b.Args()...)
	var i Party
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.IsPublic,
		wkb.Scanner(&i.Location),
		&i.StreetAddress,
		&i.PostalCode,
		&i.State,
		&i.Country,
		&i.StartDate,
		&i.EndDate,
	)

	return i, err
}

func (r PartyRepository) GetParty(ctx context.Context, id string) (Party, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	b := sqlf.
		Select(selectStmt).
		From(TABLE_NAME).
		Where("id = ?", id)

	row := r.pool.QueryRow(ctx, b.String(), b.Args()...)
	var i Party
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.IsPublic,
		wkb.Scanner(&i.Location),
		&i.StreetAddress,
		&i.PostalCode,
		&i.State,
		&i.Country,
		&i.StartDate,
		&i.EndDate,
	)
	return i, err
}

type GetManyPartiesParams struct {
	Ids   []string
	Limit uint64
}

func (r PartyRepository) GetManyParties(ctx context.Context, arg GetManyPartiesParams) ([]Party, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	b := sqlf.
		Select(selectStmt).
		From(TABLE_NAME).
		Where("id = ANY(?)", arg.Ids).
		Limit(arg.Limit)

	rows, err := r.pool.Query(ctx, b.String(), b.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Party
	for rows.Next() {
		var i Party
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.IsPublic,
			wkb.Scanner(&i.Location),
			&i.StreetAddress,
			&i.PostalCode,
			&i.State,
			&i.Country,
			&i.StartDate,
			&i.EndDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type GetPartiesByUserParams struct {
	UserID string
	Limit  uint64
	Offset uint64
}

func (r PartyRepository) GetPartiesByUser(ctx context.Context, arg GetPartiesByUserParams) ([]Party, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	b := sqlf.
		Select(selectStmt).
		From(TABLE_NAME).
		Where("user_id = ?", arg.UserID).
		OrderBy("id desc")
	if arg.Limit == 0 {
		b = b.Limit(10)
	} else {
		b = b.Limit(arg.Limit)
	}
	b = b.Offset(arg.Offset)

	rows, err := r.pool.Query(ctx, b.String(), b.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Party
	for rows.Next() {
		var i Party
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.IsPublic,
			wkb.Scanner(&i.Location),
			&i.StreetAddress,
			&i.PostalCode,
			&i.State,
			&i.Country,
			&i.StartDate,
			&i.EndDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type GeoSearchParams struct {
	Lat    float32
	Long   float32
	Radius int32
	Limit  int32
	Offset int32
}

func (r PartyRepository) GeoSearch(ctx context.Context, arg GeoSearchParams) ([]Party, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	b := sqlf.
		Select(selectStmt).
		From(TABLE_NAME).
		Where("ST_DWithin(location,ST_MakePoint(?, ?)::geography, ?)", arg.Long, arg.Lat, arg.Radius)

	if arg.Limit == 0 {
		b = b.Limit(10)
	} else {
		b = b.Limit(arg.Limit)
	}
	b = b.Offset(arg.Offset)

	rows, err := r.pool.Query(ctx, b.String(), b.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Party
	for rows.Next() {
		var i Party
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.IsPublic,
			wkb.Scanner(&i.Location),
			&i.StreetAddress,
			&i.PostalCode,
			&i.State,
			&i.Country,
			&i.StartDate,
			&i.EndDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r PartyRepository) DeleteParty(ctx context.Context, arg DeletePartyParams) error {
	return r.querier.DeleteParty(ctx, arg)
}
