package cbor2iter

import (
	"context"
	"iter"

	cm "github.com/takanoriyanagitani/go-cbor-map"
)

type CborToIter func() iter.Seq[[]any]

func (c CborToIter) ToArraySource() cm.ArraySource {
	return func(_ context.Context) iter.Seq[[]any] {
		return c()
	}
}
