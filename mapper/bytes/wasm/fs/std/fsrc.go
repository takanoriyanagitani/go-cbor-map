package stdfs

import (
	"os"

	wf "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes/wasm/fs"
)

type FsStdSource struct{ Dirname string }

func (s FsStdSource) ToFsSource() wf.FsSource {
	return wf.FsSource{
		FS:         os.DirFS(s.Dirname),
		WasmExtStr: string(wf.WasmExtDefault),
	}
}
