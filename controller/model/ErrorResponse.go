package model

type TextResponse struct {
	Message     string `json:"err"`
	Description string `json:"description"`
}

func CreateResponseFromError(code string, err error) TextResponse {
	return CreatePlainResponse(code, err.Error())
}

func CreatePlainResponse(msg string, desc string) TextResponse {
	return TextResponse{Message: msg, Description: desc}
}
