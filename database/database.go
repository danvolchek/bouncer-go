package database

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"time"
)

// New creates a new database. path is the path to the database.
func New(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: &logger{log: log.Logger},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to db: %s", err)
	}

	err = db.AutoMigrate(allTables...)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate tables: %s", err)
	}

	return db, nil
}

// WithLogger returns a DB session from the given session using the specified logger.
func WithLogger(db *gorm.DB, log zerolog.Logger) *gorm.DB {
	return db.Session(&gorm.Session{Logger: &logger{
		log: log,
	}})
}

// logger is a logger for gorm that uses zerolog.
type logger struct {
	log zerolog.Logger
}

func (l logger) LogMode(_ gormLog.LogLevel) gormLog.Interface {
	// Let zerolog control the log level - do nothing here
	return l
}

func (l logger) Info(_ context.Context, fmt string, args ...interface{}) {
	l.log.Info().Msgf(fmt, args)
}

func (l logger) Warn(_ context.Context, fmt string, args ...interface{}) {
	l.log.Warn().Msgf(fmt, args)

}

func (l logger) Error(_ context.Context, fmt string, args ...interface{}) {
	l.log.Error().Msgf(fmt, args)
}

func (l logger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// This is based off of
	// https://github.com/go-gorm/gorm/blob/master/logger/logger.go#L153
	elapsed := time.Since(begin)

	if err != nil {
		sql, rows := fc()
		l.log.Error().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Err(err).
			Msg("query failed")

		return
	}

	if elapsed >= 200*time.Millisecond {
		sql, rows := fc()
		l.log.Warn().
			Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Msg("slow query")

		return
	}

	if t := l.log.Trace(); t.Enabled() {
		sql, rows := fc()
		t.Dur("elapsed", elapsed).
			Int64("rows", rows).
			Str("sql", sql).
			Msg("sql query")
	}
}
