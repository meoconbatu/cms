package users

import "testing"

var username, password = "test", "test"

func TestNewUser(t *testing.T) {
	err := NewUser(username, password)

	if err != nil {
		t.Errorf("Failed to create new user: %s\n", err.Error())
	}
}
func TestNewUserExists(t *testing.T) {
	err := NewUser(username, password)
	if err != ErrUserAlreadyExists {
		t.Errorf("Expect an error. got=%s, want=%s", err.Error(), ErrUserAlreadyExists.Error())
	}
}
func TestAuthenticateUser(t *testing.T) {
	err := AuthenticateUser(username, password)
	if err != nil {
		t.Errorf("Expect no error. got=%s\n", err.Error())
	}
}
func TestOverrideOldPassword(t *testing.T) {
	newPassword := "newtest"
	err := OverrideOldPassword(username, newPassword)
	if err != nil {
		t.Fatalf("Failed to override old password: %s\n", err.Error())
	}
	err = AuthenticateUser(username, newPassword)
	if err != nil {
		t.Errorf("Expect no error. got=%s\n", err.Error())
	}
}
