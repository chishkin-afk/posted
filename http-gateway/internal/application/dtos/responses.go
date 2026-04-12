package dtos

// Token represents the authentication token response.
// @Description Token represents the authentication token response.
type Token struct {
	Access    string `json:"access" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	AccessTTL int64  `json:"access_ttl" example:"3600"`
}

// User represents the user profile data.
// @Description User represents the user profile data.
type User struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email" example:"user@example.com"`
	Nickname  string `json:"nickname" example:"john_doe"`
	CreatedAt int64  `json:"created_at" example:"1678886400"`
	UpdatedAt int64  `json:"updated_at" example:"1678886400"`
}

// Post represents a single post entity.
// @Description Post represents a single post entity.
type Post struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	OwnerID   string `json:"owner_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Title     string `json:"title" example:"My First Post"`
	Body      string `json:"body" example:"Content of the post goes here."`
	PostedAt  int64  `json:"posted_at" example:"1678886400"`
	UpdatedAt int64  `json:"updated_at" example:"1678886400"`
}

// Posts represents a list of posts.
// @Description Posts represents a list of posts.
type Posts struct {
	Posts []Post `json:"posts"`
}

// ErrMsg represents an error response.
// @Description ErrMsg represents an error response.
type ErrMsg struct {
	Error string `json:"error" example:"invalid request"`
}