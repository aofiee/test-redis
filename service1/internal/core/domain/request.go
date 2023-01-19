package domain

type MachineRequest struct {
	ID        *int    `json:"id" form:"id" query:"id"`
	Machine   *string `json:"machine" form:"machine" query:"machine"`
	Stock     *int    `json:"stock" form:"stock" query:"stock"`
	CreatedBy *string `json:"created_by" form:"created_by" query:"created_by"`
	UpdatedBy *string `json:"updated_by" form:"updated_by" query:"updated_by"`
	DeletedBy *string `json:"deleted_by" form:"deleted_by" query:"deleted_by"`
}

/////////////////////////////////////

type QueryMachineRequest struct {
	ID        *int    `json:"id" form:"id" query:"id"`
	Machine   *string `json:"machine" form:"machine" query:"machine"`
	Stock     *int    `json:"stock" form:"stock" query:"stock"`
	CreatedBy *string `json:"created_by" form:"created_by" query:"created_by"`
	UpdatedBy *string `json:"updated_by" form:"updated_by" query:"updated_by"`
	DeletedBy *string `json:"deleted_by" form:"deleted_by" query:"deleted_by"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	DeletedAt *string `json:"deleted_at,omitempty"`

	IDs        []uint      `json:"ids,omitempty" form:"ids" query:"ids"`
	Limit      *int        `json:"limit,omitempty" form:"limit" query:"limit"`
	Page       *int        `json:"page,omitempty" form:"page" query:"page"`
	OrderBy    *string     `json:"order_by,omitempty" form:"order_by" query:"order_by"`
	Asc        *bool       `json:"asc,omitempty" form:"asc" query:"asc"`
	Pagination *Pagination `json:"-"`
	SortMethod *SortMethod `json:"-"`
}

// ///////////////////////////////////
type Pagination struct {
	Limit  int `json:"limit" query:"limit" validate:"gte=-1,lte=100"`
	Offset int `json:"offset" query:"offset"`
}

type SortMethod struct {
	Asc     bool   `json:"asc" query:"asc"`
	OrderBy string `json:"order_by" query:"order_by"`
}
