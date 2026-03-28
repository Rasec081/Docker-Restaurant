package handlers

import "errors"

var errTest = errors.New("test error")

type mockKeycloakService struct {
	createErr error
}

func (m *mockKeycloakService) CreateUser(username, email, password, role string) (string, error) {
	if m.createErr != nil {
		return "", m.createErr
	}
	return "user-id", nil
}
