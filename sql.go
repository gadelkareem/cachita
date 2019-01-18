package cachita

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack"
)

var sCache Cache

type sqlCache struct {
	db         *sql.DB
	tableName  string
	ttl        time.Duration
	isPostgres bool
}

type row struct {
	Id        string
	Value     []byte
	ExpiredAt int64
}

func Sql(driverName, dataSourceName string) (Cache, error) {
	if sCache == nil {
		sqlDriver, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		sCache, err = NewSqlCache(24*time.Hour, 5*time.Hour, sqlDriver, "cachita_cache", strings.Contains(dataSourceName, "postgres"))
		if err != nil {
			return nil, err
		}
	}
	return sCache, nil
}

func NewSqlCache(ttl, tickerTtl time.Duration, sql *sql.DB, tableName string, isPostgres ...bool) (Cache, error) {
	c := &sqlCache{
		db:        sql,
		tableName: tableName,
		ttl:       ttl,
	}
	if len(isPostgres) > 0 && isPostgres[0] {
		c.isPostgres = true
	}
	err := c.createTable()
	if err != nil {
		return nil, err
	}

	runEvery(tickerTtl, func() {
		c.deleteExpired()
	})

	return c, nil
}

func (c *sqlCache) Get(key string, i interface{}) error {
	r, err := c.row(Id(key))

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		} else {
			return err
		}
	}
	expiredAt := time.Unix(r.ExpiredAt, 0)
	if expiredAt.Before(time.Now()) {
		return ErrExpired
	}

	return msgpack.Unmarshal(r.Value, i)
}

func (c *sqlCache) row(id string) (*row, error) {
	r := new(row)
	r.Id = id
	query := "SELECT data, expired_at FROM " + c.tableName + " WHERE id = " + c.placeholder(1)
	err := c.db.QueryRow(query, r.Id).Scan(&r.Value, &r.ExpiredAt)
	return r, err
}

func (c *sqlCache) Put(key string, i interface{}, ttl time.Duration) error {
	r, err := c.row(Id(key))
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	data, err := msgpack.Marshal(i)
	if err != nil {
		return err
	}
	var query string
	if r.Value == nil {
		query = "INSERT INTO " + c.tableName + " (data, id, expired_at) VALUES(" + c.placeholder(1) + ", " + c.placeholder(2) + ", " + c.placeholder(3) + ")"
	} else {
		query = "UPDATE " + c.tableName + " SET data = " + c.placeholder(1) + ", expired_at= " + c.placeholder(3) + " WHERE id = " + c.placeholder(2)
	}

	_, err = c.db.Exec(query, data, r.Id, expiredAt(ttl, c.ttl).Unix())
	return err
}

func (c *sqlCache) Incr(key string, ttl time.Duration) (int64, error) {
	var n int64
	r, err := c.row(Id(key))
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if r.Value != nil && time.Unix(r.ExpiredAt, 0).After(time.Now()) {
		err = msgpack.Unmarshal(r.Value, &n)
		if err != nil {
			return 0, err
		}
	}
	n++
	data, err := msgpack.Marshal(n)
	if err != nil {
		return 0, err
	}

	var query string
	if r.Value == nil {
		query = "INSERT INTO " + c.tableName + " (data, id, expired_at) VALUES(" + c.placeholder(1) + ", " + c.placeholder(2) + ", " + c.placeholder(3) + ")"
	} else {
		query = "UPDATE " + c.tableName + " SET data = " + c.placeholder(1) + ", expired_at= " + c.placeholder(3) + " WHERE id = " + c.placeholder(2)
	}

	_, err = c.db.Exec(query, data, r.Id, expiredAt(ttl, c.ttl).Unix())
	return n, err
}

func (c *sqlCache) Invalidate(key string) error {
	_, err := c.db.Exec("DELETE FROM "+c.tableName+" WHERE id = "+c.placeholder(1), Id(key))
	return err
}

func (c *sqlCache) Exists(key string) bool {
	r, _ := c.row(Id(key))
	if r.Value != nil {
		expiredAt := time.Unix(r.ExpiredAt, 0)
		if expiredAt.Before(time.Now()) {
			return false
		}
		return true
	}
	return false
}

func (c *sqlCache) deleteExpired() {
	c.db.Exec("DELETE FROM "+c.tableName+" WHERE expired_at <= "+c.placeholder(1), time.Now().Unix())
}

func (c *sqlCache) createTable() error {
	dataColumnType := "blob"
	if c.isPostgres {
		dataColumnType = "bytea"
	}
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS " + c.tableName + " (id CHAR(32) NOT NULL PRIMARY KEY, data " + dataColumnType + " NOT NULL, expired_at int NOT NULL)")
	if err != nil {
		return err
	}
	return nil
}

func (c *sqlCache) placeholder(index int) string {
	if c.isPostgres {
		return "$" + strconv.Itoa(index)
	}
	return "?"
}
