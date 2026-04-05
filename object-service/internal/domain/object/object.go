package object

import (
	"github.com/chishkin-afk/posted/object-service/pkg/errs"
	"github.com/google/uuid"
)

type Object struct {
	id       uuid.UUID
	filename Filename
	body     []byte
}

func New(filename Filename, body []byte) (*Object, error) {
	filename.Norm()
	if !filename.IsValid() {
		return nil, errs.ErrInvalidFilename
	}

	if len(body) < 1 || len(body) > 4*1024*1024 {
		return nil, errs.ErrInvalidBody
	}

	return &Object{
		id:       uuid.New(),
		filename: filename,
		body:     body,
	}, nil
}

func (o *Object) ID() uuid.UUID {
	return o.id
}

func (o *Object) Filename() Filename {
	return o.filename
}

func (o *Object) Body() []byte {
	return o.body
}
