package cbor2arr

import (
	"io"
	"iter"

	ca "github.com/fxamacker/cbor/v2"

	ci "github.com/takanoriyanagitani/go-cbor-map/iter/cbor2iter"
)

type CborToArr struct {
	*ca.Decoder
}

func (c CborToArr) ToArrays() iter.Seq[[]any] {
	return func(yield func([]any) bool) {
		var buf []any
		var err error
		for {
			clear(buf)
			buf = buf[:0]

			err = c.Decoder.Decode(&buf)
			if nil != err {
				return
			}

			if !yield(buf) {
				return
			}
		}
	}
}

func (c CborToArr) AsCborToIter() ci.CborToIter { return c.ToArrays }

func CborToArrNew(rdr io.Reader) CborToArr {
	return CborToArr{Decoder: ca.NewDecoder(rdr)}
}
