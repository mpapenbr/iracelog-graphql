package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

var DbPool *pgxpool.Pool

func InitDB() {
	dbConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to parse database config %v\n", err)
	}

	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}

	DbPool, err = pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database %v\n", err)
	}
}

func CloseDb() {
	DbPool.Close()
}
