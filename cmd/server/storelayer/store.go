package storelayer

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/lib/pq"

	/*"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	*/
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/models"
)

type store struct {
	conn *sql.DB
}

type Store interface {
	Ping() error
	Close() error
	Create(ctx context.Context, metric models.Metric) error
	Update(ctx context.Context, metric models.Metric) error
	GetByID(ctx context.Context, id string, mtype string) (models.Metric, error)
	Fetch(ctx context.Context) (map[string]models.Metric, error)
}

const (
	metricsTable   = "metrics"
	migrationsPath = "file:///Users/denis/metric-service/db/migrations"
)

func New() *store {
	if config.ServerOptions.Mode == config.DBMode {
		db, err := sql.Open("postgres", config.ServerOptions.DBaddress)
		if err != nil {
			log.Fatal(err)
		}

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		m, err := migrate.NewWithDatabaseInstance(
			migrationsPath,
			"postgres", driver)
		m.Up()
		if err != nil {
			log.Fatal(err)
		}

		storl := &store{
			conn: db,
		}

		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}
		return storl
	}

	storl := &store{
		conn: nil,
	}

	return storl
}

func (storl *store) Ping() error {
	if storl.conn == nil {
		return errors.New("db: not started")
	}

	err := storl.conn.Ping()

	if err != nil {
		pgerr, _ := errors.Unwrap(err).(*pgconn.PgError)
		isCl08 := pgerrcode.IsConnectionException(string(pgerr.Code))
		if isCl08 {
			var timeout = 1
			for i := 0; i < config.ServerOptions.ReconnectCount; i++ {
				time.Sleep(time.Duration(timeout) * time.Second)

				err = storl.conn.Ping()
				if err != nil {
					timeout += config.ServerOptions.ReconnectDelta
				} else {
					break
				}

			}
			return err
		}
	}
	return nil
}

func (storl *store) Close() error {
	return storl.conn.Close()
}

func (storl *store) Create(ctx context.Context, metric models.Metric) error {
	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	INSERT INTO metrics 
	(ID, MType, Delta, Value)
	VALUES
	($1, $2, $3, $4);
    `, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		// если ошибка, то откатываем изменения
		return errors.Join(err, tx.Rollback())
	}
	return tx.Commit()
}

func (storl *store) Update(ctx context.Context, metric models.Metric) error {
	tx, err := storl.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE metrics
	SET Delta = $3,
		Value = $4
	WHERE ID = $1 and MType = $2;		
    `, metric.ID, metric.MType, metric.Delta, metric.Value)
	if err != nil {
		// если ошибка, то откатываем изменения
		return errors.Join(err, tx.Rollback())
	}
	// коммитим транзакцию
	return tx.Commit()
}

func (storl *store) GetByID(ctx context.Context, id string, mtype string) (models.Metric, error) {
	row := storl.conn.QueryRowContext(ctx, `	
	SELECT 
		ID,
		MType,
		Delta,
		Value
	FROM metrics
	WHERE id = $1
	AND   mtype  = $2
    `,
		id, mtype,
	)

	var metric models.Metric
	err := row.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
	if err != nil {
		return metric, err
	}

	return metric, nil
}

func (storl *store) Fetch(ctx context.Context) (map[string]models.Metric, error) {
	rows, err := storl.conn.QueryContext(ctx, `
	SELECT
		ID,
		MType,
		Delta,
		Value
	FROM metrics	
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	metrics := make(map[string]models.Metric)

	for rows.Next() {
		var m models.Metric
		if err := rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err != nil {
			return nil, err
		}
		metrics[m.ID] = m
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}
