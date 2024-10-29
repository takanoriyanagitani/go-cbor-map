package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	gm "github.com/takanoriyanagitani/go-cbor-map"
	cm "github.com/takanoriyanagitani/go-cbor-map/mapper"

	bw "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes/wasm"

	wf "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes/wasm/fs"
	sf "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes/wasm/fs/std"

	b0 "github.com/takanoriyanagitani/go-cbor-map/mapper/bytes/wasm/wazero"

	ci "github.com/takanoriyanagitani/go-cbor-map/iter/cbor2iter"
	ca "github.com/takanoriyanagitani/go-cbor-map/iter/cbor2iter/amacker"

	ic "github.com/takanoriyanagitani/go-cbor-map/iter/iter2cbor"
	ia "github.com/takanoriyanagitani/go-cbor-map/iter/iter2cbor/amacker"
)

func rdr2wtr(
	ctx context.Context,
	wasmModuleDir string,
	indices []uint32,
	idx2name func(uint32) string,
	rdr io.Reader,
	wtr io.Writer,
) error {
	fs2 := sf.FsStdSource{Dirname: wasmModuleDir}
	var fsrc wf.FsSource = fs2.ToFsSource()
	var wsrc bw.WasmSource = fsrc.ToWasmSource()

	var names []bw.Name
	for _, idx := range indices {
		names = append(names, bw.Name(idx2name(idx)))
	}

	var cf b0.ConverterFactory = b0.ConverterFactoryNewDefault(ctx)
	defer cf.Close(ctx)

	convs, e := bw.NamesToConverters(
		ctx,
		wsrc,
		names,
		cf.ToConverter,
		func(ctx context.Context, m map[bw.Name]b0.Converter) error {
			return b0.ConverterMap(m).Close(ctx)
		},
	)
	if nil != e {
		return e
	}

	var we wf.WasmExt = wf.WasmExtDefault
	var i2n bw.IndexToName = we.ToIndexToNameDefault()

	var mm cm.MapperMap = convs.ToMapperMapDefault(
		ctx,
		indices,
		i2n,
	)

	var a2m gm.ArrayToMapd = mm.ToArrayToMapd()

	var c2a ca.CborToArr = ca.CborToArrNew(rdr)
	var c2i ci.CborToIter = c2a.AsCborToIter()
	var src gm.ArraySource = c2i.ToArraySource()

	var a2c ia.ArrToCbor = ia.ArrToCborNew(wtr)
	var ac ic.ArrayToCbor = a2c.AsArrayToCbor()
	var ao gm.ArrayOutput = ac.ToArrayOutput()

	om := gm.OutputMapd{
		ArraySource: src,
		ArrayToMapd: a2m,
		ArrayOutput: ao,
	}

	return om.MapAll(ctx)
}

func stdin2stdout(
	ctx context.Context,
	wasmModuleDir string,
	indices string,
	idx2name func(uint32) string,
) error {
	var r io.Reader = os.Stdin
	var br io.Reader = bufio.NewReader(r)

	var w io.Writer = os.Stdout
	var bw *bufio.Writer = bufio.NewWriter(w)
	defer bw.Flush()

	var iarr []uint32
	var splited []string = strings.Split(indices, ",")
	for _, idx := range splited {
		parsed, e := strconv.Atoi(idx)
		if nil != e {
			if "" == indices {
				break
			}
			return e
		}
		iarr = append(iarr, uint32(parsed))
	}

	return rdr2wtr(
		ctx,
		wasmModuleDir,
		iarr,
		idx2name,
		br,
		bw,
	)
}

func sub(ctx context.Context) error {
	var wmdir string = os.Getenv("ENV_WASM_MODULE_DIR")
	var indices string = os.Getenv("ENV_MAP_INDICES")
	idx2name := func(idx uint32) string {
		var s string = strconv.Itoa(int(idx))
		return s
	}
	return stdin2stdout(
		ctx,
		wmdir,
		indices,
		idx2name,
	)
}

func main() {
	e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
