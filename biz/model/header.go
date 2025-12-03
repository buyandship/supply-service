package model

type Header struct {
	XRequestID string `header:"X-Request-ID" validate:"required" example:"adc3d4cb-d4d0-4a74-9908-5b95bee4d62b"`
	Timestamp  string `header:"timestamp" validate:"required" example:"1762155727995"`
	Hmac       string `header:"hmac" validate:"required" example:"a7062ca18a39e7ec551499958684745f3bd28227c7ae52b5246492c738fa7989"`
}
