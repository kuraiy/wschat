package dto

type ChangeUsernameDTO struct {
	Username string `json:"username" binding:"required,username"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,password"`
}

type GetMeDTO struct {
	ID       int64
	Username string
}
