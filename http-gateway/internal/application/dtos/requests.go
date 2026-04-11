package dtos

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Nickname string `json:"nickname" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Updates struct {
		Email    string `json:"email"`
		Nickname string `json:"nickname"`
	} `json:"updates"`
	UpdateMask []string `json:"update_mask"`
}

type CreatePostRequest struct {
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}

type UpdatePostRequest struct {
	PostID  string `json:"post_id" validate:"required,uuid"`
	Updates struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"updates"`
	UpdateMask []string `json:"update_mask"`
}
