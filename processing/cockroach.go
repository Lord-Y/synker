// Package processing provide all requirements to process change data capture
package processing

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Lord-Y/synker/commons"
	"github.com/Lord-Y/synker/models"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// countChangeFeed permit to check if change feed exist or not
func countChangeFeed(full_table_name, status string) (count int, err error) {
	cfg, err := pgxpool.ParseConfig(commons.GetPGURI())
	if err != nil {
		return
	}

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
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

// createChangeFeed with create the change feed in the database
func createChangeFeed(changefeed models.ChangeFeed) (err error) {
	cfg, err := pgxpool.ParseConfig(commons.GetPGURI())
	if err != nil {
		return
	}

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
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

	if !slices.Contains(changefeed.Options, "diff") {
		changefeed.Options = append(changefeed.Options, "diff")
	}
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

// query perform the advanced query on the database in order to retrieve data
func (c *Validate) query(query string, table string, value map[string]interface{}) (z map[string]interface{}, fquery string, err error) {
	var (
		q []string
	)
	for k, v := range value {
		switch c := v.(type) {
		case string:
			q = append(q, fmt.Sprintf("%s.%s = '%s'", table, k, c))
		case int8, int16, int32, int64:
			q = append(q, fmt.Sprintf("%s.%s = %d", table, k, c))
		case float32, float64:
			q = append(q, fmt.Sprintf("%s.%s = %f", table, k, c))
		}
	}

	cfg, err := pgxpool.ParseConfig(commons.GetPGURI())
	if err != nil {
		return
	}

	ctx := context.Background()
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return
	}
	defer db.Close()

	if strings.Contains(strings.ToLower(query), " where ") {
		fquery = query + " AND " + strings.Join(q, " AND ") + " LIMIT 1"
	} else {
		fquery = query + " WHERE " + strings.Join(q, " AND ") + " LIMIT 1"
	}

	rows, err := db.Query(
		ctx,
		fquery,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	var columns []string
	for _, col := range fieldDescriptions {
		columns = append(columns, string(col.Name))
	}
	result := make(map[string]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		values[i] = new(interface{})
	}
	for rows.Next() {
		if err = rows.Scan(values...); err != nil {
			return
		}
	}
	for i, column := range columns {
		val := *(values[i].(*interface{}))
		switch c := val.(type) {
		case string, time.Time, int8, int16, int32, int64, float32, float64:
			result[column] = c
		case pgtype.Numeric:
			var (
				x pgtype.Numeric
				f float64
			)
			x = c
			buf := make([]byte, 0, 32)

			buf = append(buf, x.Int.String()...)
			buf = append(buf, 'e')
			buf = append(buf, strconv.FormatInt(int64(x.Exp), 10)...)

			f, err = strconv.ParseFloat(string(buf), 64)
			if err != nil {
				return
			}
			result[column] = f
		case [16]uint8:
			var (
				x pgtype.UUID
				b []byte
			)
			x.Bytes = c
			x.Status = 2
			b, err = x.MarshalJSON()
			if err != nil {
				return
			}
			result[column] = string(b)
		default:
			result[column] = c
		}
	}
	return result, fquery, nil
}
