package database

import (
	"database/sql"

	"github.com/jhphon0730/crawler_auto_dcdc/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// InitDB는 SQLite3 데이터베이스를 초기화하고 연결을 반환
func InitDB(filepath string) error {
	var err error
	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}

	// posts 테이블 생성
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS posts (
			post_number INTEGER PRIMARY KEY,
			title TEXT,
			content TEXT,
			writer TEXT,
			write_date TEXT,
			data_type TEXT
		)`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

// SavePosts는 맵에 있는 게시글을 DB에 저장
func SavePosts(posts map[int]*model.Post) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT OR IGNORE INTO posts (post_number, title, content, writer, write_date, data_type)
		VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, post := range posts {
		_, err := stmt.Exec(post.PostNumber, post.Title, post.Content, post.Writer, post.WriteDate, post.DataType)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// LoadPosts는 DB에서 데이터를 읽어 맵에 채움
func LoadPosts(posts map[int]*model.Post) error {
	rows, err := db.Query("SELECT post_number, title, content, writer, write_date, data_type FROM posts")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		post := &model.Post{}
		err := rows.Scan(&post.PostNumber, &post.Title, &post.Content, &post.Writer, &post.WriteDate, &post.DataType)
		if err != nil {
			return err
		}
		posts[post.PostNumber] = post
	}

	return rows.Err()
}

// CloseDB는 데이터베이스 연결을 닫음 (필요 시 호출)
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
