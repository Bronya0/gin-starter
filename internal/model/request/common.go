package request

type Page struct {
	Page int `form:"page,default=1"`
	Size int `form:"size,default=10"`
}
