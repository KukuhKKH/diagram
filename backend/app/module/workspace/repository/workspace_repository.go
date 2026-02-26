package repository

import (
	"git.dev.siap.id/kukuhkkh/app-diagram/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-diagram/internal/bootstrap/database"
)

// WorkspaceRepository
type WorkspaceRepository interface {
	Create(workspace *schema.Workspace) (*schema.Workspace, error)
	FindByID(id uint64) (*schema.Workspace, error)
	FindByOwnerID(ownerID uint64, limit, offset int) ([]schema.Workspace, error)
	CountByOwnerID(ownerID uint64) (int64, error)
	Update(workspace *schema.Workspace) error
	Delete(id uint64) error
	CheckNameExists(name string, ownerID uint64, excludeID uint64) bool
}

type workspaceRepository struct {
	db *database.Database
}

func NewWorkspaceRepository(db *database.Database) WorkspaceRepository {
	return &workspaceRepository{
		db: db,
	}
}

func (_i *workspaceRepository) Create(workspace *schema.Workspace) (*schema.Workspace, error) {
	if err := _i.db.DB.Create(workspace).Error; err != nil {
		return nil, err
	}
	return workspace, nil
}

func (_i *workspaceRepository) FindByID(id uint64) (*schema.Workspace, error) {
	var workspace schema.Workspace
	if err := _i.db.DB.Where("id = ?", id).First(&workspace).Error; err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (_i *workspaceRepository) FindByOwnerID(ownerID uint64, limit, offset int) ([]schema.Workspace, error) {
	var workspaces []schema.Workspace
	if err := _i.db.DB.Where("owner_id = ?", ownerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&workspaces).Error; err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (_i *workspaceRepository) CountByOwnerID(ownerID uint64) (int64, error) {
	var count int64
	if err := _i.db.DB.Model(&schema.Workspace{}).
		Where("owner_id = ?", ownerID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (_i *workspaceRepository) Update(workspace *schema.Workspace) error {
	return _i.db.DB.Model(workspace).
		Updates(workspace).Error
}

func (_i *workspaceRepository) Delete(id uint64) error {
	return _i.db.DB.Model(&schema.Workspace{}).
		Where("id = ?", id).
		Delete(&schema.Workspace{}).Error
}

func (_i *workspaceRepository) CheckNameExists(name string, ownerID uint64, excludeID uint64) bool {
	var count int64
	query := _i.db.DB.Model(&schema.Workspace{}).
		Where("name = ? AND owner_id = ?", name, ownerID)

	// Exclude current workspace jika sedang update
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	query.Count(&count)
	return count > 0
}
