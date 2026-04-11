package dtos

type Token struct {
	Access    string `json:"access"`
	AccessTTL int64  `json:"access_ttl"`
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Post struct {
	ID        string `json:"id"`
	OwnerID   string `json:"owner_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	PostedAt  int64  `json:"posted_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type ErrMsg struct {
	Error string `json:"error"`
}
