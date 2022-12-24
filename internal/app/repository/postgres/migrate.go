package postgres

import (
	"embed"
	"errors"
	"os"
	"time"

	logKey "github.com/apolsh/yapr-url-shortener/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const _defaultAttempts = 5
const _defaultTimeout = 10 * time.Second

//go:embed migrations/*.sql
var fs embed.FS

func RunMigration(databaseURL string) {

	databaseURL += "?sslmode=disable"

	logger := log.With().
		Str(logKey.ComponentKey, "migrations").
		Logger()

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	d, err := iofs.New(fs, "migrations")
	if err != nil {
		logger.Error().Err(err).Msgf("postgres connect error: %s", err)
		os.Exit(1)
	}

	for attempts > 0 {
		//m, err = migrate.New("file://migrations", databaseURL)
		m, err = migrate.NewWithSourceInstance("iofs", d, databaseURL)
		if err == nil {
			break
		}

		logger.Info().Msgf("postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		logger.Error().Err(err).Msgf("postgres connect error: %s", err)
		os.Exit(1)
	}

	if m == nil {
		os.Exit(1)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error().Err(err).Msgf("up error: %s", err)
		os.Exit(1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info().Msg("no change")
		return
	}

	logger.Info().Msg("up success")
}
