package models

type ErrorResponse struct {
	ErrorCode    int    `json:"error_code" validate:"required"`
	ErrorMessage string `json:"error_message" validate:"required"`
}

type LaunchURL struct {
	GameURL string `json:"game_url" validate:"required"`
}

type SuccessResponse struct {
	Status  int         `json:"status" validate:"required"`
	Message string      `json:"message" validate:"required"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseMessage struct {
	Status  int         `json:"status" validate:"required"`
	Message interface{} `json:"message" validate:"required"`
}

type PaginationFilters struct {
	Page    int64  `json:"page" form:"page" query:"page"`
	PerPage int64  `json:"per_page" form:"per_page" query:"per_page"`
	Sort    string `json:"sort" form:"sort" query:"sort"`
	Start   string `json:"start" form:"start" query:"start"`
	End     string `json:"end" form:"end" query:"end"`
	Period  int64  `json:"period" form:"period" query:"period"`
}

type Pagination struct {
	Total       int         `json:"total"`
	PerPage     int         `json:"per_page"`
	NextPageUrl string      `json:"next_page_url"`
	PrevPageUrl string      `json:"prev_page_url"`
	CurrentPage int         `json:"current_page"`
	LastPage    int         `json:"last_page"`
	From        int         `json:"from"`
	To          int         `json:"to"`
	Data        interface{} `json:"data"`
}
