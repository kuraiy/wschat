package dto

type ChangeUsernameDTO struct {
	Username string `json:"username" binding:"required,username"`
}

type GetMeDTO struct {
	ID       int64
	Username string
}
