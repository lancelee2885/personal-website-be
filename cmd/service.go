package cmd

import (
	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lancelee2885/personal-website-be/config"
	internal "github.com/lancelee2885/personal-website-be/internal/service"
	"github.com/lancelee2885/personal-website-be/internal/storage"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-retry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func service() *cobra.Command {
	service := &cobra.Command{
		Use: "service",
		RunE: func(cmd *cobra.Command, args []string) error {
			config.InitializeConfig()

			logger := log.With().Str("service", "service").Logger()

			db, err := sql.Open("postgres", config.FormatDBConfig())
			if err != nil {
				logger.Error().Err(err).Msg("Failed to open database")
				return err
			}

			ctx := context.Background()
			cancel := func() {}
			defer cancel()

			if err := retry.Fibonacci(ctx, 1*time.Second, func(ctx context.Context) error {
				if err := db.PingContext(ctx); err != nil {
					logger.Error().Err(err).Msg("Failed to ping database")
					return retry.RetryableError(err)
				}
				return nil
			}); err != nil {
				logger.Error().Err(err).Msg("Failed to ping database")
			}

			if err := storage.Migrate(db); err != nil {
				logger.Error().Err(err).Msg("Failed to migrate database")
				return err
			}

			store, err := storage.NewPostgresStore(logger, db)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create store")
				return err
			}

			backend := internal.NewService(&internal.ServiceConfig{
				Logger:  logger,
				Storage: store,
			})

			router := gin.New()
			router.Use(gin.Logger())
			router.Use(gin.Recovery())

			// Register routes and handlers
			router.POST("/entities", backend.CreateEntity)
			router.GET("/entities/:id", backend.GetEntity)
			router.PUT("/entities/:id", backend.UpdateEntity)
			router.DELETE("/entities/:id", backend.DeleteEntity)
			router.PATCH("/entities/:id/archive", backend.ArchiveEntity)
			router.GET("/entities", backend.ListEntities)

			// Start the HTTP server
			port := viper.GetString(config.ServerPort)
			logger.Info().Msgf("Starting server on address: %s", port)
			return router.Run(":" + port)

		},
	}
	return service
}
