package main

import (
	"context"
	"log"
	"net/http"

	"github.com/andersonribeir0/mouchat/internal/web/database"
	"github.com/andersonribeir0/mouchat/internal/web/server"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
	"go.uber.org/fx"
)

type Config struct {
	// Supabase config
	SbHost   string `env:"SB_HOST"`
	SbSecret string `env:"SB_SECRET"`

	// Web server config
	ListenPort    string `env:"LISTEN_PORT"`
	CallbackURL   string `env:"CALLBACK_URL"`
	SessionSecret string `env:"SESSION_SECRET"`
}

func ProvideConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func ProvideDatabase(cfg *Config) *supabase.Client {
	return database.New(database.Params{
		Host:   cfg.SbHost,
		Secret: cfg.SbSecret,
	})
}

func ProvideServer(cfg *Config, db *supabase.Client) *http.Server {
	return server.NewServer(cfg.ListenPort, db)
}

func StartServer(lc fx.Lifecycle, srv *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Starting server on port %s", srv.Addr)
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					log.Fatalf("Server failed: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server...")
			return srv.Shutdown(ctx)
		},
	})
}

func main() {
	app := fx.New(
		fx.Provide(
			ProvideConfig,
			ProvideDatabase,
			ProvideServer,
		),
		fx.Invoke(StartServer),
	)

	app.Run()
}
