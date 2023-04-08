package api

type Response struct {
	Data any `json:"data"`
}

func OK(data any) Response {
	return Response{data}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(message string) ErrorResponse {
	return ErrorResponse{message}
}
