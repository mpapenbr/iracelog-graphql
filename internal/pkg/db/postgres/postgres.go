package database

import (
	"context"
	"log"
	"os"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var DbPool *pgxpool.Pool

func InitDB() *pgxpool.Pool {
	return InitWithUrl(os.Getenv("DATABASE_URL"))
}

func InitWithUrl(url string) *pgxpool.Pool {
	dbConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Unable to parse database config %v\n", err)
	}

	logger := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.TraceLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	dbConfig.ConnConfig.Tracer = &myQueryTracer{logger, logrus.TraceLevel}

	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	DbPool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database %v\n", err)
	}
	return DbPool
}

func CloseDb() {
	DbPool.Close()
}

type myQueryTracer struct {
	log   logrus.FieldLogger
	level logrus.Level
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	// do the logging
	tracer.log.WithField("sql", data.SQL).WithField("args", data.Args).Log(tracer.level, "Query started")
	return ctx
}

//nolint:whitespace // can't make the linters happy
func (tracer *myQueryTracer) TraceQueryEnd(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
}
