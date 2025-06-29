package request

type (
	UpdateMenuAvailabilityRequest struct {
		IsAvailable *bool `json:"is_available" form:"is_available" binding:"required"`
	}
)
