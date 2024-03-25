package postgres

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	// The blank import is used to ensure that the "github.com/jackc/pgx/v5/stdlib" package's init functions are executed, which registers the PostgreSQL driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	timezone                 = "utc"
	pingTimeout              = 1000 * time.Millisecond
	maxPingAttempts          = 10
	sleepTimeBetweenAttempts = 100
)

type (
	Config struct {
		User       string
		Password   string
		Host       string
		DBName     string
		DisableTLS bool
	}

	Client struct {
		*sqlx.DB
	}
)

func NewClient(config Config) (*Client, error) {
	sslMode := "require"
	if config.DisableTLS {
		sslMode = "disable"
	}

	query := make(url.Values)

	query.Set("sslmode", sslMode)
	query.Set("timezone", timezone)

	URL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.User, config.Password),
		Host:     config.Host,
		Path:     config.DBName,
		RawQuery: query.Encode(),
	}

	db, err := sqlx.Open("pgx", URL.String())
	if err != nil {
		return &Client{}, err
	}

	client := &Client{db}

	if err := client.Ping(context.Background()); err != nil {
		return &Client{}, errors.New("failed to connect to the database")
	}

	return client, nil
}

func (client *Client) Ping(ctx context.Context) error {
	if _, ok := ctx.Deadline(); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, pingTimeout)
		defer cancel()
	}

	for attempts := 1; ; attempts++ {
		err := client.DB.Ping()
		if err == nil {
			break
		}

		if attempts == maxPingAttempts {
			return err
		}

		time.Sleep(time.Duration(attempts) * sleepTimeBetweenAttempts * time.Millisecond)
	}

	var result bool
	return client.DB.QueryRowContext(ctx, "SELECT true").Scan(&result)
}

func (client *Client) Close() error {
	return client.DB.Close()
}
