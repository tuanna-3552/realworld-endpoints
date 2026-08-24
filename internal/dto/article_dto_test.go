package dto

import (
	"testing"
	"realworld-endpoints/internal/models"
)

func TestToArticleDTO(t *testing.T) {
	article := &models.Article{
		Title:       "Test Title",
		Description: "Test Desc",
		Body:        "Test Body",
		Slug:        "test-title",
		Author: models.User{
			Username: "author",
		},
	}

	result := ToArticleDTO(article)

	if result.Title != article.Title {
		t.Errorf("expected title %s, got %s", article.Title, result.Title)
	}
	if result.Slug != article.Slug {
		t.Errorf("expected slug %s, got %s", article.Slug, result.Slug)
	}
	if result.Description != article.Description {
		t.Errorf("expected description %s, got %s", article.Description, result.Description)
	}
	if result.Body != article.Body {
		t.Errorf("expected body %s, got %s", article.Body, result.Body)
	}
	if result.Favorited != false {
		t.Errorf("expected favorited false")
	}
	if result.Author.Following != false {
		t.Errorf("expected following false")
	}
	if result.Author.Username != "author" {
		t.Errorf("expected author username 'author', got %s", result.Author.Username)
	}
}

func TestToArticleDTOWithStatus(t *testing.T) {
	article := &models.Article{
		Title: "Title",
		Author: models.User{
			Username: "testuser",
		},
	}

	result := ToArticleDTOWithStatus(article, true, true)

	if result.Favorited != true {
		t.Errorf("expected favorited true")
	}
	if result.Author.Following != true {
		t.Errorf("expected following true")
	}
}

func TestToArticleDTOWithStatus_TagList(t *testing.T) {
	article := &models.Article{
		Tags: []models.Tag{
			{Name: "Go"},
			{Name: "Echo"},
		},
		Author: models.User{
			Username: "testuser",
		},
	}

	result := ToArticleDTOWithStatus(article, false, false)

	if len(result.TagList) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(result.TagList))
	}
	if result.TagList[0] != "Go" {
		t.Errorf("expected first tag Go, got %s", result.TagList[0])
	}
	if result.TagList[1] != "Echo" {
		t.Errorf("expected second tag Echo, got %s", result.TagList[1])
	}
}

func TestToArticlesDTO(t *testing.T) {
	articles := []models.Article{
		{Title: "Article 1"},
		{Title: "Article 2"},
	}

	dtos := ToArticlesDTO(articles)

	if len(dtos) != 2 {
		t.Fatalf("expected 2 dtos, got %d", len(dtos))
	}
	if dtos[0].Title != "Article 1" {
		t.Errorf("expected title Article 1, got %s", dtos[0].Title)
	}
	if dtos[1].Title != "Article 2" {
		t.Errorf("expected title Article 2, got %s", dtos[1].Title)
	}
}
