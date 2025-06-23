package app

type GetUserQuery struct {
	ID string
}

type GetUserByEmailQuery struct {
	Email string
}

type ListUsersQuery struct {
	Page  int
	Limit int
}
