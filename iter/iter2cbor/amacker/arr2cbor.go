package arr2cbor

import (
	"io"

	ca "github.com/fxamacker/cbor/v2"

	ac "github.com/takanoriyanagitani/go-cbor-map/iter/iter2cbor"
)

type ArrToCbor struct {
	*ca.Encoder
}

func (a ArrToCbor) Encode(arr []any) error { return a.Encoder.Encode(arr) }

func (a ArrToCbor) AsArrayToCbor() ac.ArrayToCbor { return a.Encode }

func ArrToCborNew(wtr io.Writer) ArrToCbor {
	return ArrToCbor{Encoder: ca.NewEncoder(wtr)}
}
