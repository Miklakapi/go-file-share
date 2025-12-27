package dto

type CreateRoomRequest struct {
	Password string `json:"password" form:"password" binding:"required"`
	Lifespan int    `json:"lifespan" form:"lifespan"`
}
