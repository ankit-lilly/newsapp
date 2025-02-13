package repository

import (
	"database/sql"
	"time"

	"github.com/ankit-lilly/newsapp/pkg/errors"
)

type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	User        int64     `json:"user"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	IsFavorite  bool      `json:"is_favorite"`
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

func (a *ArticleRepository) CreateFavoriteArticle(article_id, user_id int64) error {
	query := `INSERT INTO favorites (article_id, user_id) VALUES($1, $2)`

	_, err := a.DB.Exec(query, article_id, user_id)

	return err
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

func (a *ArticleRepository) GetArticleByID(id int64) (*Article, error) {
	query := `SELECT a.id, a.title, a.description, a.link, a.body, a.created_at, f.user_id FROM articles a LEFT JOIN favorites f ON a.id = f.article_id WHERE a.id = $1`
	var article Article
	var userID sql.NullInt64
	err := a.DB.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Description, &article.Link, &article.Body, &article.CreatedAt, &userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("Article not found")
		}
		return nil, err
	}

	if userID.Valid {
		article.User = userID.Int64
	}

	return &article, nil
}

func (a *ArticleRepository) UpdateArticleByID(article Article) error {
	query := `UPDATE articles SET title = $1, description = $2, link = $3, body = $4 WHERE id = $5`
	_, err := a.DB.Exec(query, article.Title, article.Description, article.Link, article.Body, article.ID)
	return err
}

func (a *ArticleRepository) IsFavoriteArticle(article_id, user_id int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE article_id = $1 AND user_id = $2)`
	var isFavorite bool
	err := a.DB.QueryRow(query, article_id, user_id).Scan(&isFavorite)
	if err != nil {
		return false, err
	}
	return isFavorite, nil
}

func (a *ArticleRepository) DeleteFavoriteArticle(article_id, user_id int64) error {
	query := `DELETE FROM favorites WHERE article_id = $1 AND user_id = $2`
	_, err := a.DB.Exec(query, article_id, user_id)
	return err
}

func (a *ArticleRepository) GetFavoritesByUser(user int64) ([]Article, error) {

	query := `SELECT a.id, a.title, a.link, a.description, a.body, a.created_at FROM favorites as f INNER JOIN articles as a ON f.article_id = a.id  WHERE f.user_id = $1`

	rows, err := a.DB.Query(query, user)

	if err != nil {
		return []Article{}, err
	}

	defer rows.Close()

	articles := []Article{}

	for rows.Next() {
		article := Article{}
		rows.Scan(
			&article.ID,
			&article.Title,
			&article.Link,
			&article.Description,
			&article.Body,
			&article.CreatedAt,
		)
		articles = append(articles, article)
	}

	return articles, nil
}
