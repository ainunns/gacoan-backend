package presentation

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Response
	Data any `json:"data"`
	Meta any `json:"meta,omitempty"`
}

type FailedResponse struct {
	Response
	Error any `json:"error"`
	Data  any `json:"data,omitempty"`
}

type EmptyObj struct{}

func BuildResponseSuccess(message string, data any, meta ...any) SuccessResponse {
	res := SuccessResponse{
		Response: Response{
			Status:  true,
			Message: message,
		},
		Data: data,
	}

	if len(meta) > 0 {
		res.Meta = meta[0]
	}

	return res
}

func BuildResponseFailed(message string, err string, data any) FailedResponse {
	res := FailedResponse{
		Response: Response{
			Status:  false,
			Message: message,
		},
		Error: err,
		Data:  data,
	}
	return res
}
