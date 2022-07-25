package requests

type UserGetByIdRequest struct {
	Id int64 `uri:"id" binding:"required"`
}
