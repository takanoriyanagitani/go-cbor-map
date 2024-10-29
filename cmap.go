package cbor2mapd

import (
	"context"
	"iter"
)

type Original []any
type Mapd []any

type ArrayToMapd func(context.Context, Original) (Mapd, error)

type ArraySource func(context.Context) iter.Seq[[]any]
type ArrayOutput func(context.Context, []any) error

type OutputMapd struct {
	ArraySource
	ArrayToMapd
	ArrayOutput
}

func (o OutputMapd) MapAll(ctx context.Context) error {
	var i iter.Seq[[]any] = o.ArraySource(ctx)
	for arr := range i {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		mapd, e := o.ArrayToMapd(ctx, arr)
		if nil != e {
			return e
		}

		e = o.ArrayOutput(ctx, mapd)
		if nil != e {
			return e
		}
	}
	return nil
}
