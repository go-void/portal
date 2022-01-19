package collector

import (
	"database/sql"
	"fmt"
	"strings"
)

// MySQLStore stores flushed entries to a MySLQ / MariaDB database
type MySQLStore struct {
	conn *sql.DB
}

func (s *MySQLStore) StoreEntries(entries []Entry) error {
	placeholders := make([]string, len(entries))
	values := make([]interface{}, len(entries)*8)

	for _, entry := range entries {
		placeholders = append(placeholders, "(?,?,?,?,?,?,?,?)")
		values = append(values,
			entry.ID,
			entry.Question,
			entry.Answer,
			entry.QueryTime,
			entry.ClientIP,
			entry.Filtered,
			entry.Cached,
		)
	}

	st := fmt.Sprintf(
		"INSERT INTO queries (id, question, answer, query_time, client_ip, filtered, cached) VALUES %s",
		strings.Join(placeholders, ","),
	)

	stmt, err := s.conn.Prepare(st)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	return err
}
