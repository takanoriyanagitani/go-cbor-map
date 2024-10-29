package b2bo0

import (
	"context"
	"errors"
	"iter"
	"maps"

	w0 "github.com/tetratelabs/wazero"
	wa "github.com/tetratelabs/wazero/api"

	mb "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes"
	bw "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes/wasm"
)

var (
	ErrUnableToSetInputLength error = errors.New("unable to set input length")
	ErrUnableToSetInput       error = errors.New("unable to set input data")

	ErrUnableToSetOutputLength   error = errors.New("unable to set output length")
	ErrUnableToGetOutput         error = errors.New("unable to get output data")
	ErrUnableToGetOutputEstimate error = errors.New("unable to get output size estimate")

	ErrUnableToConvert   error = errors.New("unable to convert")
	ErrUnableToGetOffset error = errors.New("unable to get offset")

	ErrInvalidModule error = errors.New("invalid module")

	ErrConvFailure error = errors.New("unable to convert")
)

type BytesToBytesRaw struct {
	SetInputSize   wa.Function
	GetOutEstimate wa.Function
	SetOutputSize  wa.Function
	Converter      wa.Function
	OffsetI        wa.Function
	OffsetO        wa.Function
	wa.Memory
}

func (b BytesToBytesRaw) Validate() (valid bool) {
	oks := []bool{
		nil != b.SetInputSize,
		nil != b.GetOutEstimate,
		nil != b.SetOutputSize,
		nil != b.Converter,
		nil != b.OffsetI,
		nil != b.OffsetO,
		nil != b.Memory,
	}
	for _, ok := range oks {
		var ng bool = !ok
		if ng {
			return false
		}
	}
	return true
}

func (b BytesToBytesRaw) SetInputLength(ctx context.Context, sz uint32) error {
	var encoded uint64 = wa.EncodeU32(sz)
	results, e := b.SetInputSize.Call(ctx, encoded)
	switch {
	case nil != e:
		return e
	case 1 == len(results):
		var i int32 = wa.DecodeI32(results[0])
		if i < 0 {
			return ErrUnableToSetInputLength
		}
		return nil
	default:
		return ErrUnableToSetInputLength
	}
}

func (b BytesToBytesRaw) GetOutputEstimate(ctx context.Context) (uint32, error) {
	results, e := b.GetOutEstimate.Call(ctx)
	switch {
	case nil != e:
		return 0, e
	case 1 == len(results):
		var i int32 = wa.DecodeI32(results[0])
		if i < 0 {
			return 0, ErrUnableToGetOutputEstimate
		}
		return uint32(i), nil
	default:
		return 0, ErrUnableToGetOutputEstimate
	}
}

func (b BytesToBytesRaw) SetOutputLength(ctx context.Context, sz uint32) error {
	var encoded uint64 = wa.EncodeU32(sz)
	results, e := b.SetOutputSize.Call(ctx, encoded)
	switch {
	case nil != e:
		return e
	case 1 == len(results):
		var i int32 = wa.DecodeI32(results[0])
		if i < 0 {
			return ErrUnableToSetOutputLength
		}
		return nil
	default:
		return ErrUnableToSetOutputLength
	}
}

func (b BytesToBytesRaw) Convert(ctx context.Context) (uint32, error) {
	results, e := b.Converter.Call(ctx)
	switch {
	case nil != e:
		return 0, e
	case 1 == len(results):
		var i int32 = wa.DecodeI32(results[0])
		if i < 0 {
			return 0, ErrConvFailure
		}
		return uint32(i), nil
	default:
		return 0, ErrUnableToConvert
	}
}

func (b BytesToBytesRaw) GetOffsetI(ctx context.Context) (uint32, error) {
	results, e := b.OffsetI.Call(ctx)
	switch {
	case nil != e:
		return 0, e
	case 1 == len(results):
		var i int32 = wa.DecodeI32(results[0])
		if i < 0 {
			return 0, ErrUnableToGetOffset
		}
		return wa.DecodeU32(results[0]), nil
	default:
		return 0, ErrUnableToGetOffset
	}
}

func (b BytesToBytesRaw) GetOffsetO(ctx context.Context) (uint32, error) {
	results, e := b.OffsetO.Call(ctx)
	switch {
	case nil != e:
		return 0, e
	case 1 == len(results):
		var i int32 = wa.DecodeI32(results[0])
		if i < 0 {
			return 0, ErrUnableToGetOffset
		}
		return wa.DecodeU32(results[0]), nil
	default:
		return 0, ErrUnableToGetOffset
	}
}

func (b BytesToBytesRaw) SetInput(ctx context.Context, i []byte) error {
	oi, e := b.GetOffsetI(ctx)
	if nil != e {
		return e
	}

	var ok bool = b.Memory.Write(oi, i)
	switch ok {
	case true:
		return nil
	default:
		return ErrUnableToSetInput
	}
}

func (b BytesToBytesRaw) GetOutput(ctx context.Context, sz uint32) ([]byte, error) {
	oo, e := b.GetOffsetO(ctx)
	if nil != e {
		return nil, e
	}
	out, ok := b.Memory.Read(oo, sz)
	switch ok {
	case true:
		return out, nil
	default:
		return nil, ErrUnableToGetOutput
	}
}

func (b BytesToBytesRaw) Map(
	ctx context.Context,
	input []byte,
) (output []byte, e error) {
	var isz uint32 = uint32(len(input))
	e = b.SetInputLength(ctx, isz)
	if nil != e {
		return nil, e
	}

	osz, e := b.GetOutputEstimate(ctx)
	if nil != e {
		return nil, e
	}

	e = b.SetOutputLength(ctx, osz)
	if nil != e {
		return nil, e
	}

	e = b.SetInput(ctx, input)
	if nil != e {
		return nil, e
	}

	csz, e := b.Convert(ctx)
	if nil != e {
		return nil, e
	}

	return b.GetOutput(ctx, csz)
}

type Converter struct {
	w0.CompiledModule
	wa.Module
	BytesToBytesRaw
}

func (c Converter) Close(ctx context.Context) error {
	return errors.Join(c.Module.Close(ctx), c.CompiledModule.Close(ctx))
}

func (c Converter) Convert(
	ctx context.Context,
	input []byte,
) (output []byte, e error) {
	return c.BytesToBytesRaw.Map(ctx, input)
}

func (c Converter) AsBytesToBytes() mb.BytesToBytes { return c }

type ConverterFactory struct {
	w0.Runtime
	w0.ModuleConfig

	SetInputSize   string
	GetOutEstimate string
	SetOutputSize  string
	Converter      string
	OffsetI        string
	OffsetO        string
}

func (f ConverterFactory) Close(ctx context.Context) error {
	return f.Runtime.Close(ctx)
}

func (f ConverterFactory) ToConverter(
	ctx context.Context,
	wasmBytes []byte,
) (Converter, error) {
	compiled, e := f.Runtime.CompileModule(ctx, wasmBytes)
	if nil != e {
		return Converter{}, e
	}

	instance, e := f.Runtime.InstantiateModule(
		ctx,
		compiled,
		f.ModuleConfig,
	)
	if nil != e {
		return Converter{}, errors.Join(e, compiled.Close(ctx))
	}

	b2br := BytesToBytesRaw{
		SetInputSize:   instance.ExportedFunction(f.SetInputSize),
		GetOutEstimate: instance.ExportedFunction(f.GetOutEstimate),
		SetOutputSize:  instance.ExportedFunction(f.SetOutputSize),
		Converter:      instance.ExportedFunction(f.Converter),
		OffsetI:        instance.ExportedFunction(f.OffsetI),
		OffsetO:        instance.ExportedFunction(f.OffsetO),
		Memory:         instance.Memory(),
	}

	var valid bool = b2br.Validate()

	converter := Converter{
		CompiledModule:  compiled,
		Module:          instance,
		BytesToBytesRaw: b2br,
	}

	switch valid {
	case true:
		return converter, nil
	default:
		return Converter{}, errors.Join(ErrInvalidModule, converter.Close(ctx))
	}
}

type ConverterMap map[bw.Name]Converter

func (m ConverterMap) Close(ctx context.Context) error {
	var cnvs iter.Seq[Converter] = maps.Values(m)
	var earr []error
	for cnv := range cnvs {
		earr = append(earr, cnv.Close(ctx))
	}
	return errors.Join(earr...)
}

func ConverterFactoryNewDefault(ctx context.Context) ConverterFactory {
	return ConverterFactory{
		Runtime:      w0.NewRuntime(ctx),
		ModuleConfig: w0.NewModuleConfig().WithName(""),

		SetInputSize:   bw.SetInputSize,
		GetOutEstimate: bw.GetOutEstimate,
		SetOutputSize:  bw.SetOutputSize,
		Converter:      bw.Converter,
		OffsetI:        bw.OffsetI,
		OffsetO:        bw.OffsetO,
	}
}
