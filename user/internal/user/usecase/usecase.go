package usecase

import (
	"context"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
)

// UserUseCase
type UserUseCase struct {
	userPGRepo user.PGRepository
}

// NewUserUseCase
func NewUserUseCase(userPGRepo user.PGRepository) *UserUseCase {
	return &UserUseCase{userPGRepo: userPGRepo}
}

func (u *UserUseCase) Register(ctx context.Context) {
	panic("implement me")
}
