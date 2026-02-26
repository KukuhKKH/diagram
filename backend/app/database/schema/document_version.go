package schema

import (
	"time"
)

// DocumentVersion represents a version/revision of a document
type DocumentVersion struct {
	ID                uint64    `gorm:"primaryKey" json:"id"`
	DocumentID        uint64    `gorm:"column:document_id;type:bigint;not null;index:idx_version_document_version,unique" json:"document_id"`
	Content           string    `gorm:"column:content;type:text;not null" json:"content"`
	VersionNumber     int       `gorm:"column:version_number;type:integer;not null;index:idx_version_document_version,unique" json:"version_number"`
	AuthorID          *uint64   `gorm:"column:author_id;type:bigint" json:"author_id"`
	ChangeDescription *string   `gorm:"column:change_description;type:varchar(500)" json:"change_description"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime;index:idx_version_document_created" json:"created_at"`

	// Relations
	Document *Document `gorm:"foreignKey:DocumentID;references:ID;OnDelete:CASCADE" json:"-"`
	Author   *User     `gorm:"foreignKey:AuthorID;references:ID;OnDelete:SET NULL" json:"-"`
}

// TableName specifies the table name for DocumentVersion
func (DocumentVersion) TableName() string {
	return "document_versions"
}
