package command

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/handler/routes"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/model"
	"github.com/kickplate/api/repository"
	"github.com/spf13/cobra"
)

type ServeCommand struct{}

func (c ServeCommand) Short() string {
	return "Start the server"
}

func (c ServeCommand) Setup(cmd *cobra.Command) {}

func (s *ServeCommand) Run() lib.CommandRunner {
	return func(
		env lib.Env,
		logger lib.Logger,
		handler lib.RequestHandler,
		accountRepo repository.AccountRepository,
		badgeRepo repository.BadgeRepository,
		r routes.Routes,
	) {
		seedBadges(env, logger, badgeRepo)

		// CORS middleware first
		handler.Mux.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
				next.ServeHTTP(w, r)
			})
		})

		// Optional authentication (non-blocking, just extracts token if present)
		handler.Mux.Use(middleware.Authenticate(env, logger))
		handler.Mux.Use(middleware.HeaderAuth(env, accountRepo, logger))

		r.Setup()

		addr := fmt.Sprintf(":%s", env.ServerPort)
		logger.Info("Running server on port ", env.ServerPort)

		if err := http.ListenAndServe(addr, handler.Mux); err != nil {
			logger.Fatal("Server failed: ", err)
		}
	}
}

func seedBadges(env lib.Env, logger lib.Logger, badgeRepo repository.BadgeRepository) {
	if len(env.Badges) == 0 {
		return
	}
	ctx := context.Background()
	logger.Info("Seeding badges from config...")
	for _, bc := range env.Badges {
		existing, err := badgeRepo.GetBySlug(ctx, bc.Slug)
		if err != nil {
			logger.Errorf("  ✗ %s: %v", bc.Slug, err)
			continue
		}
		if existing != nil {
			continue
		}
		badge := model.Badge{
			ID:          uuid.New(),
			Slug:        bc.Slug,
			Name:        bc.Name,
			Description: bc.Description,
			Icon:        bc.Icon,
			Tier:        model.BadgeTier(bc.Tier),
		}
		if err := badgeRepo.Create(ctx, &badge); err != nil {
			logger.Errorf("  ✗ %s: %v", bc.Slug, err)
			continue
		}
		logger.Infof("  ✓ seeded badge: %s", bc.Slug)
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
