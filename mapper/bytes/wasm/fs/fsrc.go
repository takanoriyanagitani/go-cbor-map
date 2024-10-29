package fsrc

import (
	"context"
	"io/fs"
	"strconv"

	bw "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes/wasm"
)

type FsSource struct {
	fs.FS
	WasmExtStr string
}

func (f FsSource) ReadWasmBytesTrusted(basename string) ([]byte, error) {
	return fs.ReadFile(f.FS, basename) // TODO: use limited reader
}

func (f FsSource) ToWasmSource() bw.WasmSource {
	return func(_ context.Context, n bw.Name) ([]byte, error) {
		var nm string = string(n)
		var full string = nm + "." + f.WasmExtStr
		return f.ReadWasmBytesTrusted(full)
	}
}

func (f FsSource) IsWasmExists(basename string) bool {
	_, e := fs.Stat(f.FS, basename)
	return nil == e
}

type WasmExt string

func (w WasmExt) ToIndexToNameDefault() bw.IndexToName {
	return func(_ context.Context, idx uint32) (bw.Name, error) {
		return bw.Name(strconv.Itoa(int(idx))), nil
	}
}

var WasmExtDefault WasmExt = WasmExt("wasm")
