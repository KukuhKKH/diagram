package schema

import (
	"time"

	"gorm.io/gorm"
)

// DocumentType represents the type of document content
type DocumentType string

const (
	DocumentTypeMermaid  DocumentType = "mermaid"
	DocumentTypeMarkdown DocumentType = "markdown"
)

// Document represents a diagram or markdown file
type Document struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	WorkspaceID uint64         `gorm:"column:workspace_id;type:bigint;not null;index:idx_document_workspace;index:idx_document_workspace_created" json:"workspace_id"`
	Title       string         `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Type        DocumentType   `gorm:"column:type;type:varchar(50);not null;default:'mermaid'" json:"type"`
	Slug        string         `gorm:"column:slug;type:varchar(255);not null;index:idx_document_workspace_slug,unique" json:"slug"`
	IsPublic    bool           `gorm:"column:is_public;type:boolean;default:false;index" json:"is_public"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime;index:idx_document_workspace_created" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Relations
	Workspace  *Workspace        `gorm:"foreignKey:WorkspaceID;references:ID;OnDelete:CASCADE" json:"-"`
	Versions   []DocumentVersion `gorm:"foreignKey:DocumentID;references:ID;OnDelete:CASCADE" json:"-"`
	SharedWith []SharedAccess    `gorm:"foreignKey:DocumentID;references:ID;OnDelete:CASCADE" json:"-"`
}

// TableName specifies the table name for Document
func (Document) TableName() string {
	return "documents"
}
