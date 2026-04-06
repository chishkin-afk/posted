package object

import (
	"strings"

	"github.com/chishkin-afk/posted/object-service/pkg/errs"
	"github.com/google/uuid"
)

type Object struct {
	id   uuid.UUID
	ext  string
	body []byte
}

func New(id uuid.UUID, ext string, body []byte) (*Object, error) {
	ext = strings.TrimSpace(ext)

	if ext == "" || !strings.HasPrefix(ext, ".") {
		return nil, errs.ErrInvalidExtension
	}

	for i, r := range ext {
		if i == 0 {
			continue
		}
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_') {
			return nil, errs.ErrInvalidExtension
		}
	}

	if len(body) < 1 || len(body) > 10*1024*1024 {
		return nil, errs.ErrInvalidBody
	}

	return &Object{
		id:   id,
		ext:  ext,
		body: body,
	}, nil
}

func (o *Object) ID() uuid.UUID {
	return o.id
}

func (o *Object) Ext() string {
	return o.ext
}

func (o *Object) Body() []byte {
	return o.body
}

func (o *Object) Filename() string {
	return o.id.String() + o.ext
}
