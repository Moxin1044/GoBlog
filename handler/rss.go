package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
)

// GetRSS 生成RSS XML
func GetRSS(c *gin.Context) {
	var articles []model.Article
	database.DB.Preload("Category").Preload("Author").
		Where("status = ?", "published").
		Order("published_at DESC").
		Limit(20).
		Find(&articles)

	siteName := getConfigValue("site_name")
	if siteName == "" {
		siteName = "GoBlog"
	}
	siteURL := getConfigValue("site_url")
	if siteURL == "" {
		siteURL = "http://localhost:8080"
	}
	siteDesc := getConfigValue("site_description")
	if siteDesc == "" {
		siteDesc = "A blog powered by GoBlog"
	}

	now := time.Now().Format(time.RFC1123Z)

	xml := `<?xml version="1.0" encoding="UTF-8"?>`
	xml += `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`
	xml += `<channel>`
	xml += fmt.Sprintf(`<title>%s</title>`, siteName)
	xml += fmt.Sprintf(`<link>%s</link>`, siteURL)
	xml += fmt.Sprintf(`<description>%s</description>`, siteDesc)
	xml += fmt.Sprintf(`<lastBuildDate>%s</lastBuildDate>`, now)
	xml += fmt.Sprintf(`<atom:link href="%s/api/rss" rel="self" type="application/rss+xml"/>`, siteURL)

	for _, article := range articles {
		xml += `<item>`
		xml += fmt.Sprintf(`<title><![CDATA[%s]]></title>`, article.Title)
		xml += fmt.Sprintf(`<link>%s/article/%d</link>`, siteURL, article.ID)
		xml += fmt.Sprintf(`<description><![CDATA[%s]]></description>`, article.Summary)
		xml += fmt.Sprintf(`<pubDate>%s</pubDate>`, article.PublishedAt.Format(time.RFC1123Z))
		if article.Category.ID > 0 {
			xml += fmt.Sprintf(`<category>%s</category>`, article.Category.Name)
		}
		xml += fmt.Sprintf(`<guid isPermaLink="true">%s/article/%d</guid>`, siteURL, article.ID)
		xml += `</item>`
	}

	xml += `</channel></rss>`

	c.Header("Content-Type", "application/xml; charset=utf-8")
	c.String(200, xml)
}
