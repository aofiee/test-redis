package domain

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

var (
	Success             = Status{Code: http.StatusOK, Message: []string{"Success"}}
	BadRequest          = Status{Code: http.StatusBadRequest, Message: []string{"Sorry, Not responding because of incorrect syntax"}}
	Unauthorized        = Status{Code: http.StatusUnauthorized, Message: []string{"Sorry, We are not able to process your request. Please try again"}}
	Forbidden           = Status{Code: http.StatusForbidden, Message: []string{"Sorry, Permission denied"}}
	InternalServerError = Status{Code: http.StatusInternalServerError, Message: []string{"Internal Server Error"}}
	ConFlict            = Status{Code: http.StatusBadRequest, Message: []string{"Sorry, Data is conflict"}}
	FieldsPermission    = Status{Code: http.StatusBadRequest, Message: []string{"Sorry, Fields are not able to update"}}
)

type ResponseBody struct {
	Status Status      `json:"status,omitempty"`
	Data   interface{} `json:"data,omitempty"`

	CurrentPage *int   `json:"current_page,omitempty"`
	PerPage     *int   `json:"per_page,omitempty"`
	TotalItem   *int64 `json:"total_item,omitempty"`
}

type Status struct {
	Code    int      `json:"code,omitempty"`
	Message []string `json:"message,omitempty"`
}

// ///////////////////////////////////
type (
	MachineResponse struct {
		ID        *uint           `json:"id,omitempty"`
		Machine   *string         `json:"machine,omitempty"`
		Stock     *int            `json:"stock" form:"stock" query:"stock"`
		CreatedAt *time.Time      `json:"created_at,omitempty"`
		UpdatedAt *time.Time      `json:"updated_at,omitempty"`
		DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
		CreatedBy *string         `json:"created_by,omitempty"`
		UpdatedBy *string         `json:"updated_by,omitempty"`
		DeletedBy *string         `json:"deleted_by,omitempty"`
	}
	MachineListResult struct {
		Machines []MachineResponse

		CurrentPage *int   `json:"current_page,omitempty"`
		PerPage     *int   `json:"per_page,omitempty"`
		TotalItem   *int64 `json:"total_item,omitempty"`
	}
)

/////////////////////////////////////
