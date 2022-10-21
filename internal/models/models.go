package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no match found")

type SnippetModel struct {
	DB *sql.DB
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	statement := `INSERT INTO snippets (title, content, created, expires) VALUES (
		?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY)
		);`

	result, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets WHERE id = ? AND expires > UTC_TIMESTAMP();`

	row := m.DB.QueryRow(statement, id)

	snippet := &Snippet{}

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	statement := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10;`

	rows, err := m.DB.Query(statement)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		snippet := &Snippet{}
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return snippets, nil
}
