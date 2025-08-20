package mappers

import (
	"bifur.app/core/internal/adapters/postgres/dbmodels"
	"bifur.app/core/internal/domain"
)

type UserMapper struct {
}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToDbModel(user *domain.User) *dbmodels.User {
	return &dbmodels.User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func (m *UserMapper) ToDomain(u *dbmodels.User) *domain.User {
	return &domain.User{
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
