package presenter

type RoleInput struct {
	Name string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}
