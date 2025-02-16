package usecase

import (
	"github.com/Avazbek-02/Online-Hotel-System/config"
	"github.com/Avazbek-02/Online-Hotel-System/internal/usecase/repo"
	"github.com/Avazbek-02/Online-Hotel-System/pkg/logger"
	"github.com/Avazbek-02/Online-Hotel-System/pkg/postgres"
)

// UseCase -.
type UseCase struct {
	UserRepo         UserRepoI
	SessionRepo      SessionRepoI
}

// New -.
func New(pg *postgres.Postgres, config *config.Config, logger *logger.Logger) *UseCase {
	return &UseCase{
		UserRepo:         repo.NewUserRepo(pg, config, logger),
		SessionRepo:      repo.NewSessionRepo(pg, config, logger),
	}
}
