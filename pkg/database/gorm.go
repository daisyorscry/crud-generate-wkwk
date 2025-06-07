package database

import (
	"context"
	"daisy/config"
	"fmt"
	"log"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connection interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
	Get(ctx context.Context) *gorm.DB
}

type SlogWriter struct{}

type connection struct {
	db *gorm.DB
}

func (w *SlogWriter) Printf(format string, v ...interface{}) {
	slog.Info(fmt.Sprintf(format, v...))
}

func setTransactionContext(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, "tx", tx)
}

func (c *connection) Get(ctx context.Context) *gorm.DB {
	tx := ctx.Value("tx")
	if tx != nil {
		return tx.(*gorm.DB)
	}

	return c.db.WithContext(ctx)
}

func (c *connection) Run(ctx context.Context, fn func(context.Context) error) error {
	tx := c.db.Begin()
	defer tx.Rollback()

	newCtx := setTransactionContext(ctx, tx)
	err := fn(newCtx)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func NewConnection(cfg config.Database) (*connection, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConn)
	sqlDB.SetMaxIdleConns(cfg.MinConn)
	sqlDB.SetConnMaxLifetime(cfg.ConnLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.IdleTimeout)

	return &connection{db: db}, nil
}
