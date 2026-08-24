package cache

import (
	"testing"
	"time"
)

func TestArticleCacheKey(t *testing.T) {
	key := ArticleCacheKey("how-to-train-your-dragon")
	expected := "cache:article:how-to-train-your-dragon"
	if key != expected {
		t.Errorf("ArticleCacheKey() = %q, want %q", key, expected)
	}
}

func TestArticlesCacheKey(t *testing.T) {
	key := ArticlesCacheKey("golang", "johndoe", "", 10, 0)
	expected := "cache:articles:tag=golang:author=johndoe:fav=:lim=10:off=0"
	if key != expected {
		t.Errorf("ArticlesCacheKey() = %q, want %q", key, expected)
	}
}

func TestArticlesCacheKey_Empty(t *testing.T) {
	key := ArticlesCacheKey("", "", "", 20, 0)
	expected := "cache:articles:tag=:author=:fav=:lim=20:off=0"
	if key != expected {
		t.Errorf("ArticlesCacheKey() = %q, want %q", key, expected)
	}
}

func TestCacheConstants(t *testing.T) {
	if TagsCacheTTL != 10*time.Minute {
		t.Errorf("TagsCacheTTL = %v, want %v", TagsCacheTTL, 10*time.Minute)
	}
	if ArticleCacheTTL != 5*time.Minute {
		t.Errorf("ArticleCacheTTL = %v, want %v", ArticleCacheTTL, 5*time.Minute)
	}
	if ArticlesCacheTTL != 2*time.Minute {
		t.Errorf("ArticlesCacheTTL = %v, want %v", ArticlesCacheTTL, 2*time.Minute)
	}
	if TagsCacheKey != "cache:tags:all" {
		t.Errorf("TagsCacheKey = %q, want %q", TagsCacheKey, "cache:tags:all")
	}
}
