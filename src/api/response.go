package api

type Response struct {
	Data any `json:"data"`
}

func OK(data any) Response {
	return Response{data}
}
