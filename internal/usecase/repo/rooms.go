package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Avazbek-02/Online-Hotel-System/config"
	"github.com/Avazbek-02/Online-Hotel-System/internal/entity"
	"github.com/Avazbek-02/Online-Hotel-System/pkg/logger"
	"github.com/Avazbek-02/Online-Hotel-System/pkg/postgres"
	"github.com/google/uuid"
)

type RoomsRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewRoomsRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *RoomsRepo {
	return &RoomsRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *RoomsRepo) Create(ctx context.Context, req entity.Room) (entity.Room, error) {
	req.ID = uuid.NewString()
	query, args, err := r.pg.Builder.Insert("rooms").
		Columns(`id, type, category, status, price, availability, rating`).
		Values(req.ID, req.Type, req.Category, req.Status, req.Price, req.Availability, req.Rating).ToSql()
	if err != nil {
		return entity.Room{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Room{}, err
	}

	return req, nil
}

func (r *RoomsRepo) GetSingle(ctx context.Context, req entity.Id) (entity.Room, error) {
	var response entity.Room
	var createdAt, updatedAt time.Time

	queryBuilder := r.pg.Builder.
		Select(`id, type, category, status, price, availability, rating, created_at, updated_at`).
		From("rooms")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	default:
		return entity.Room{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.Room{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.Type, &response.Category, &response.Status, &response.Price, &response.Availability, &response.Rating, &createdAt, &updatedAt)
	if err != nil {
		return entity.Room{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	return response, nil
}

func (r *RoomsRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.RoomList, error) {
	var (
		response             = entity.RoomList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select("id, type, category, status, price, availability, rating, created_at, updated_at").
		From("rooms")

	queryBuilder, where := PrepareGetListQuery(queryBuilder, req)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := r.pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return response, err
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.Room
		err = rows.Scan(&item.ID, &item.Type, &item.Category, &item.Status, &item.Price, &item.Availability, &item.Rating, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)

		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").From("users").Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (r *RoomsRepo) Update(ctx context.Context, req entity.Room) (entity.Room, error) {
	updateFields := make(map[string]interface{})

	if req.Type != "" && req.Type != "string" {
		updateFields["type"] = req.Type
	}
	if req.Category != "" && req.Category != "string" {
		updateFields["category"] = req.Category
	}
	if req.Status != "" && req.Status != "string" {
		updateFields["status"] = req.Status
	}
	if req.Price != 0{
		updateFields["price"] = req.Price
	}
	if req.Availability{
		updateFields["availability"] = req.Availability
	}

	updateFields["updated_at"] = "now()"

	if len(updateFields) == 0 {
		return entity.Room{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("rooms").SetMap(updateFields).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.Room{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.Room{}, err
	}

	return r.GetSingle(ctx, entity.Id{ID: req.ID})
}

func (r *RoomsRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("rooms").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	return err
}
