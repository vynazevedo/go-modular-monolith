package http

type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Total int            `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
