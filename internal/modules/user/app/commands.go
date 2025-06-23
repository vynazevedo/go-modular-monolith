package app

type CreateUserCommand struct {
	Email string
	Name  string
}

type UpdateUserCommand struct {
	ID   string
	Name string
}

type DeleteUserCommand struct {
	ID string
}

type ActivateUserCommand struct {
	ID string
}

type DeactivateUserCommand struct {
	ID string
}
