package users

import (
	"context"

	"github.com/Financial-Times/gourmet/example/usersapp/pkg/storage"
)

type Service interface {
	RegisterUser(ctx context.Context, c *User) (*User, error)
	UnregisterUser(ctx context.Context, cID int) error
	GetAllUsers(ctx context.Context, opts *storage.QueryOptions) ([]User, error)
	GetUserByID(ctx context.Context, cID int) (User, error)
}

type User struct {
	UserID      int    `json:"userId" db:"cid" goqu:"skipinsert"`
	FirstName   string `json:"firstName" db:"first_name"`
	LastName    string `json:"lastName" db:"last_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Created     int64  `json:"created"`
	LastUpdated int64  `json:"lastUpdated" db:"last_updated"`
}

type UserService struct {
	userRepo Repository
}

func NewUserService(repo Repository) Service {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, c *User) (*User, error) {
	return s.userRepo.AddUser(c)
}

func (s *UserService) UnregisterUser(ctx context.Context, cID int) error {
	return s.userRepo.RemoveUser(cID)
}

func (s *UserService) GetAllUsers(ctx context.Context, opts *storage.QueryOptions) ([]User, error) {
	return s.userRepo.FindAllUsers(opts)
}

func (s *UserService) GetUserByID(ctx context.Context, cID int) (User, error) {
	return s.userRepo.FindUserByID(cID)
}
