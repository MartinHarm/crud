package service

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"cruder/internal/model"
)

type MockUserRepository struct {
	getAllFunc        func() ([]model.User, error)
	getByUsernameFunc func(username string) (*model.User, error)
	getByIDFunc       func(id int64) (*model.User, error)
	getByUUIDFunc     func(uuid string) (*model.User, error)
	createFunc        func(user *model.User) (*model.User, error)
	updateFunc        func(uuid string, user *model.User) (*model.User, error)
	deleteFunc        func(uuid string) error
}

func (m *MockUserRepository) GetAll() ([]model.User, error) {
	return m.getAllFunc()
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	return m.getByUsernameFunc(username)
}

func (m *MockUserRepository) GetByID(id int64) (*model.User, error) {
	return m.getByIDFunc(id)
}

func (m *MockUserRepository) GetByUUID(uuid string) (*model.User, error) {
	return m.getByUUIDFunc(uuid)
}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	return m.createFunc(user)
}

func (m *MockUserRepository) Update(uuid string, user *model.User) (*model.User, error) {
	return m.updateFunc(uuid, user)
}

func (m *MockUserRepository) Delete(uuid string) error {
	return m.deleteFunc(uuid)
}

func TestGetAll_Success(t *testing.T) {
	users := []model.User{
		{
			ID:        1,
			UUID:      "123e4567-e89b-12d3-a456-426614174000",
			Username:  "jdoe",
			Email:     "jdoe@example.com",
			FullName:  "John Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			UUID:      "223e4567-e89b-12d3-a456-426614174001",
			Username:  "asmith",
			Email:     "asmith@example.com",
			FullName:  "Alice Smith",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo := &MockUserRepository{
		getAllFunc: func() ([]model.User, error) {
			return users, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.GetAll()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 users, got %d", len(result))
	}
	if result[0].Username != "jdoe" {
		t.Errorf("expected username jdoe, got %s", result[0].Username)
	}
}

func TestGetAll_EmptyResult(t *testing.T) {
	mockRepo := &MockUserRepository{
		getAllFunc: func() ([]model.User, error) {
			return []model.User{}, nil
		},
	}

	service := NewUserService(mockRepo)

	// When: Calling GetAll
	result, err := service.GetAll()

	// Then: Should return empty slice without error
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result == nil || len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestGetAll_DatabaseError(t *testing.T) {
	mockRepo := &MockUserRepository{
		getAllFunc: func() ([]model.User, error) {
			return nil, errors.New("database connection failed")
		},
	}

	service := NewUserService(mockRepo)

	// When: Calling GetAll
	result, err := service.GetAll()

	// Then: Should return error
	if err == nil {
		t.Error("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestGetByUsername_Success(t *testing.T) {
	expectedUser := &model.User{
		ID:        1,
		UUID:      "123e4567-e89b-12d3-a456-426614174000",
		Username:  "jdoe",
		Email:     "jdoe@example.com",
		FullName:  "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockUserRepository{
		getByUsernameFunc: func(username string) (*model.User, error) {
			if username == "jdoe" {
				return expectedUser, nil
			}
			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.GetByUsername("jdoe")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Username != "jdoe" {
		t.Errorf("expected username jdoe, got %s", result.Username)
	}
}

func TestGetByUsername_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		getByUsernameFunc: func(username string) (*model.User, error) {
			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.GetByUsername("nonexistent")

	if err == nil {
		t.Error("expected error, got nil")
	}
	if _, ok := err.(*model.ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found', got %s", err.Error())
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestGetByID_Success(t *testing.T) {
	expectedUser := &model.User{
		ID:        1,
		UUID:      "123e4567-e89b-12d3-a456-426614174000",
		Username:  "jdoe",
		Email:     "jdoe@example.com",
		FullName:  "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockUserRepository{
		getByIDFunc: func(id int64) (*model.User, error) {
			if id == 1 {
				return expectedUser, nil
			}
			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.GetByID(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
}

func TestGetByID_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		getByIDFunc: func(id int64) (*model.User, error) {
			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.GetByID(999)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if _, ok := err.(*model.ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestGetByUUID_Success(t *testing.T) {
	expectedUser := &model.User{
		ID:        1,
		UUID:      "123e4567-e89b-12d3-a456-426614174000",
		Username:  "jdoe",
		Email:     "jdoe@example.com",
		FullName:  "John Doe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockUserRepository{
		getByUUIDFunc: func(uuid string) (*model.User, error) {
			if uuid == "123e4567-e89b-12d3-a456-426614174000" {
				return expectedUser, nil
			}
			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.GetByUUID("123e4567-e89b-12d3-a456-426614174000")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.UUID != "123e4567-e89b-12d3-a456-426614174000" {
		t.Errorf("expected UUID 123e4567-e89b-12d3-a456-426614174000, got %s", result.UUID)
	}
}

func TestGetByUUID_UserNotFound(t *testing.T) {
	mockRepo := &MockUserRepository{
		getByUUIDFunc: func(uuid string) (*model.User, error) {
			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.GetByUUID("invalid-uuid")

	if err == nil {
		t.Error("expected error, got nil")
	}
	if _, ok := err.(*model.ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestCreate_Success(t *testing.T) {
	req := &model.CreateUserRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		FullName: "New User",
	}

	createdUser := &model.User{
		ID:        3,
		UUID:      "323e4567-e89b-12d3-a456-426614174002",
		Username:  req.Username,
		Email:     req.Email,
		FullName:  req.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockUserRepository{
		createFunc: func(user *model.User) (*model.User, error) {
			return createdUser, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.Create(req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Username != "newuser" {
		t.Errorf("expected username newuser, got %s", result.Username)
	}
	if result.UUID == "" {
		t.Error("expected UUID to be set, got empty string")
	}
}

func TestCreate_DatabaseError(t *testing.T) {
	req := &model.CreateUserRequest{
		Username: "jdoe",
		Email:    "jdoe@example.com",
		FullName: "John Doe",
	}

	mockRepo := &MockUserRepository{
		createFunc: func(user *model.User) (*model.User, error) {
			return nil, errors.New("duplicate username")
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.Create(req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestUpdate_Success(t *testing.T) {
	uuid := "123e4567-e89b-12d3-a456-426614174000"
	req := &model.UpdateUserRequest{
		Username: "jdoe_updated",
		Email:    "jdoe_updated@example.com",
		FullName: "John Doe Updated",
	}

	updatedUser := &model.User{
		ID:        1,
		UUID:      uuid,
		Username:  "jdoe_updated",
		Email:     "jdoe_updated@example.com",
		FullName:  "John Doe Updated",
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now(),
	}

	mockRepo := &MockUserRepository{
		updateFunc: func(u string, user *model.User) (*model.User, error) {
			if u == uuid {
				return updatedUser, nil
			}
			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.Update(uuid, req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result.Username != "jdoe_updated" {
		t.Errorf("expected username jdoe_updated, got %s", result.Username)
	}
	if result.UUID != uuid {
		t.Errorf("expected UUID %s, got %s", uuid, result.UUID)
	}
}

func TestUpdate_UserNotFound(t *testing.T) {
	uuid := "invalid-uuid"
	req := &model.UpdateUserRequest{
		Username: "updated",
	}

	mockRepo := &MockUserRepository{
		updateFunc: func(u string, user *model.User) (*model.User, error) {
			return nil, nil
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.Update(uuid, req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if _, ok := err.(*model.ValidationError); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestUpdate_DatabaseError(t *testing.T) {
	uuid := "123e4567-e89b-12d3-a456-426614174000"
	req := &model.UpdateUserRequest{
		Email: "duplicate@example.com",
	}

	mockRepo := &MockUserRepository{
		updateFunc: func(u string, user *model.User) (*model.User, error) {
			return nil, errors.New("duplicate email")
		},
	}

	service := NewUserService(mockRepo)

	result, err := service.Update(uuid, req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if _, ok := err.(*model.ValidationError); ok {
		t.Errorf("expected database error, got ValidationError")
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}

func TestDelete_Success(t *testing.T) {
	uuid := "123e4567-e89b-12d3-a456-426614174000"

	mockRepo := &MockUserRepository{
		deleteFunc: func(u string) error {
			if u == uuid {
				return nil
			}
			return sql.ErrNoRows
		},
	}

	service := NewUserService(mockRepo)

	err := service.Delete(uuid)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDelete_UserNotFound(t *testing.T) {
	uuid := "invalid-uuid"

	mockRepo := &MockUserRepository{
		deleteFunc: func(u string) error {
			return sql.ErrNoRows
		},
	}

	service := NewUserService(mockRepo)

	// When: Calling Delete with invalid UUID
	err := service.Delete(uuid)

	// Then: Should return sql.ErrNoRows
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err != sql.ErrNoRows {
		t.Errorf("expected sql.ErrNoRows, got %v", err)
	}
}

// TestDelete_DatabaseError tests error handling from database
func TestDelete_DatabaseError(t *testing.T) {
	// Given: Database returns error during deletion
	uuid := "123e4567-e89b-12d3-a456-426614174000"

	mockRepo := &MockUserRepository{
		deleteFunc: func(u string) error {
			return errors.New("database connection failed")
		},
	}

	service := NewUserService(mockRepo)

	// When: Calling Delete
	err := service.Delete(uuid)

	// Then: Should return error
	if err == nil {
		t.Error("expected error, got nil")
	}
}
