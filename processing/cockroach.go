// Package processing provide all requirements to process change data capture
package processing

import (
	"context"
	"fmt"
	"strings"

	"github.com/Lord-Y/synker/commons"
	"github.com/Lord-Y/synker/models"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

// countChangeFeed permit to check if change feed exist or not
func countChangeFeed(full_table_name, status string) (count int, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.GetPGURI())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT COUNT(job_id) FROM [SHOW CHANGEFEED JOBS] WHERE full_table_names = $1 AND status = $2 and sink_uri = $3 LIMIT 1",
		fmt.Sprintf("{%s}", full_table_name),
		status,
		fmt.Sprintf("kafka://%s", commons.GetKafkaURI()),
	).Scan(&count)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return
}

func createChangeFeed(changefeed models.ChangeFeed) (err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.GetPGURI())
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	//golangci-lint fail on this check while the transaction error is checked
	defer tx.Rollback(ctx) //nolint
	q := fmt.Sprintf(
		"CREATE CHANGEFEED FOR TABLE %s INTO '%s' WITH %s",
		changefeed.FullTableName,
		fmt.Sprintf("kafka://%s", commons.GetKafkaURI()),
		strings.Join(changefeed.Options, ","),
	)

	_, err = tx.Exec(ctx, q)
	if err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		return
	}
	return
}
