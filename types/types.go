package types

import (
	"time"
)

type Users struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Mail      string    `gorm:"size:255;unique;not null" json:"mail"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Username  string    `gorm:"size:20" json:"username,omitempty"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedOn time.Time `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn time.Time `gorm:"autoUpdateTime" json:"updated_on"`
}

type Menu struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Options   string    `gorm:"size:50" json:"option"`
	Is_active bool      `gorm:"default:true" json:"is_active"`
	Is_navbar bool      `gorm:"default:true" json:"is_navbar"`
	CreatedOn time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedOn time.Time `gorm:"autoUpdateTime" json:"-"`
}

type Blogs struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int       `gorm:"foreignKey;not null" json:"user_id"`
	Title      string    `gorm:"size:255;not null" json:"title"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Blog_image []byte    `gorm:"type:bytea" json:"blog_image"`
	IsActive   bool      `gorm:"default:true" json:"-"`
	CreatedOn  time.Time `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn  time.Time `gorm:"autoUpdateTime" json:"updated_on"`
}

type Comments struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int       `gorm:"foreignKey;not null" json:"user_id"`
	BlogID    int       `gorm:"not null" json:"blog_id"`
	Comment   string    `gorm:"size:100;not null" json:"comment"`
	IsActive  bool      `gorm:"default:true" json:"-"`
	CreatedOn time.Time `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn time.Time `gorm:"autoUpdateTime" json:"updated_on"`
}

type Likes struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"-"`
	BlogID    int       `gorm:"not null" json:"blog_id"`
	Likes     int       `gorm:"default:0" json:"likes"`
	CreatedOn time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedOn time.Time `gorm:"autoUpdateTime" json:"-"`
}

type DetailedBlog struct {
	BlogData BlogWithName       `json:"blog_data"`
	Likes    int                `gorm:"default:0" json:"likes"`
	Comments []DetailedComments `gorm:"foreignKey:BlogID" json:"comments"`
}

type DetailedComments struct {
	Comments
	Name string `json:"name"`
	Mail string `json:"mail"`
}

type BlogWithName struct {
	Blogs
	Name string `json:"name"`
}

type Socials struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      int       `gorm:"foreignKey;not null" json:"user_id"`
	SocialMedia string    `gorm:"size:255;not null" json:"social_media"`
	SocialURL   string    `gorm:"type:text;not null" json:"social_url"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedOn   time.Time `gorm:"autoCreateTime" json:"created_on"`
	UpdatedOn   time.Time `gorm:"autoUpdateTime" json:"updated_on"`
}

type LinkedSocials struct {
	UserID   int                      `json:"user_id"`
	Socials  []map[string]interface{} `json:"socials"`
	IsActive bool                     `json:"is_active"`
}

type RecommendationRequest struct {
	Text   string `json:"text"`
	BlogID int    `json:"user_id"`
}

type BlogRecommendation struct {
	BlogData        BlogWithName       `json:"blog_data"`
	Likes           int                `json:"likes"`
	Comments        []DetailedComments `json:"comments,omitempty"`
	SimilarityScore float64            `json:"similarity_score"`
}

type RecommendationResponse struct {
	Recommendations []BlogRecommendation `json:"recommendations"`
}
