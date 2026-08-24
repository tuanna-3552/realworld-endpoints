package dto

import (
	"time"

	"realworld-endpoints/internal/models"
)

type ArticleDTO struct {
	Slug           string     `json:"slug"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Body           string     `json:"body"`
	TagList        []string   `json:"tagList"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	Favorited      bool       `json:"favorited"`
	FavoritesCount int        `json:"favoritesCount"`
	Author         ProfileDTO `json:"author"`
}

type ArticleResponse struct {
	Article ArticleDTO `json:"article"`
}

type ArticlesResponse struct {
	Articles      []ArticleDTO `json:"articles"`
	ArticlesCount int          `json:"articlesCount"`
}

func ToArticleDTOWithStatus(article *models.Article, favorited bool, following bool) ArticleDTO {
	tagList := make([]string, 0)
	for _, tag := range article.Tags {
		tagList = append(tagList, tag.Name)
	}

	return ArticleDTO{
		Slug:           article.Slug,
		Title:          article.Title,
		Description:    article.Description,
		Body:           article.Body,
		TagList:        tagList,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		Favorited:      favorited,
		FavoritesCount: article.FavoritesCount,
		Author:         ToProfileDTO(&article.Author, following),
	}
}

func ToArticleDTO(article *models.Article) ArticleDTO {
	return ToArticleDTOWithStatus(article, false, false)
}

func ToArticlesDTO(articles []models.Article) []ArticleDTO {
	list := make([]ArticleDTO, 0, len(articles))
	for _, article := range articles {
		list = append(list, ToArticleDTO(&article))
	}
	return list
}
