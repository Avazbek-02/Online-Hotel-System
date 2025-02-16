package handler

import (
	"github.com/Avazbek-02/Online-Hotel-System/config"
	"github.com/Avazbek-02/Online-Hotel-System/internal/usecase"
	minio "github.com/Avazbek-02/Online-Hotel-System/pkg/MinIO"
	"github.com/Avazbek-02/Online-Hotel-System/pkg/logger"
	rediscache "github.com/golanguzb70/redis-cache"
)

type Handler struct {
	Logger  *logger.Logger
	Config  *config.Config
	UseCase *usecase.UseCase
	Redis   rediscache.RedisCache
	MinIO   *minio.MinIO
}

func NewHandler(l *logger.Logger, c *config.Config, useCase *usecase.UseCase, redis rediscache.RedisCache, minio *minio.MinIO) *Handler {
	return &Handler{
		Logger:  l,
		Config:  c,
		UseCase: useCase,
		Redis:   redis,
		MinIO:   minio,
	}
}
