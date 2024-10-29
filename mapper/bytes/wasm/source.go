package wasm

import (
	"context"
	"errors"

	cm "github.com/takanoriyanagitani/go-cbor-map/mapper"
	mb "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes"
)

const (
	DefaultWasmExt string = "wasm"
)

var (
	ErrNotFound error = errors.New("no wasm found")
)

type Name string

type WasmSource func(context.Context, Name) ([]byte, error)

func NamesToConverters[C mb.BytesToBytes](
	ctx context.Context,
	ws WasmSource,
	names []Name,
	bytes2conv func(context.Context, []byte) (C, error),
	onError func(context.Context, map[Name]C) error,
) (ConverterMap[C], error) {
	m := map[Name]C{}
	for _, name := range names {
		wbytes, e := ws(ctx, name)
		if nil != e {
			return nil, errors.Join(e, onError(ctx, m))
		}

		conv, e := bytes2conv(ctx, wbytes)
		if nil != e {
			return nil, errors.Join(e, onError(ctx, m))
		}

		m[name] = conv
	}
	return m, nil
}

type ConverterMap[C mb.BytesToBytes] map[Name]C

func (m ConverterMap[C]) GetMapperByName(n Name, alt cm.Mapper) cm.Mapper {
	mapper, found := m[n]
	switch found {
	case true:
		return mb.BytesToBytesFn(mapper.Convert).ToMapper()
	default:
		return alt
	}
}

type IndexToName func(context.Context, uint32) (Name, error)

func (m ConverterMap[C]) GetMapperByIndex(
	ctx context.Context,
	index uint32,
	i2n IndexToName,
	alt cm.Mapper,
) cm.Mapper {
	name, e := i2n(ctx, index)
	if nil != e {
		return alt
	}
	return m.GetMapperByName(name, alt)
}

func (m ConverterMap[C]) ToMapperMap(
	ctx context.Context,
	indices []uint32,
	i2n IndexToName,
	alt cm.Mapper,
) cm.MapperMap {
	mm := map[uint32]cm.Mapper{}
	for _, idx := range indices {
		mapper := m.GetMapperByIndex(
			ctx,
			idx,
			i2n,
			alt,
		)
		mm[idx] = mapper
	}
	return mm
}

func (m ConverterMap[C]) ToMapperMapDefault(
	ctx context.Context,
	indices []uint32,
	i2n IndexToName,
) cm.MapperMap {
	return m.ToMapperMap(
		ctx,
		indices,
		i2n,
		cm.MapperIdentity,
	)
}
