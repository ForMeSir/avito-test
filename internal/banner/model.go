package banner

import (
	"time"

	"github.com/LukaGiorgadze/gonull"
	"github.com/lib/pq"
)

type Content struct {
	Title    string `json:"title" db:"content_title"`
	Text     string `json:"text" db:"content_text"`
	Url      string `json:"url" db:"content_url"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type ContentWithOut struct{
	Title    string `json:"title" db:"content_title"`
	Text     string `json:"text" db:"content_text"`
	Url      string `json:"url" db:"content_url"`
}

type Banner struct {
	Token     string  `json:"token"`
	TagIds    []int   `json:"tag_ids"    db:"tag_ids"`
	FeatureId int     `json:"feature_id" db:"feature_id"`
	Content   Content `json:"content"`
}

type BannerCache struct {
	TagId     int
	FeatureId int
}

type ContentCache struct {
	Time  int64  `json:"time,omitempty"`
	Title string `json:"title" db:"content_title"`
	Text  string `json:"text" db:"content_text"`
	Url   string `json:"url" db:"content_url"`
	IsActive bool `json:"is_active"  db:"is_active"`
}

type FullBanner struct {
	Id        int               `json:"id"         db:"id"`
	TagIds    pq.Int64Array     `json:"tag_ids"    db:"tag_ids"`
	FeatureId int               `json:"feature_id" db:"feature_id"`
	Title     string            `json:"title"      db:"content_title"`
	Text      string            `json:"text"       db:"content_text"`
	Url       string            `json:"url"        db:"content_url"`
	IsActive  bool              `json:"is_active"  db:"is_active"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}

type UpdateBody struct {
	TagIds    gonull.Nullable[[]int]        `json:"tag_ids"    binding:"required"`
	FeatureId gonull.Nullable[int]          `json:"feature_id" binding:"required"`
	Content   gonull.Nullable[ContentWithOut] `json:"content"    binding:"required"`
	IsActive  gonull.Nullable[bool]         `json:"is_active"  binding:"required"`
}
