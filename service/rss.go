package service

import (
	"fmt"

	"github.com/moxin/GoBlog/database"
	"github.com/moxin/GoBlog/model"
)

func GenerateRSS() (string, error) {
	var articles []model.Article
	database.DB.Where("status = ?", "published").
		Order("published_at DESC").Limit(20).Find(&articles)

	// 获取站点配置
	siteName := getConfigValue("site_name", "GoBlog")
	siteURL := getConfigValue("site_url", "http://localhost:8080")
	siteDesc := getConfigValue("site_description", "A personal blog powered by GoBlog")

	rss := `<?xml version="1.0" encoding="UTF-8"?>`
	rss += `<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">`
	rss += `<channel>`
	rss += fmt.Sprintf(`<title>%s</title>`, siteName)
	rss += fmt.Sprintf(`<link>%s</link>`, siteURL)
	rss += fmt.Sprintf(`<description>%s</description>`, siteDesc)
	rss += fmt.Sprintf(`<atom:link href="%s/api/rss" rel="self" type="application/rss+xml"/>`, siteURL)

	for _, article := range articles {
		rss += `<item>`
		rss += fmt.Sprintf(`<title>%s</title>`, article.Title)
		rss += fmt.Sprintf(`<link>%s/article/%d</link>`, siteURL, article.ID)
		rss += fmt.Sprintf(`<description>%s</description>`, article.Summary)
		rss += fmt.Sprintf(`<pubDate>%s</pubDate>`, article.PublishedAt.Format("Mon, 02 Jan 2006 15:04:05 -0700"))
		rss += fmt.Sprintf(`<guid>%s/article/%d</guid>`, siteURL, article.ID)
		rss += `</item>`
	}

	rss += `</channel></rss>`
	return rss, nil
}

func getConfigValue(key, defaultVal string) string {
	var cfg model.SiteConfig
	if err := database.DB.Where("`key` = ?", key).First(&cfg).Error; err != nil {
		return defaultVal
	}
	if cfg.Value == "" {
		return defaultVal
	}
	return cfg.Value
}
