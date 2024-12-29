package api

type DataResponse struct {
	Data        interface{} `json:"data"`
	TotalPages  int         `json:"totalPages"`
	CurrentPage int         `json:"currentPage"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
