package api

type Response struct {
	Error *string `json:"error"`
	Data  any     `json:"data"`
}

func OK(data any) Response {
	return Response{nil, data}
}
