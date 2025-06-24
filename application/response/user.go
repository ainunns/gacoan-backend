package response

type (
	User struct {
		ID          string `json:"id"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number,omitempty"`
		Role        string `json:"role"`
	}

	UserRegister struct {
		ID          string `json:"id"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number,omitempty"`
		Role        string `json:"role"`
	}
)
