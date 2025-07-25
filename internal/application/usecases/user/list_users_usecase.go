package user

import (
	"context"
	"user-management/internal/application/dto"
	"user-management/internal/application/mappers"
	"user-management/internal/domain/entities"
	"user-management/internal/domain/interfaces/repositories"
)

type ListUsersUseCase struct {
	repo repositories.IUserRepository
}

func NewListUsersUseCase(repo repositories.IUserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{repo: repo}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, input *dto.ListUserQueryParam) (*dto.UserListResponseDTO, error) {
	var users []*entities.User
	var total int64
	var err error

	// Se há um termo de busca, usa a busca filtrada
	if input.Search != "" {
		users, err = uc.repo.Search(ctx, input.Search, input.Page, input.PerPage)
		if err != nil {
			return nil, err
		}

		total, err = uc.repo.CountSearch(ctx, input.Search)
		if err != nil {
			return nil, err
		}
	} else {
		// Caso contrário, lista todos os usuários
		users, err = uc.repo.List(ctx, input.Page, input.PerPage)
		if err != nil {
			return nil, err
		}

		total, err = uc.repo.Count(ctx)
		if err != nil {
			return nil, err
		}
	}

	userListDTO := mappers.ToUserListResponseDTO(users, total, input.Page, input.PerPage)
	return userListDTO, nil
}
