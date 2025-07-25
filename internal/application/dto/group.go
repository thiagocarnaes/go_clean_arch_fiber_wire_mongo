package dto

type CreateGroupRequestDTO struct {
	Name    string   `json:"name" validate:"required,min=2,max=100"`
	Members []string `json:"members"`
}

type ListGroupResponseDTO struct {
	Data []*GroupResponseDTO `json:"groups"`
	Meta Meta                `json:"meta"`
}

type ListGroupQueryParam struct {
	Page    int64 `query:"page" default:"1" validate:"min=0"`
	PerPage int64 `query:"per_page" default:"10" validate:"min=1,max=100"`
}

type GroupResponseDTO struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}
