package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Role      string         `gorm:"size:20;default:user" json:"role"`     // guest/user/admin/superadmin
	Status    string         `gorm:"size:20;default:active" json:"status"` // active/disabled
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Article 文章模型
type Article struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Title        string         `gorm:"size:200;not null" json:"title"`
	Cover        string         `gorm:"size:255" json:"cover"`
	Content      string         `gorm:"type:longtext" json:"content"`
	Summary      string         `gorm:"size:500" json:"summary"`
	CategoryID   uint           `gorm:"index" json:"category_id"`
	AuthorID     uint           `gorm:"index" json:"author_id"`
	ViewCount    int            `gorm:"default:0" json:"view_count"`
	LikeCount    int            `gorm:"default:0" json:"like_count"`
	CommentCount int            `gorm:"default:0" json:"comment_count"`
	Status       string         `gorm:"size:20;default:draft;index" json:"status"` // draft/published/offline
	PublishedAt  *time.Time     `json:"published_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Author   Admin    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags     []Tag    `gorm:"many2many:article_tags" json:"tags,omitempty"`
}

// Category 分类模型
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	NameEn    string         `gorm:"size:50" json:"name_en"`
	Sort      int            `gorm:"default:0" json:"sort"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Tag 标签模型
type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	NameEn    string         `gorm:"size:50" json:"name_en"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ArticleTag 文章标签关联模型
type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey" json:"article_id"`
	TagID     uint `gorm:"primaryKey" json:"tag_id"`
}

// Comment 评论模型
type Comment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ArticleID  uint           `gorm:"index;not null" json:"article_id"`
	UserID     uint           `gorm:"index" json:"user_id"`
	Nickname   string         `gorm:"size:50" json:"nickname"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	ParentID   uint           `gorm:"index;default:0" json:"parent_id"`
	Status     string         `gorm:"size:20;default:pending;index" json:"status"` // pending/approved/rejected
	IP         string         `gorm:"size:50" json:"ip"`
	ReviewedAt *time.Time     `json:"reviewed_at"`
	ReviewerID uint           `json:"reviewer_id"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Article Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
}

// Like 点赞模型
type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"index;not null" json:"article_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	IP        string    `gorm:"size:50" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

// Subscription 订阅模型
type Subscription struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	Email        string    `gorm:"size:100" json:"email"`
	FeishuToken  string    `gorm:"size:255" json:"feishu_token"`
	NotifyEmail  bool      `gorm:"default:false" json:"notify_email"`
	NotifyFeishu bool      `gorm:"default:false" json:"notify_feishu"`
	CreatedAt    time.Time `json:"created_at"`

	Categories []Category `gorm:"many2many:subscription_categories" json:"categories,omitempty"`
}

// SubscriptionCategory 订阅分类关联模型
type SubscriptionCategory struct {
	SubscriptionID uint `gorm:"primaryKey" json:"subscription_id"`
	CategoryID     uint `gorm:"primaryKey" json:"category_id"`
}

// AIModel AI模型配置
type AIModel struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Provider     string    `gorm:"size:50" json:"provider"`
	APIUrl       string    `gorm:"size:255" json:"api_url"`
	APIType      string    `gorm:"size:10;default:POST" json:"api_type"` // GET/POST
	Headers      string    `gorm:"type:json" json:"headers"`
	BodyTemplate string    `gorm:"type:json" json:"body_template"`
	MaxContext   int       `gorm:"default:10" json:"max_context"`
	Temperature  float64   `gorm:"type:decimal(3,2);default:0.70" json:"temperature"`
	Enabled      bool      `gorm:"default:true" json:"enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserAIConfig 用户AI配置
type UserAIConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	APIToken    string    `gorm:"size:500" json:"api_token"`
	APIUrl      string    `gorm:"size:255" json:"api_url"`
	ModelID     uint      `gorm:"index" json:"model_id"`
	Temperature float64   `gorm:"type:decimal(3,2);default:0.70" json:"temperature"`
	MaxContext  int       `gorm:"default:10" json:"max_context"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Model AIModel `gorm:"foreignKey:ModelID" json:"model,omitempty"`
}

// ChatMessage AI对话消息
type ChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	ArticleID uint      `gorm:"index" json:"article_id"`
	Role      string    `gorm:"size:20;not null" json:"role"` // user/assistant
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// SiteConfig 站点配置KV存储
type SiteConfig struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Admin 管理员模型
type Admin struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email       string         `gorm:"size:100" json:"email"`
	Password    string         `gorm:"size:255;not null" json:"-"`
	Role        string         `gorm:"size:20;default:admin" json:"role"`    // admin/superadmin
	Status      string         `gorm:"size:20;default:active" json:"status"` // active/disabled
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// OperationLog 操作日志
type OperationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	AdminName string    `gorm:"size:50" json:"admin_name"`
	Action    string    `gorm:"size:100;not null" json:"action"`
	Target    string    `gorm:"size:200" json:"target"`
	Result    string    `gorm:"size:500" json:"result"`
	IP        string    `gorm:"size:50" json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

// VerificationCode 验证码
type VerificationCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"index;size:100;not null" json:"email"`
	Code      string    `gorm:"size:10;not null" json:"code"`
	Type      string    `gorm:"size:20;not null" json:"type"` // register/reset
	ExpiredAt time.Time `json:"expired_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// VisitLog 访问日志
type VisitLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	IP        string    `gorm:"index;size:50" json:"ip"`
	Location  string    `gorm:"size:100" json:"location"`
	Path      string    `gorm:"size:255" json:"path"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// Backup 备份
type Backup struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Filename  string    `gorm:"size:255;not null" json:"filename"`
	Size      int64     `json:"size"`
	Type      string    `gorm:"size:20;default:manual" json:"type"` // manual/auto
	CreatedAt time.Time `json:"created_at"`
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{}, &Article{}, &Category{}, &Tag{}, &ArticleTag{},
		&Comment{}, &Like{}, &Subscription{}, &SubscriptionCategory{},
		&AIModel{}, &UserAIConfig{}, &ChatMessage{}, &SiteConfig{},
		&Admin{}, &OperationLog{}, &VerificationCode{}, &VisitLog{},
		&Backup{},
	)
}
