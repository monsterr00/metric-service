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
	"github.com/monsterr00/metric-service.gittest_client/internal/config"
	"github.com/monsterr00/metric-service.gittest_client/internal/helpers"
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
	Fetch(ctx context.Context) (map[string]models.Metric, error)
}

const (
	metricsTable   = "metrics"
	migrationsPath = "db/migrations"
)

// New инициализирует соединение с БД и соотвествующие настройки.
func New() *store {
	if config.ServerOptions.Mode == config.DBMode {
		db, err := sql.Open("postgres", config.ServerOptions.DBaddress)
		if err != nil {
			log.Fatal(err)
		}

		filePath := helpers.AbsolutePath("file:///", migrationsPath)

		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Fatal(err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			filePath,
			"postgres", driver)
		if err != nil {
			log.Fatal(err)
		}

		if err = m.Up(); err != nil {
			if err != migrate.ErrNoChange {
				log.Fatal(err)
			}
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

// Ping возвращает статус соединения с БД.
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

// Close закрывает соединение с БД.
func (storl *store) Close() error {
	return storl.conn.Close()
}

// Create создает новую запись в таблице БД metrics.
func (storl *store) Create(ctx context.Context, metric models.Metric) error {
	errPing := storl.conn.Ping()
	if errPing != nil {
		return errPing
	}

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

// Update обновляет запись в таблице БД metrics.
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

// Fetch возвращает полный набор данных из таблицы БД metrics.
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
