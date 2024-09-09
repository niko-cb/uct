package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sqlcommentercore "github.com/google/sqlcommenter/go/core"
	gosql "github.com/google/sqlcommenter/go/database/sql"
	"github.com/niko-cb/uct/internal/infrastructure/monitor/log"
	"github.com/niko-cb/uct/internal/infrastructure/web/config"
)

type MySQLClient struct {
	*sql.DB
	dbHost string
	dbPort string
	dbUser string
	dbPass string
	dbName string
	once   sync.Once
}

func NewMySQLClient(cfg *config.Config) *MySQLClient {
	return &MySQLClient{
		DB:     &sql.DB{},
		dbUser: cfg.DbUser,
		dbPass: cfg.DbPass,
		dbHost: cfg.DbHost,
		dbPort: cfg.DbPort,
		dbName: cfg.DbName,
	}
}

func (c *MySQLClient) Connect() {
	// We only want to establish the connection once
	// Other calls to Connect() should not re-establish the connection, but just reuse the existing one
	c.once.Do(func() {
		dbPool, err := gosql.Open("mysql", c.uri(),
			sqlcommentercore.CommenterOptions{},
		)
		if err != nil {
			log.Error(context.Background(), fmt.Errorf("failed to open database connection: %+v", err))
			return
		}

		dbPool.SetMaxIdleConns(10)
		dbPool.SetMaxOpenConns(90)
		dbPool.SetConnMaxLifetime(100 * time.Second)

		c.DB = dbPool
	})
}

// uri returns the connection string for the MySQL database
func (c *MySQLClient) uri() string {
	connectionType := "tcp"
	connectionDetails := fmt.Sprintf("(%s:%s)", c.dbHost, c.dbPort)

	dbURI := fmt.Sprintf("%s:%s@%s%s/%s?parseTime=true",
		c.dbUser, c.dbPass, connectionType, connectionDetails, c.dbName)

	return dbURI
}

func (c *MySQLClient) CtxTxKey() interface{} {
	return "tx"
}

// DoInTx executes the given function in a transaction
func (c *MySQLClient) DoInTx(ctx context.Context, f func(context.Context) error) error {
	c.Connect() // Ensure connection is established

	tx, err := c.Begin()
	if err != nil {
		log.Error(ctx, fmt.Errorf("%+v\n", err))
		return err
	}
	ctx = context.WithValue(ctx, c.CtxTxKey(), tx)
	log.Debug(ctx, "transaction started")

	var done bool
	defer func() {
		if !done {
			log.Debug(ctx, "transaction rollback")
			err = tx.Rollback()
			if err != nil {
				log.Error(ctx, fmt.Errorf("%+v\n", err))
			}
		}
	}()

	if err = f(ctx); err != nil {
		log.Error(ctx, fmt.Errorf("%+v\n", err))
		return err
	}

	done = true
	if err = tx.Commit(); err != nil {
		log.Error(ctx, fmt.Errorf("%+v\n", err))
		return err
	}

	log.Debug(ctx, "transaction done and committed")
	return nil
}
