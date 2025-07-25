package dto

type CreateUserRequestDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	IsActive bool   `json:"is_active"`
}

type ListUserQueryParam struct {
	Page    int64  `query:"page" default:"1" validate:"min=0"`
	PerPage int64  `query:"per_page" default:"10" validate:"min=1,max=100"`
	Search  string `query:"search" validate:"max=100"`
}

type UserListResponseDTO struct {
	Data []UserResponseDTO `json:"users"`
	Meta Meta              `json:"meta"`
}

type UserResponseDTO struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}
