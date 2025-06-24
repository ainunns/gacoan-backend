package request

type (
	UserRegister struct {
		Email       string `json:"email" form:"email" binding:"required,email"`
		Password    string `json:"password" form:"password" binding:"required,min=8"`
		Name        string `json:"name" form:"name" binding:"required,min=2,max=100"`
		PhoneNumber string `json:"phone_number" form:"phone_number" binding:"omitempty,min=8,max=20"`
	}

	UserLogin struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required,min=8"`
	}
)
