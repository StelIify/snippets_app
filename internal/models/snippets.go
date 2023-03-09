package models

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Snippet struct {
	Id int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	statement := `INSERT INTO snippets (title, content, created, expires) VALUES (
						$1,
						$2,
						CURRENT_TIMESTAMP,
						CURRENT_TIMESTAMP + INTERVAL '` + strconv.Itoa(expires) + ` days'
					) RETURNING id`
	var id int64
	err := m.DB.QueryRow(context.Background(), statement, title, content).Scan(&id)
	if err != nil{
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error){
	statement := `SELECT id, title, content, created, expires FROM snippets
				  WHERE expires > CURRENT_TIMESTAMP and id=$1`
	row := m.DB.QueryRow(context.Background(), statement, id)

	s := &Snippet{}

	err := row.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrNoRecord
		}else{
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error){
	statement := `SELECT id, title, content, created, expires FROM snippets WHERE expires > CURRENT_TIMESTAMP ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(context.Background(), statement)

	if err != nil{
		return nil, err
	}

	defer rows.Close()
	snippets := []*Snippet{}

	for rows.Next(){
		s := &Snippet{}
		err := rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return snippets, nil
}

