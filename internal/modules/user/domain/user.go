package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

func (s Status) String() string {
	return string(s)
}

type User struct {
	id        string
	email     string
	name      string
	status    Status
	createdAt time.Time
}

func NewUser(email, name string) (*User, error) {
	if email == "" || name == "" {
		return nil, errors.New("email and name required")
	}

	return &User{
		id:        uuid.New().String(),
		email:     email,
		name:      name,
		status:    StatusActive,
		createdAt: time.Now(),
	}, nil
}

func ReconstructUser(id, email, name, status string, createdAt time.Time) (*User, error) {
	if id == "" || email == "" || name == "" {
		return nil, errors.New("id, email and name required")
	}

	return &User{
		id:        id,
		email:     email,
		name:      name,
		status:    Status(status),
		createdAt: createdAt,
	}, nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Status() Status {
	return u.status
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdateName(name string) error {
	if name == "" {
		return errors.New("name required")
	}
	u.name = name
	return nil
}

func (u *User) Activate() {
	u.status = StatusActive
}

func (u *User) Deactivate() {
	u.status = StatusInactive
}
