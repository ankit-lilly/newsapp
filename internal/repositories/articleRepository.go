package repositories

import (
	"context"
	"database/sql"

	"github.com/ankit-lilly/newsapp/internal/models"
)

type ArticleRepository interface {
	CreateFavoriteArticle(ctx context.Context, article *models.Article) (int64, error)
	GetFavoriteArticleByUser(ctx context.Context, userId int64) ([]models.Article, error)
	GetFavoriteArticle(ctx context.Context, articleId, userId int64) (*models.Article, error)
	DeleteFavoriteArticle(ctx context.Context, articleId, userId int64) error
	IsFavorite(ctx context.Context, link string, userId int64) (int64, error)
}

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *articleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetFavoriteArticleByUser(ctx context.Context, userId int64) ([]models.Article, error) {
	query := "SELECT id, title, content,portal,link, description, user_id, published_at FROM articles WHERE user_id = ?"
	rows, err := r.db.Query(query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Article{}, nil
		}
		return nil, err
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		article := models.Article{}
		if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Portal, &article.Link, &article.Description, &article.UserID, &article.PublishedAt); err != nil {
			return nil, err
		}
		article.IsFavorited = true
		articles = append(articles, article)
	}
	return articles, nil
}

func (r *articleRepository) CreateFavoriteArticle(ctx context.Context, article *models.Article) (int64, error) {
	query := "INSERT INTO articles (title, content, portal, link, description, user_id, published_at) VALUES (?,?,?,?,?,?,?)"
	res, err := r.db.Exec(query, article.Title, article.Content, article.Portal, article.Link, article.Description, article.UserID, article.PublishedAt)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *articleRepository) GetFavoriteArticle(ctx context.Context, articleId, userId int64) (*models.Article, error) {
	query := "SELECT id, title, content, portal, link, description, user_id, published_at FROM articles WHERE id = ? AND user_id = ? LIMIT 1"

	rows := r.db.QueryRow(query, articleId, userId)

	article := models.Article{}

	if err := rows.Scan(&article.ID, &article.Title, &article.Content, &article.Portal, &article.Link, &article.Description, &article.UserID, &article.PublishedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) DeleteFavoriteArticle(ctx context.Context, articleId, userId int64) error {
	query := "DELETE FROM articles WHERE id = ? AND user_id = ?"
	_, err := r.db.Exec(query, articleId, userId)

	return err
}

func (r *articleRepository) IsFavorite(ctx context.Context, link string, userId int64) (int64, error) {
	query := "SELECT id FROM articles WHERE link = ? AND user_id = ?"
	var id int64
	err := r.db.QueryRow(query, link, userId).Scan(&id)
	return id, err
}
