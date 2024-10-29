package mapper

import (
	"context"

	cm "github.com/takanoriyanagitani/go-cbor-map"
)

type Original any
type Mapd any

type Mapper func(context.Context, Original) (Mapd, error)

var MapperIdentity Mapper = func(_ context.Context, o Original) (Mapd, error) {
	return Mapd(o), nil
}

type MapperMap map[uint32]Mapper

func (m MapperMap) GetMapper(idx uint32) Mapper {
	mapper, found := m[idx]
	switch found {
	case true:
		return mapper
	default:
		return MapperIdentity
	}
}

func (m MapperMap) MapArray(ctx context.Context, arr []any) ([]any, error) {
	var mapd []any = make([]any, 0, len(arr))
	for i, input := range arr {
		var mapper Mapper = m.GetMapper(uint32(i))
		output, e := mapper(ctx, input)
		if nil != e {
			return nil, e
		}
		mapd = append(mapd, output)
	}
	return mapd, nil
}

func (m MapperMap) ToArrayToMapd() cm.ArrayToMapd {
	return func(ctx context.Context, o cm.Original) (cm.Mapd, error) {
		return m.MapArray(ctx, o)
	}
}
