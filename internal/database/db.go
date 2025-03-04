package database

import (
	"fmt"
	"sync"
	"strconv"
	"database/sql"

	"github.com/jhphon0730/crawler_auto_dcdc/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
	syncDB sync.Once
)


// InitDB는 SQLite3 데이터베이스를 초기화하고 연결을 반환
func InitDB(filepath string) error {
	var err error

	if db != nil {
		return nil
	}

	syncDB.Do(func() {
		db, err = sql.Open("sqlite3", filepath)
		if err != nil {
			return
		}
	})

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

// 데이터를 배열로 반환 ( 필요 시 사용 ) 
// * 페이징 처리를 위해 사용
func LoadPostsByArray(limitStr, pageStr string) ([]model.Post, error) {
	var posts []model.Post

	// 1. 문자열을 정수로 변환
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return posts, fmt.Errorf("invalid limit: %v", limitStr)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return posts, fmt.Errorf("invalid page: %v", pageStr)
	}

	// 2. OFFSET 계산: (page - 1) * limit
	offset := (page - 1) * limit

	// 3. 쿼리 실행
	rows, err := db.Query(`
		SELECT post_number, title, content, writer, write_date, data_type 
		FROM posts 
		ORDER BY post_number DESC 
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return posts, err
	}
	defer rows.Close() // 반드시 rows를 닫아야 함

	// 4. 결과 처리
	for rows.Next() {
		var p model.Post
		err := rows.Scan(&p.PostNumber, &p.Title, &p.Content, &p.Writer, &p.WriteDate, &p.DataType)
		if err != nil {
			return posts, err
		}
		posts = append(posts, p)
	}

	// 5. 에러 체크
	if err = rows.Err(); err != nil {
		return posts, err
	}

	return posts, nil
}

// CloseDB는 데이터베이스 연결을 닫음 (필요 시 호출)
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
