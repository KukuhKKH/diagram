package controller

import (
	"strconv"

	"git.dev.siap.id/kukuhkkh/app-diagram/app/middleware"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/request"
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/service"
	"git.dev.siap.id/kukuhkkh/app-diagram/utils/response"
	"github.com/gofiber/fiber/v2"
)

// WorkspaceController
type workspaceController struct {
	workspaceService service.WorkspaceService
}

type WorkspaceControllerI interface {
	CreateWorkspace(c *fiber.Ctx) error
	GetWorkspace(c *fiber.Ctx) error
	ListWorkspaces(c *fiber.Ctx) error
	UpdateWorkspace(c *fiber.Ctx) error
	DeleteWorkspace(c *fiber.Ctx) error
}

func NewWorkspaceController(workspaceService service.WorkspaceService) WorkspaceControllerI {
	return &workspaceController{
		workspaceService: workspaceService,
	}
}

func (_i *workspaceController) CreateWorkspace(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusUnauthorized,
			Messages: response.Messages{"user not authenticated"},
		})
	}

	var req request.CreateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusBadRequest,
			Messages: response.Messages{"invalid request body"},
		})
	}

	// Validasi input
	if err := response.ValidateStruct(req); err != nil {
		return err
	}

	result, err := _i.workspaceService.CreateWorkspace(userID, &req)
	if err != nil {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusBadRequest,
			Messages: response.Messages{err.Error()},
		})
	}

	return response.Resp(c, response.Response{
		Code:     fiber.StatusCreated,
		Messages: response.Messages{"workspace created successfully"},
		Data:     result,
	})
}

// GetWorkspace handler untuk get single workspace
func (_i *workspaceController) GetWorkspace(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusUnauthorized,
			Messages: response.Messages{"user not authenticated"},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusBadRequest,
			Messages: response.Messages{"invalid workspace id"},
		})
	}

	result, err := _i.workspaceService.GetWorkspace(id, userID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "workspace not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "you don't have permission to access this workspace" {
			statusCode = fiber.StatusForbidden
		}
		return response.Resp(c, response.Response{
			Code:     statusCode,
			Messages: response.Messages{err.Error()},
		})
	}

	return response.Resp(c, response.Response{
		Code:     fiber.StatusOK,
		Messages: response.Messages{"workspace retrieved successfully"},
		Data:     result,
	})
}

// ListWorkspaces handler untuk list workspaces
func (_i *workspaceController) ListWorkspaces(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusUnauthorized,
			Messages: response.Messages{"user not authenticated"},
		})
	}

	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	result, err := _i.workspaceService.ListWorkspaces(userID, page, limit)
	if err != nil {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusInternalServerError,
			Messages: response.Messages{err.Error()},
		})
	}

	return response.Resp(c, response.Response{
		Code:     fiber.StatusOK,
		Messages: response.Messages{"workspaces retrieved successfully"},
		Data:     result,
	})
}

// UpdateWorkspace handler untuk update workspace
func (_i *workspaceController) UpdateWorkspace(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusUnauthorized,
			Messages: response.Messages{"user not authenticated"},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusBadRequest,
			Messages: response.Messages{"invalid workspace id"},
		})
	}

	var req request.UpdateWorkspaceRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusBadRequest,
			Messages: response.Messages{"invalid request body"},
		})
	}

	// Validasi input (omitempty fields tidak di-validasi jika kosong)
	if err := response.ValidateStruct(req); err != nil {
		return err
	}

	result, err := _i.workspaceService.UpdateWorkspace(id, userID, &req)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "workspace not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "you don't have permission to update this workspace" {
			statusCode = fiber.StatusForbidden
		} else {
			statusCode = fiber.StatusBadRequest
		}
		return response.Resp(c, response.Response{
			Code:     statusCode,
			Messages: response.Messages{err.Error()},
		})
	}

	return response.Resp(c, response.Response{
		Code:     fiber.StatusOK,
		Messages: response.Messages{"workspace updated successfully"},
		Data:     result,
	})
}

// DeleteWorkspace handler untuk delete workspace
func (_i *workspaceController) DeleteWorkspace(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusUnauthorized,
			Messages: response.Messages{"user not authenticated"},
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return response.Resp(c, response.Response{
			Code:     fiber.StatusBadRequest,
			Messages: response.Messages{"invalid workspace id"},
		})
	}

	err = _i.workspaceService.DeleteWorkspace(id, userID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err.Error() == "workspace not found" {
			statusCode = fiber.StatusNotFound
		} else if err.Error() == "you don't have permission to delete this workspace" {
			statusCode = fiber.StatusForbidden
		}
		return response.Resp(c, response.Response{
			Code:     statusCode,
			Messages: response.Messages{err.Error()},
		})
	}

	return response.Resp(c, response.Response{
		Code:     fiber.StatusOK,
		Messages: response.Messages{"workspace deleted successfully"},
	})
}
