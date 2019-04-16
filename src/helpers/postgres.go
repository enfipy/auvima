package helpers

import (
	"time"

	"github.com/jackc/pgx"
)

type PostgresConnection struct {
	conn *pgx.ConnPool
}

func InitPostgres() *PostgresConnection {
	config, err := pgx.ParseEnvLibpq()
	PanicOnError(err)

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 3,
		AcquireTimeout: 3 * time.Second,
	}

	pc := PostgresConnection{}
	for i := 0; i < 10; i++ {
		pc.conn, err = pgx.NewConnPool(poolConfig)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	PanicOnError(err)

	return &pc
}

func (pc *PostgresConnection) Exec(query string, params ...interface{}) {
	_, err := pc.conn.Exec(query, params...)
	PanicOnError(err)
}

func (pc *PostgresConnection) ExecWithError(query string, params ...interface{}) error {
	_, err := pc.conn.Exec(query, params...)
	if err != nil {
		return err
	}
	return nil
}

func (pc *PostgresConnection) Query(query string, params ...interface{}) func(...interface{}) {
	rows, err := pc.conn.Query(query, params...)
	PanicOnError(err)

	return func(results ...interface{}) {
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(results...)
			PanicOnError(err)
		}
	}
}

func (pc *PostgresConnection) QueryMany(query string, params ...interface{}) *pgx.Rows {
	rows, err := pc.conn.Query(query, params...)
	PanicOnError(err)
	return rows
}
