package request

type (
	CalculateTotalPrice struct {
		Orders []Order `json:"orders" form:"orders" binding:"required"`
	}
)
