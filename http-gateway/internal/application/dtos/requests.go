package dtos

// RegisterRequest represents the payload for user registration.
// @Description RegisterRequest represents the payload for user registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"strongpassword123"`
	Nickname string `json:"nickname" validate:"required" example:"john_doe"`
}

// LoginRequest represents the payload for user login.
// @Description LoginRequest represents the payload for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required" example:"strongpassword123"`
}

// UpdateUserRequest represents the payload for updating user profile.
// @Description UpdateUserRequest represents the payload for updating user profile.
type UpdateUserRequest struct {
	Updates struct {
		Email    string `json:"email" example:"newemail@example.com"`
		Nickname string `json:"nickname" example:"new_nickname"`
	} `json:"updates"`
	UpdateMask []string `json:"update_mask" example:"[\"email\",\"nickname\"]"`
}

// CreatePostRequest represents the payload for creating a new post.
// @Description CreatePostRequest represents the payload for creating a new post.
type CreatePostRequest struct {
	Title string `json:"title" validate:"required" example:"My First Post"`
	Body  string `json:"body" validate:"required" example:"Content of the post goes here."`
}

// UpdatePostRequest represents the payload for updating an existing post.
// @Description UpdatePostRequest represents the payload for updating an existing post.
type UpdatePostRequest struct {
	PostID  string `json:"post_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Updates struct {
		Title string `json:"title" example:"Updated Title"`
		Body  string `json:"body" example:"Updated content."`
	} `json:"updates"`
	UpdateMask []string `json:"update_mask" example:"[\"title\"]"`
}