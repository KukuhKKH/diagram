package request

type CreateWorkspaceRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
}

type UpdateWorkspaceRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	IsPublic    *bool   `json:"is_public" validate:"omitempty"`
}
