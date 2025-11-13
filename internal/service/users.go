package service

import (
	"cruder/internal/model"
	"cruder/internal/repository"
	"cruder/pkg/validation"
)

type UserService interface {
	GetAll() ([]model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByID(id int64) (*model.User, error)
	GetByUUID(uuid string) (*model.User, error)
	Create(req *model.CreateUserRequest) (*model.User, error)
	Update(uuid string, req *model.UpdateUserRequest) (*model.User, error)
	Delete(uuid string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll() ([]model.User, error) {
	return s.repo.GetAll()
}

func (s *userService) GetByUsername(username string) (*model.User, error) {
	if err := validation.ValidateUsername(username); err != nil {
		return nil, err
	}
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, model.NewValidationError("user not found")
	}
	return user, nil
}

func (s *userService) GetByID(id int64) (*model.User, error) {
	if err := validation.ValidateID(id); err != nil {
		return nil, err
	}
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, model.NewValidationError("user not found")
	}
	return user, nil
}

func (s *userService) GetByUUID(uuid string) (*model.User, error) {
	if err := validation.ValidateUUID(uuid); err != nil {
		return nil, err
	}
	user, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, model.NewValidationError("user not found")
	}
	return user, nil
}

func (s *userService) Create(req *model.CreateUserRequest) (*model.User, error) {
	if err := validation.ValidateCreateUserInput(req.Username, req.Email, req.FullName); err != nil {
		return nil, err
	}
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
	}
	return s.repo.Create(user)
}

func (s *userService) Update(uuid string, req *model.UpdateUserRequest) (*model.User, error) {
	if err := validation.ValidateUUID(uuid); err != nil {
		return nil, err
	}
	if err := validation.ValidateUpdateUserInput(req.Username, req.Email, req.FullName); err != nil {
		return nil, err
	}
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
	}
	updatedUser, err := s.repo.Update(uuid, user)
	if err != nil {
		return nil, err
	}
	if updatedUser == nil {
		return nil, model.NewValidationError("user not found")
	}
	return updatedUser, nil
}

func (s *userService) Delete(uuid string) error {
	if err := validation.ValidateUUID(uuid); err != nil {
		return err
	}
	return s.repo.Delete(uuid)
}
