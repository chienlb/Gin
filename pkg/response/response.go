package response

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int64       `json:"total"`
}

func Success(message string, data interface{}) Response {
	return Response{
		Status:  200,
		Message: message,
		Data:    data,
	}
}

func Created(message string, data interface{}) Response {
	return Response{
		Status:  201,
		Message: message,
		Data:    data,
	}
}

func BadRequest(message string) Response {
	return Response{
		Status:  400,
		Message: message,
		Error:   message,
	}
}

func Unauthorized(message string) Response {
	return Response{
		Status:  401,
		Message: message,
		Error:   message,
	}
}

func NotFound(message string) Response {
	return Response{
		Status:  404,
		Message: message,
		Error:   message,
	}
}

func InternalServerError(message string) Response {
	return Response{
		Status:  500,
		Message: message,
		Error:   message,
	}
}
