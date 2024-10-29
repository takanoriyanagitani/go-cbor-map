package iter2cbor

import (
	"context"

	cm "github.com/takanoriyanagitani/go-cbor-map"
)

type ArrayToCbor func([]any) error

func (a ArrayToCbor) ToArrayOutput() cm.ArrayOutput {
	return func(_ context.Context, mapd []any) error {
		return a(mapd)
	}
}
