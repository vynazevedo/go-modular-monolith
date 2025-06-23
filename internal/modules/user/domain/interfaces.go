// Package domain provides the interfaces for user-related operations.
package domain

import (
	"context"
)

type UserInfo struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type UserQueryService interface {
	GetUserInfo(ctx context.Context, id string) (*UserInfo, error)
	GetUserByEmail(ctx context.Context, email string) (*UserInfo, error)
}
