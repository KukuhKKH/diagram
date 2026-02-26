package schema

import (
	"time"

	"gorm.io/gorm"
)

// Permission represents the level of access for a shared document
type Permission string

const (
	PermissionView Permission = "view"
	PermissionEdit Permission = "edit"
)

// SharedAccess represents sharing permissions for documents
type SharedAccess struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	DocumentID  uint64         `gorm:"column:document_id;type:bigint;not null;index:idx_shared_document;index:idx_shared_document_user,unique" json:"document_id"`
	UserID      *uint64        `gorm:"column:user_id;type:bigint;index:idx_shared_document_user,unique" json:"user_id"`
	AccessToken string         `gorm:"column:access_token;type:varchar(255);not null;uniqueIndex:idx_access_token;index" json:"access_token"`
	Permission  Permission     `gorm:"column:permission;type:varchar(50);not null;default:'view'" json:"permission"`
	ExpiresAt   *time.Time     `gorm:"column:expires_at;type:timestamp" json:"expires_at"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Relations
	Document *Document `gorm:"foreignKey:DocumentID;references:ID;OnDelete:CASCADE" json:"-"`
	User     *User     `gorm:"foreignKey:UserID;references:ID;OnDelete:SET NULL" json:"-"`
}

// TableName specifies the table name for SharedAccess
func (SharedAccess) TableName() string {
	return "shared_access"
}
