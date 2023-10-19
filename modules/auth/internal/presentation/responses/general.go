package responses

type GeneralResponse struct {
	Code    int         `json:"code" example:"123"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data,omitempty"`
}
