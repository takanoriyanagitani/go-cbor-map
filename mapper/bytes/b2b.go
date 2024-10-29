package bytes2bytes

import (
	"context"
	"errors"
	"log"

	cm "github.com/takanoriyanagitani/go-cbor-map/mapper"
)

var (
	ErrInvalidInput error = errors.New("invalid input type")
)

type BytesToBytes interface {
	Convert(ctx context.Context, input []byte) (output []byte, e error)
}

type BytesToBytesFn func(context.Context, []byte) ([]byte, error)

func (f BytesToBytesFn) ToMapper() cm.Mapper {
	return func(ctx context.Context, o cm.Original) (cm.Mapd, error) {
		switch input := o.(type) {
		case []byte:
			return f(ctx, input)
		case string:
			return f(ctx, []byte(input))
		default:
			log.Printf("unexpected input: %v\n", input)
			return nil, ErrInvalidInput
		}
	}
}
