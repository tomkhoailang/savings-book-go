package presenter

type RoleInput struct {
	Name string `json:"name" validate:"required,min=2"`
	Description string `json:"description" validate:"required"`
}
