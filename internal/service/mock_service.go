package service

import (
	"ImportAndSearchCsvFile/pkg/models"
	"io"
)

type MockService struct {
	ImportUsersFn    func(reader io.Reader) error
	GetUserByEmailFn func(email string) (models.User, bool)
}

func (m *MockService) ImportUsers(reader io.Reader) error {
	if m.ImportUsersFn != nil {
		return m.ImportUsersFn(reader)
	}
	return nil
}

func (m *MockService) GetUserByEmail(email string) (models.User, bool) {
	if m.GetUserByEmailFn != nil {
		return m.GetUserByEmailFn(email)
	}
	return models.User{}, false
}
