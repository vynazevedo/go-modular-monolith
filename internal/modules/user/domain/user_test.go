package domain

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		userName    string
		expectError bool
	}{
		{
			name:        "Valid user",
			email:       "usuario@teste.com",
			userName:    "Test User",
			expectError: false,
		},
		{
			name:        "Empty email",
			email:       "",
			userName:    "Test User",
			expectError: true,
		},
		{
			name:        "Empty name",
			email:       "usuario@teste.com",
			userName:    "",
			expectError: true,
		},
		{
			name:        "Empty email and name",
			email:       "",
			userName:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.userName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user.ID() == "" {
				t.Errorf("Expected non-empty ID")
			}

			if user.Email() != tt.email {
				t.Errorf("Expected email %s, got %s", tt.email, user.Email())
			}

			if user.Name() != tt.userName {
				t.Errorf("Expected name %s, got %s", tt.userName, user.Name())
			}

			if user.Status() != StatusActive {
				t.Errorf("Expected status %s, got %s", StatusActive, user.Status())
			}

			now := time.Now()
			diff := now.Sub(user.CreatedAt())
			if diff < 0 || diff > time.Second {
				t.Errorf("CreatedAt time is not close to current time. Diff: %v", diff)
			}
		})
	}
}

func TestReconstructUser(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		email       string
		userName    string
		status      string
		createdAt   time.Time
		expectError bool
	}{
		{
			name:        "Valid user",
			id:          "user-123",
			email:       "usuario@teste.com",
			userName:    "Test User",
			status:      "active",
			createdAt:   time.Now(),
			expectError: false,
		},
		{
			name:        "Empty ID",
			id:          "",
			email:       "usuario@teste.com",
			userName:    "Test User",
			status:      "active",
			createdAt:   time.Now(),
			expectError: true,
		},
		{
			name:        "Empty email",
			id:          "user-123",
			email:       "",
			userName:    "Test User",
			status:      "active",
			createdAt:   time.Now(),
			expectError: true,
		},
		{
			name:        "Empty name",
			id:          "user-123",
			email:       "usuario@teste.com",
			userName:    "",
			status:      "active",
			createdAt:   time.Now(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := ReconstructUser(tt.id, tt.email, tt.userName, tt.status, tt.createdAt)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user.ID() != tt.id {
				t.Errorf("Expected ID %s, got %s", tt.id, user.ID())
			}

			if user.Email() != tt.email {
				t.Errorf("Expected email %s, got %s", tt.email, user.Email())
			}

			if user.Name() != tt.userName {
				t.Errorf("Expected name %s, got %s", tt.userName, user.Name())
			}

			if user.Status().String() != tt.status {
				t.Errorf("Expected status %s, got %s", tt.status, user.Status())
			}

			if !user.CreatedAt().Equal(tt.createdAt) {
				t.Errorf("Expected createdAt %v, got %v", tt.createdAt, user.CreatedAt())
			}
		})
	}
}

func TestUpdateName(t *testing.T) {
	user, _ := NewUser("usuario@teste.com", "Original Name")

	tests := []struct {
		name        string
		newName     string
		expectError bool
	}{
		{
			name:        "Valid name update",
			newName:     "New Name",
			expectError: false,
		},
		{
			name:        "Empty name",
			newName:     "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.UpdateName(tt.newName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if user.Name() != tt.newName {
				t.Errorf("Expected name %s, got %s", tt.newName, user.Name())
			}
		})
	}
}

func TestActivateDeactivate(t *testing.T) {
	t.Run("Activate inactive user", func(t *testing.T) {
		user, _ := NewUser("usuario@teste.com", "Test User")
		user.Deactivate()

		if user.Status() != StatusInactive {
			t.Errorf("Expected status %s after deactivate, got %s", StatusInactive, user.Status())
		}

		user.Activate()

		if user.Status() != StatusActive {
			t.Errorf("Expected status %s after activate, got %s", StatusActive, user.Status())
		}
	})

	t.Run("Deactivate active user", func(t *testing.T) {
		user, _ := NewUser("usuario@teste.com", "Test User")

		if user.Status() != StatusActive {
			t.Errorf("Expected initial status %s, got %s", StatusActive, user.Status())
		}

		user.Deactivate()

		if user.Status() != StatusInactive {
			t.Errorf("Expected status %s after deactivate, got %s", StatusInactive, user.Status())
		}
	})
}

func TestStatusString(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{StatusActive, "active"},
		{StatusInactive, "inactive"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.status.String(); got != tt.want {
				t.Errorf("Status.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
