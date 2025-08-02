package mappers

import (
	"scheduly.io/core/internal/adapters/postgres/dbmodels"
	"scheduly.io/core/internal/domain"
)

type UserMapper struct {
}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) DomainToDBModel(user *domain.User) *dbmodels.User {
	return &dbmodels.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func (m *UserMapper) DBModelToDomain(u *dbmodels.User) *domain.User {
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
