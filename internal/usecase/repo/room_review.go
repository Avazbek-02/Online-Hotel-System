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

// RoomReviewRepo - room_reviews uchun repository
type RoomReviewRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

// NewRoomReviewRepo - room_reviews repoga yangi instance qaytaradi
func NewRoomReviewRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *RoomReviewRepo {
	return &RoomReviewRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

// Create - yangi room review qo'shadi
func (r *RoomReviewRepo) Create(ctx context.Context, req entity.RoomReview) (entity.RoomReview, error) {
	req.ID = uuid.NewString()
	query, args, err := r.pg.Builder.Insert("room_reviews").
		Columns("id", "user_id", "room_id", "rating", "comment").
		Values(req.ID, req.UserID, req.RoomID, req.Rating, req.Comment).
		ToSql()
	if err != nil {
		return entity.RoomReview{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.RoomReview{}, err
	}

	return req, nil
}

// GetSingle - bitta room review ni ID bo'yicha qaytaradi
func (r *RoomReviewRepo) GetSingle(ctx context.Context, req entity.Id) (entity.RoomReview, error) {
	var response entity.RoomReview
	var createdAt, updatedAt time.Time

	queryBuilder := r.pg.Builder.
		Select("id, user_id, room_id, rating, comment, created_at, updated_at").
		From("room_reviews")

	if req.ID == "" {
		return entity.RoomReview{}, fmt.Errorf("GetSingle - invalid request")
	}

	queryBuilder = queryBuilder.Where("id = ?", req.ID)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.RoomReview{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.UserID, &response.RoomID, &response.Rating, &response.Comment, &createdAt, &updatedAt)
	if err != nil {
		return entity.RoomReview{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	return response, nil
}

// GetList - room review lar ro'yhatini qaytaradi
func (r *RoomReviewRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.RoomReviewList, error) {
	var (
		response           = entity.RoomReviewList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select("id, user_id, room_id, rating, comment, created_at, updated_at").
		From("room_reviews")

	// PrepareGetListQuery funksiyasi mavjud bo'lsa, undan foydalaning.
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
		var item entity.RoomReview
		err = rows.Scan(&item.ID, &item.UserID, &item.RoomID, &item.Rating, &item.Comment, &createdAt, &updatedAt)
		if err != nil {
			return response, err
		}

		item.CreatedAt = createdAt.Format(time.RFC3339)
		item.UpdatedAt = updatedAt.Format(time.RFC3339)
		response.Items = append(response.Items, item)
	}

	countQuery, args, err := r.pg.Builder.Select("COUNT(1)").
		From("room_reviews").
		Where(where).ToSql()
	if err != nil {
		return response, err
	}

	err = r.pg.Pool.QueryRow(ctx, countQuery, args...).Scan(&response.Count)
	if err != nil {
		return response, err
	}

	return response, nil
}

// Update - mavjud room review ni yangilaydi
func (r *RoomReviewRepo) Update(ctx context.Context, req entity.RoomReview) (entity.RoomReview, error) {
	updateFields := make(map[string]interface{})

	if req.Rating != 0 {
		updateFields["rating"] = req.Rating
	}
	if req.Comment != "" && req.Comment != "string" {
		updateFields["comment"] = req.Comment
	}

	updateFields["updated_at"] = "now()"

	if len(updateFields) == 0 {
		return entity.RoomReview{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("room_reviews").SetMap(updateFields).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.RoomReview{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.RoomReview{}, err
	}

	return r.GetSingle(ctx, entity.Id{ID: req.ID})
}

// Delete - room review ni ID bo'yicha o'chiradi
func (r *RoomReviewRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("room_reviews").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	return err
}
