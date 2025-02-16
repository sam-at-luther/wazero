package wasi

import (
	"testing"

	"github.com/tetratelabs/wazero/internal/testing/require"
	"github.com/tetratelabs/wazero/internal/wasm"
)

var testMem = &wasm.MemoryInstance{
	Min: 1,
	Buffer: []byte{
		0,                // environBuf is after this
		'a', '=', 'b', 0, // null terminated "a=b",
		'b', '=', 'c', 'd', 0, // null terminated "b=cd"
		0,          // environ is after this
		1, 0, 0, 0, // little endian-encoded offset of "a=b"
		5, 0, 0, 0, // little endian-encoded offset of "b=cd"
		0,
	},
}

func Test_EnvironGet(t *testing.T) {
	sys, err := newSysContext(nil, []string{"a=b", "b=cd"}, nil)
	require.NoError(t, err)

	m := newModule(make([]byte, 20), sys)
	environGet := newSnapshotPreview1(testCtx).EnvironGet

	require.Equal(t, ErrnoSuccess, environGet(testCtx, m, 11, 1))
	require.Equal(t, m.Memory(), testMem)
}

func Benchmark_EnvironGet(b *testing.B) {
	sys, err := newSysContext(nil, []string{"a=b", "b=cd"}, nil)
	if err != nil {
		b.Fatal(err)
	}

	m := newModule([]byte{
		0,                // environBuf is after this
		'a', '=', 'b', 0, // null terminated "a=b",
		'b', '=', 'c', 'd', 0, // null terminated "b=cd"
		0,          // environ is after this
		1, 0, 0, 0, // little endian-encoded offset of "a=b"
		5, 0, 0, 0, // little endian-encoded offset of "b=cd"
		0,
	}, sys)

	environGet := newSnapshotPreview1(testCtx).EnvironGet
	b.Run("EnvironGet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if environGet(testCtx, m, 0, 4) != ErrnoSuccess {
				b.Fatal()
			}
		}
	})
}

func newModule(buf []byte, sys *wasm.SysContext) *wasm.CallContext {
	return wasm.NewCallContext(nil, &wasm.ModuleInstance{
		Memory: &wasm.MemoryInstance{Min: 1, Buffer: buf},
	}, sys)
}
