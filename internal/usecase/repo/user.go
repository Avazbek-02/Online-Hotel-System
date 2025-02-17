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
		Columns(`id, fullname, username, email, password, phone, user_status, gender, role`).
		Values(req.ID, req.FullName, req.UserName, req.Email, req.Password_hash, req.Phone, req.UserStatus, req.Gender, req.UserRole).ToSql()
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
		Select(`id, fullname, username, email, phone, user_status, gender, role, created_at, updated_at`).
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
		Scan(&response.ID, &response.FullName, &response.UserName, &response.Email, &response.Phone, &response.UserStatus, &response.Gender, &response.UserRole, &createdAt, &updatedAt)
	if err != nil {
		return entity.User{}, err
	}

	response.CreatedAt = createdAt.Format(time.RFC3339)
	response.UpdatedAt = updatedAt.Format(time.RFC3339)
	return response, nil
}

func (r *UserRepo) GetList(ctx context.Context, req entity.GetListFilter) (entity.UserList, error) {
	var (
		response             = entity.UserList{}
		createdAt, updatedAt time.Time
	)

	queryBuilder := r.pg.Builder.
		Select("id, fullname, username, email, phone, user_status, gender, role, created_at, updated_at").
		From("users")

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
		var item entity.User
		err = rows.Scan(&item.ID, &item.FullName, &item.UserName, &item.Email, &item.Phone, &item.UserStatus, &item.Gender, &item.UserRole, &createdAt, &updatedAt)
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


func (r *UserRepo) Update(ctx context.Context, req entity.User) (entity.User, error) {
	updateFields := make(map[string]interface{})

	if req.FullName != "" && req.FullName != "string"{
		updateFields["fullname"] = req.FullName
	}
	if req.UserName != "" && req.UserName != "string"{
		updateFields["username"] = req.UserName
	}
	if req.Email != "" && req.Email != "string"{
		updateFields["email"] = req.Email
	}
	if req.Phone != "" && req.Phone != "string"{
		updateFields["phone"] = req.Phone
	}
	if req.UserStatus != "" && req.UserStatus != "string"{
		updateFields["user_status"] = req.UserStatus
	}
	if req.Gender != "" && req.Gender != "string"{
		updateFields["gender"] = req.Gender
	}
	if req.UserRole != "" && req.UserRole != "string"{
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

func (r *UserRepo) Delete(ctx context.Context, req entity.Id) error {
	query, args, err := r.pg.Builder.Delete("users").Where("id = ?", req.ID).ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, query, args...)
	return err
}
