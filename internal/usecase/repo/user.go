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

type UserRepo struct {
	pg     *postgres.Postgres
	config *config.Config
	logger *logger.Logger
}

func NewUserRepo(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UserRepo {
	return &UserRepo{
		pg:     pg,
		config: config,
		logger: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, req entity.User) (entity.User, error) {
	req.ID = uuid.NewString()
	query, args, err := r.pg.Builder.Insert("users").
		Columns(`id, name, email, password, phone, user_status, gender, role`).
		Values(req.ID, req.Name, req.Email, req.Password_hash, req.Phone, req.UserStatus, req.Gender, req.UserRole).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.User{}, err
	}

	return req, nil
}

func (r *UserRepo) GetSingle(ctx context.Context, req entity.UserSingleRequest) (entity.User, error) {
	var response entity.User
	var createdAt, updatedAt time.Time

	queryBuilder := r.pg.Builder.
		Select(`id, name, email, phone, user_status, gender, role, created_at, updated_at`).
		From("users")

	switch {
	case req.ID != "":
		queryBuilder = queryBuilder.Where("id = ?", req.ID)
	case req.Email != "":
		queryBuilder = queryBuilder.Where("email = ?", req.Email)
	case req.UserRole != "":
		queryBuilder = queryBuilder.Where("role = ?", req.UserRole)
	default:
		return entity.User{}, fmt.Errorf("GetSingle - invalid request")
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	err = r.pg.Pool.QueryRow(ctx, query, args...).
		Scan(&response.ID, &response.Name, &response.Email, &response.Phone, &response.UserStatus, &response.Gender, &response.UserRole, &createdAt, &updatedAt)
	if err != nil {
		return entity.User{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	return response, nil
}

func (r *UserRepo) GetList(ctx context.Context) ([]entity.User, error) {
	var response []entity.User
	query, _, err := r.pg.Builder.
		Select("id, name, email, phone, user_status, gender, role, created_at, updated_at").
		From("users").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.User
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.UserStatus, &user.Gender, &user.UserRole, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		user.CreatedAt = createdAt.Format(time.RFC3339)
		user.UpdatedAt = updatedAt.Format(time.RFC3339)
		response = append(response, user)
	}

	return response, nil
}

func (r *UserRepo) Update(ctx context.Context, req entity.User) (entity.User, error) {
	updateFields := make(map[string]interface{})

	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Email != "" {
		updateFields["email"] = req.Email
	}
	if req.Phone != "" {
		updateFields["phone"] = req.Phone
	}
	if req.UserStatus != "" {
		updateFields["user_status"] = req.UserStatus
	}
	if req.Gender != "" {
		updateFields["gender"] = req.Gender
	}
	if req.UserRole != "" {
		updateFields["role"] = req.UserRole
	}

	updateFields["updated_at"] = "now()"

	if len(updateFields) == 0 {
		return entity.User{}, errors.New("no fields to update")
	}

	query, args, err := r.pg.Builder.Update("users").SetMap(updateFields).Where("id = ?", req.ID).ToSql()
	if err != nil {
		return entity.User{}, err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return entity.User{}, err
	}

	return r.GetSingle(ctx, entity.UserSingleRequest{ID: req.ID})
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	query, args, err := r.pg.Builder.Delete("users").Where("id = ?", id).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	return err
}
