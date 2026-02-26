package schema

import (
	"time"

	"gorm.io/gorm"
)

// Workspace represents a project/workspace grouping for documents
type Workspace struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(255);not null;index:idx_workspace_owner_name,unique" json:"name"`
	Description *string        `gorm:"column:description;type:text" json:"description"`
	OwnerID     uint64         `gorm:"column:owner_id;type:bigint;not null;index:idx_workspace_owner_name,unique;index:idx_workspace_owner" json:"owner_id"`
	IsPublic    bool           `gorm:"column:is_public;type:boolean;default:false;index" json:"is_public"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Relations
	Owner     *User      `gorm:"foreignKey:OwnerID;references:ID" json:"-"`
	Documents []Document `gorm:"foreignKey:WorkspaceID;references:ID;OnDelete:CASCADE" json:"-"`
}

// TableName specifies the table name for Workspace
func (Workspace) TableName() string {
	return "workspaces"
}
