package components

import (
	"context"
	"io"

	"github.com/a-h/templ"
)

func Render(component templ.Component, ctx context.Context, w io.Writer) error {
	return component.Render(ctx, w)
}
