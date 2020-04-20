package cachita

import (
	"database/sql"
	"fmt"
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

type tagRow struct {
	Id   string
	Keys string
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
	_, err = c.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s_tags (id CHAR(32) NOT NULL PRIMARY KEY, keys TEXT NOT NULL)", c.tableName))
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

func (c *sqlCache) InvalidateMulti(keys ...string) error {
	var ids []string
	for _, key := range keys {
		ids = append(ids, Id(key))
	}
	_, err := c.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", c.tableName, c.placeholder(1)), ids)
	return err
}

func (c *sqlCache) Tag(key string, tags ...string) (err error) {
	id := Id(key)
	var r *tagRow
	for _, t := range tags {
		r, err = c.tagRow(Id(t))
		if err != nil && err != sql.ErrNoRows {
			return
		}
		r.Keys += fmt.Sprintf(",%s", id)
		var query string
		if err == sql.ErrNoRows {
			query = fmt.Sprintf("INSERT INTO %s_tags (keys, id) VALUES(%s, %s)", c.tableName, c.placeholder(1), c.placeholder(2))
		} else {
			query = fmt.Sprintf("UPDATE %s_tags SET keys = %s WHERE id = %s ", c.tableName, c.placeholder(1), c.placeholder(2))
		}
		_, err = c.db.Exec(query, r.Keys, r.Id)
		if err != nil {
			return
		}
	}
	return nil
}

func (c *sqlCache) tagRow(id string) (r *tagRow, err error) {
	r = new(tagRow)
	r.Id = id
	query := fmt.Sprintf("SELECT keys FROM %s_tags WHERE id = %s", c.tableName, c.placeholder(1))
	err = c.db.QueryRow(query, r.Id).Scan(&r.Keys)
	return
}

func (c *sqlCache) InvalidateTags(tags ...string) (err error) {
	var keys string
	var r *tagRow
	for _, t := range tags {
		r, err = c.tagRow(Id(t))
		if err != nil && err != sql.ErrNoRows {
			return
		}
		if err == sql.ErrNoRows || r == nil {
			continue
		}
		keys += fmt.Sprintf(",%s", r.Keys)
	}
	s := sqlKeys(keys)
	_, err = c.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id IN (%s)", c.tableName, s))

	return
}
func sqlKeys(s string) string {
	ids := strings.Split(s, ",")
	var (
		l []string
		r string
	)
	for _, v := range ids {
		if !inArr(l, v) && v != "," && v != "" {
			l = append(l, v)
			r += fmt.Sprintf("'%s',", v)
		}
	}
	return strings.TrimRight(r, ",")
}
