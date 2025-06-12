package request

type Pagination struct {
	Page  int    `form:"page,default=1" binding:"gte=1"`
	Limit int    `form:"limit,default=10" binding:"gte=1,lte=100"`
	Sort  string `form:"sort"`
}
