package inspector

import (
	"bytes"
	"testing"

	testv1 "github.com/detailyang/pb-inspector-go/proto/go/proto/test/v1"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestInspectWithoutSchema(t *testing.T) {
	test := testv1.Test{
		Int32:   1,
		Int64:   2,
		Float:   3.0,
		Double:  4.0,
		Uint32:  5,
		Uint64:  6,
		Bool:    true,
		Bytes:   []byte("haha"),
		String_: "hello world",
	}

	raw, err := proto.Marshal(&test)
	require.Nil(t, err)

	w := bytes.NewBuffer(nil)
	in := NewInspector()
	err = in.InspectWithoutSchema(false, raw, w)
	require.Nil(t, err)

	o := ""
	o += `  0: t=  1 varint 1` + "\n"
	o += "  2: t=  2 varint 2" + "\n"
	o += "  4: t=  3 fix32 1077936128" + "\n"
	o += "  9: t=  4 fix64 4616189618054758400" + "\n"
	o += " 18: t=  5 varint 5" + "\n"
	o += " 20: t=  6 varint 6" + "\n"
	o += " 22: t=  7 varint 1" + "\n"
	o += " 24: t=  8 bytes [4] 68 61 68 61" + "\n"
	o += " 30: t=  9 bytes [11] 68 65 6c 72 6c 64" + "\n"

	require.Equal(t, o, w.String())
}
