package service

import (
	"errors"
	"fmt"
	"time"

	"git.dev.siap.id/kukuhkkh/app-diagram/app/database/schema"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/repository"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/request"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/response"
	"gorm.io/gorm"
)

// WorkspaceService adalah interface untuk business logic workspace
type WorkspaceService interface {
	CreateWorkspace(userID uint64, req *request.CreateWorkspaceRequest) (*response.WorkspaceResponse, error)
	GetWorkspace(id uint64, userID uint64) (*response.WorkspaceResponse, error)
	ListWorkspaces(userID uint64, page, limit int) (*response.WorkspaceListResponse, error)
	UpdateWorkspace(id uint64, userID uint64, req *request.UpdateWorkspaceRequest) (*response.WorkspaceResponse, error)
	DeleteWorkspace(id uint64, userID uint64) error
}

type workspaceService struct {
	workspaceRepo repository.WorkspaceRepository
}

// NewWorkspaceService instance
func NewWorkspaceService(workspaceRepo repository.WorkspaceRepository) WorkspaceService {
	return &workspaceService{
		workspaceRepo: workspaceRepo,
	}
}

func (_i *workspaceService) CreateWorkspace(userID uint64, req *request.CreateWorkspaceRequest) (*response.WorkspaceResponse, error) {
	// Validasi input
	if req.Name == "" {
		return nil, errors.New("workspace name is required")
	}

	// Check if nama workspace sudah ada untuk user
	if _i.workspaceRepo.CheckNameExists(req.Name, userID, 0) {
		return nil, fmt.Errorf("workspace with name '%s' already exists", req.Name)
	}

	workspace := &schema.Workspace{
		OwnerID:     userID,
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	created, err := _i.workspaceRepo.Create(workspace)
	if err != nil {
		return nil, err
	}

	return _i.toResponse(created), nil
}

func (_i *workspaceService) GetWorkspace(id uint64, userID uint64) (*response.WorkspaceResponse, error) {
	workspace, err := _i.workspaceRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workspace not found")
		}
		return nil, err
	}

	// Validasi ownership: user hanya bisa akses workspace miliknya
	if workspace.OwnerID != userID && !workspace.IsPublic {
		return nil, errors.New("you don't have permission to access this workspace")
	}

	return _i.toResponse(workspace), nil
}

func (_i *workspaceService) ListWorkspaces(userID uint64, page, limit int) (*response.WorkspaceListResponse, error) {
	// Validasi pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	workspaces, err := _i.workspaceRepo.FindByOwnerID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	total, err := _i.workspaceRepo.CountByOwnerID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]response.WorkspaceResponse, 0, len(workspaces))
	for _, ws := range workspaces {
		responses = append(responses, *_i.toResponse(&ws))
	}

	return &response.WorkspaceListResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (_i *workspaceService) UpdateWorkspace(id uint64, userID uint64, req *request.UpdateWorkspaceRequest) (*response.WorkspaceResponse, error) {
	workspace, err := _i.workspaceRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("workspace not found")
		}
		return nil, err
	}

	// Validasi ownership
	if workspace.OwnerID != userID {
		return nil, errors.New("you don't have permission to update this workspace")
	}

	// Update fields jika ada
	if req.Name != nil && *req.Name != "" {
		// Check if new name sudah ada
		if _i.workspaceRepo.CheckNameExists(*req.Name, userID, id) {
			return nil, fmt.Errorf("workspace with name '%s' already exists", *req.Name)
		}
		workspace.Name = *req.Name
	}

	if req.Description != nil {
		workspace.Description = req.Description
	}

	if req.IsPublic != nil {
		workspace.IsPublic = *req.IsPublic
	}

	workspace.UpdatedAt = time.Now()

	if err := _i.workspaceRepo.Update(workspace); err != nil {
		return nil, err
	}

	return _i.toResponse(workspace), nil
}

func (_i *workspaceService) DeleteWorkspace(id uint64, userID uint64) error {
	workspace, err := _i.workspaceRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("workspace not found")
		}
		return err
	}

	// Validasi ownership
	if workspace.OwnerID != userID {
		return errors.New("you don't have permission to delete this workspace")
	}

	if err := _i.workspaceRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

// Helper: convert schema to response
func (_i *workspaceService) toResponse(workspace *schema.Workspace) *response.WorkspaceResponse {
	return &response.WorkspaceResponse{
		ID:          workspace.ID,
		Name:        workspace.Name,
		Description: workspace.Description,
		OwnerID:     workspace.OwnerID,
		IsPublic:    workspace.IsPublic,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}
}
