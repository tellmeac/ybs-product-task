package json

type BadRequestResponse struct {
	Message string `json:"message,omitempty"`
}

type NotFoundResponse struct{}
