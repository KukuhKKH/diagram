package controller

import (
	"git.dev.siap.id/kukuhkkh/app-diagram/app/module/workspace/service"
	"go.uber.org/fx"
)

// Controller aggregator
type Controller struct {
	Workspace WorkspaceControllerI
}

// NewController
func NewController(workspaceController WorkspaceControllerI) *Controller {
	return &Controller{
		Workspace: workspaceController,
	}
}

var Module = fx.Options(
	fx.Provide(func(workspaceService service.WorkspaceService) WorkspaceControllerI {
		return NewWorkspaceController(workspaceService)
	}),
	fx.Provide(NewController),
)
