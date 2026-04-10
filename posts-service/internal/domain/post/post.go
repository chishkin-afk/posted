package post

import (
	"strings"
	"time"

	"github.com/chishkin-afk/posted/posts-service/pkg/errs"
	"github.com/google/uuid"
)

type Post struct {
	id      uuid.UUID
	ownerID uuid.UUID
	title   string
	body    string

	postedAt  time.Time
	updatedAt time.Time
}

func New(ownerID uuid.UUID, title, body string) (*Post, error) {
	title = strings.TrimSpace(title)
	nTitle := len([]rune(title))
	if nTitle < 3 || nTitle > 64 {
		return nil, errs.ErrInvalidTitle
	}

	body = strings.TrimSpace(body)
	nBody := len([]rune(body))
	if nBody < 3 || nBody > 512 {
		return nil, errs.ErrInvalidBody
	}

	now := time.Now().UTC()
	return &Post{
		id:        uuid.New(),
		ownerID:   ownerID,
		title:     title,
		body:      body,
		postedAt:  now,
		updatedAt: now,
	}, nil
}

func From(
	id, ownerID uuid.UUID,
	title, body string,
	postedAt, updatedAt time.Time,
) (*Post, error) {
	title = strings.TrimSpace(title)
	nTitle := len([]rune(title))
	if nTitle < 3 || nTitle > 64 {
		return nil, errs.ErrInvalidTitle
	}

	body = strings.TrimSpace(body)
	nBody := len([]rune(body))
	if nBody < 3 || nBody > 512 {
		return nil, errs.ErrInvalidBody
	}

	return &Post{
		id:        id,
		ownerID:   ownerID,
		title:     title,
		body:      body,
		postedAt:  postedAt,
		updatedAt: updatedAt,
	}, nil
}

func (p *Post) ChangeTitle(title string) error {
	title = strings.TrimSpace(title)
	nTitle := len([]rune(title))
	if nTitle < 3 || nTitle > 64 {
		return errs.ErrInvalidTitle
	}

	p.title = title
	p.updatedAt = time.Now().UTC()

	return nil
}

func (p *Post) ChangeBody(body string) error {
	body = strings.TrimSpace(body)
	nBody := len([]rune(body))
	if nBody < 3 || nBody > 512 {
		return errs.ErrInvalidBody
	}

	p.body = body
	p.updatedAt = time.Now().UTC()

	return nil
}

func (p *Post) ID() uuid.UUID {
	return p.id
}

func (p *Post) OwnerID() uuid.UUID {
	return p.ownerID
}

func (p *Post) Title() string {
	return p.title
}

func (p *Post) Body() string {
	return p.body
}

func (p *Post) PostedAt() time.Time {
	return p.postedAt
}

func (p *Post) UpdatedAt() time.Time {
	return p.updatedAt
}
