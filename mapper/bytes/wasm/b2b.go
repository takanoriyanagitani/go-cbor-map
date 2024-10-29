package wasm

import (
	"context"

	bb "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes"
)

const (
	SetInputSize   string = "set_input_size"
	GetOutEstimate string = "get_out_estimate"
	SetOutputSize  string = "set_output_size"
	Converter      string = "converter"
	OffsetI        string = "offset_i"
	OffsetO        string = "offset_o"
)

type ConverterFactory func(context.Context, []byte) (bb.BytesToBytes, error)
