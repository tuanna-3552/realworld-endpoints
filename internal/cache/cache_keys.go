package cache

import (
	"fmt"
	"time"
)

const (
	// TagsCacheKey is the cache key for all tags
	TagsCacheKey = "cache:tags:all"
	// TagsCacheTTL is the TTL for tags cache (10 minutes)
	TagsCacheTTL = 10 * time.Minute

	// ArticleCachePrefix is the prefix for single article cache
	ArticleCachePrefix = "cache:article:"
	// ArticleCacheTTL is the TTL for single article cache (5 minutes)
	ArticleCacheTTL = 5 * time.Minute

	// ArticlesCachePrefix is the prefix for articles list cache
	ArticlesCachePrefix = "cache:articles:"
	// ArticlesCacheTTL is the TTL for articles list cache (2 minutes)
	ArticlesCacheTTL = 2 * time.Minute
)

// ArticleCacheKey returns the cache key for a single article by slug
func ArticleCacheKey(slug string) string {
	return ArticleCachePrefix + slug
}

// ArticlesCacheKey returns the cache key for an articles list query
func ArticlesCacheKey(tag, author, favorited string, limit, offset int) string {
	return fmt.Sprintf("%stag=%s:author=%s:fav=%s:lim=%d:off=%d",
		ArticlesCachePrefix, tag, author, favorited, limit, offset)
}
