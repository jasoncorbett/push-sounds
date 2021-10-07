package libraries

import (
	"fmt"
	"path/filepath"
	"testing"
)

type mockUserConfigBase struct {
	Called bool
	Base string
	Error error
	orig func() (string, error)
}

func (m *mockUserConfigBase) fakeUserConfigBase() (string, error) {
	m.Called = true
	return m.Base, m.Error
}

func (m *mockUserConfigBase) Replace() {
	m.orig = GetUserConfigBase
	GetUserConfigBase = m.fakeUserConfigBase
}

func (m *mockUserConfigBase) Restore() {
	GetUserConfigBase = m.orig
}

func TestLocationDefaultHappyPath(t *testing.T) {
	mock := &mockUserConfigBase{
		Called: false,
		Base: "base",
		Error: nil,
	}
	mock.Replace()
	defer mock.Restore()
	expectedPath := filepath.Join(mock.Base, "push-sounds")

	actualPath := GetLocationDefault()
	if !mock.Called {
		t.Errorf("Mock config base function never called!")
	}
	if actualPath != expectedPath {
		t.Fatalf("Expected default path to be '%s', it was '%s'", expectedPath, actualPath)
	}
}
func TestLocationDefaultErrorWithBase(t *testing.T) {
	mock := &mockUserConfigBase{
		Called: false,
		Base: "something ridiculous",
		Error: fmt.Errorf("error getting base config path"),
	}
	mock.Replace()
	defer mock.Restore()
	expectedPath := filepath.Join(".", "push-sounds")

	actualPath := GetLocationDefault()
	if !mock.Called {
		t.Errorf("Mock config base function never called!")
	}
	if actualPath != expectedPath {
		t.Fatalf("Expected default path to be '%s', it was '%s'", expectedPath, actualPath)
	}
}
