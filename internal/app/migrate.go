//go:build migrate

package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bhcoder23/gin-clean-template/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	l := logger.New("info")

	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		l.Fatal("migrate: environment variable not declared: PG_URL")
	}

	databaseURL += "?sslmode=disable"

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		l.Info("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		l.Fatal(fmt.Errorf("Migrate: postgres connect error: %w", err))
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Fatal(fmt.Errorf("Migrate: up error: %w", err))
	}

	if errors.Is(err, migrate.ErrNoChange) {
		l.Info("Migrate: no change")
		return
	}

	l.Info("Migrate: up success")
}
