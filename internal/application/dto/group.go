package dto

type CreateGroupRequestDTO struct {
	Name    string   `json:"name" validate:"required,min=2,max=100"`
	Members []string `json:"members"`
}

type GroupResponseDTO struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}
