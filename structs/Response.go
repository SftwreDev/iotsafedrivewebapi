package structs

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
