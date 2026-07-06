package dto

type LoginOutput struct {
	ID       int64
	Username string

	AccessToken  string
	RefreshToken string

	AccessExp  int
	RefreshExp int
}

type SignUpDTO struct {
	Username string `json:"username" binding:"required,username"`
	Password string `json:"password" binding:"required,password"`
}

type SignInDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
