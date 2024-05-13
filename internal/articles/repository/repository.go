package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type ArticleRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{DB: db}
}

func (a *ArticleRepository) GetAllArticles() ([]Article, error) {

	query := `SELECT id, title, link, body, created_at FROM articles ORDER BY created_at DESC`

	rows, err := a.DB.Query(query)
	if err != nil {
		return []Article{}, err
	}
	// We close the resource
	defer rows.Close()

	articles := []Article{}
	for rows.Next() {
		article := Article{}
		rows.Scan(
			&article.ID,
			&article.Title,
			&article.Link,
			&article.Body,
			&article.CreatedAt,
		)

		articles = append(articles, article)
	}

	return articles, nil
}

func (a *ArticleRepository) InsertArticle(article Article) (int, error) {
	query := `INSERT INTO articles (title, link, body, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := a.DB.QueryRow(query, article.Title, article.Link, article.Body, article.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *ArticleRepository) BulkInsertArticles(articles []Article) error {
	tx, err := a.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO articles (id, title, link, description, body, created_at) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, article := range articles {
		_, err := stmt.Exec(article.ID, article.Title, article.Link, article.Description, article.Body, article.CreatedAt)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (a *ArticleRepository) DeleteArticle(id int) error {
	query := `DELETE FROM articles WHERE id = $1`

	result, err := a.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (a *ArticleRepository) GetArticleByID(id int) (*Article, error) {
	query := `SELECT id, title, description, link, body, created_at FROM articles WHERE id = $1`
	var article Article
	err := a.DB.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Description, &article.Link, &article.Body, &article.CreatedAt)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Article not found")
		}
		return nil, err
	}
	return &article, nil
}

func (a *ArticleRepository) UpdateArticleByID(article Article) error {
	query := `UPDATE articles SET title = $1, description = $2, link = $3, body = $4 WHERE id = $5`
	_, err := a.DB.Exec(query, article.Title, article.Description, article.Link, article.Body, article.ID)
	return err
}
