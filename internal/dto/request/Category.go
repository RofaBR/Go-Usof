package request

type CreateCategory struct {
	Title string `json:"title" binding:"required"`
	Desc  string `json:"description"`
}
