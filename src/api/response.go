package api

type Response[T any] struct {
	Data T `json:"data"`
}

func OK[T any](data T) Response[T] {
	return Response[T]{data}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(message string) ErrorResponse {
	return ErrorResponse{message}
}
