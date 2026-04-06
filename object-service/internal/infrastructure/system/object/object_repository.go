package objectsystem

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chishkin-afk/posted/object-service/internal/domain/object"
	"github.com/chishkin-afk/posted/object-service/internal/infrastructure/config"
	"github.com/chishkin-afk/posted/object-service/pkg/errs"
)

type objectSystemRepository struct {
	cfg *config.Config
}

func New(cfg *config.Config) *objectSystemRepository {
	return &objectSystemRepository{
		cfg: cfg,
	}
}

func (osr *objectSystemRepository) Save(ctx context.Context, object *object.Object) error {
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		prefix := object.ID().String()[:2]
		fullPath := filepath.Join(osr.cfg.Storage.BasePath, prefix, object.Filename())

		file, err := os.Create(fullPath)
		if err != nil {
			errCh <- fmt.Errorf("failed to create file: %w", err)
			return
		}
		defer file.Close()

		if _, err := file.Write(object.Body()); err != nil {
			errCh <- fmt.Errorf("failed to write file body: %w", err)
			return
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err, _ := <-errCh:
		return err
	}
}

func (osr *objectSystemRepository) Delete(ctx context.Context, object *object.Object) error {
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

		prefix := object.ID().String()[:2]
		fullPath := filepath.Join(osr.cfg.Storage.BasePath, prefix, object.Filename())

		if err := os.Remove(fullPath); err != nil {
			errCh <- fmt.Errorf("failed to delete file: %w", err)
			return
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err, _ := <-errCh:
		return err
	}
}

func (osr *objectSystemRepository) Get(ctx context.Context, find *object.Object) (*object.Object, error) {
	errCh := make(chan error, 1)
	bytesCh := make(chan []byte, 1)

	go func() {
		defer close(errCh)
		defer close(bytesCh)

		prefix := find.ID().String()[:2]
		fullPath := filepath.Join(osr.cfg.Storage.BasePath, prefix, find.Filename())

		bytes, err := os.ReadFile(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				errCh <- errs.ErrObjectNotFound
				return
			}

			errCh <- fmt.Errorf("failed to read file: %w", err)
			return
		}

		bytesCh <- bytes
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case bytes := <-bytesCh:
		return object.New(
			find.ID(),
			find.Ext(),
			bytes,
		)
	case err := <-errCh:
		return nil, err
	}
}
