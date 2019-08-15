package inspector

import (
	"errors"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
)

// errOverflow is returned when an integer is too large to be represented.
var errOverflow = errors.New("Buffer: integer overflow")

// ErrInternalBadWireType is returned by generated code when an incorrect
// wire type is encountered. It does not get returned to user code.
var ErrInternalBadWireType = errors.New("Buffer: bad wiretype for oneof")

// A Buffer is a buffer manager for marshaling and unmarshaling
// protocol buffers.  It may be reused between invocations to
// reduce memory usage.  It is not necessary to use a Buffer;
// the global functions Marshal and Unmarshal create a
// temporary Buffer and are fine for most applications.
type Buffer struct {
	buf     []byte // encode/decode byte stream
	index   int    // read point
	decoder *decoder
}

// NewBuffer allocates a new Buffer and initializes its internal data to
// the contents of the argument slice.
func NewBuffer(e []byte) *Buffer {
	return &Buffer{
		buf: e,
	}
}

// Reset resets the Buffer, ready for marshaling a new protocol buffer.
func (p *Buffer) Reset() {
	p.buf = p.buf[0:0] // for reading/writing
	p.index = 0        // for reading
}

// SetBuf replaces the internal buffer with the slice,
// ready for unmarshaling the contents of the slice.
func (p *Buffer) SetBuf(s []byte) {
	p.buf = s
	p.index = 0
}

// Bytes returns the contents of the Buffer.
func (p *Buffer) Bytes() []byte { return p.buf }

func (p *Buffer) decodeVarintSlow() (x uint64, err error) {
	i := p.index
	l := len(p.buf)

	for shift := uint(0); shift < 64; shift += 7 {
		if i >= l {
			err = io.ErrUnexpectedEOF
			return
		}
		b := p.buf[i]
		i++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			p.index = i
			return
		}
	}

	// The number is too large to represent in a 64-bit value.
	err = errOverflow
	return
}

// DecodeVarint reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
func (p *Buffer) DecodeVarint() (x uint64, err error) {
	i := p.index
	buf := p.buf

	if i >= len(buf) {
		return 0, io.ErrUnexpectedEOF
	} else if buf[i] < 0x80 {
		p.index++
		return uint64(buf[i]), nil
	} else if len(buf)-i < 10 {
		return p.decodeVarintSlow()
	}

	var b uint64
	// we already checked the first byte
	x = uint64(buf[i]) - 0x80
	i++

	b = uint64(buf[i])
	i++
	x += b << 7
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 7

	b = uint64(buf[i])
	i++
	x += b << 14
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 14

	b = uint64(buf[i])
	i++
	x += b << 21
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 21

	b = uint64(buf[i])
	i++
	x += b << 28
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 28

	b = uint64(buf[i])
	i++
	x += b << 35
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 35

	b = uint64(buf[i])
	i++
	x += b << 42
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 42

	b = uint64(buf[i])
	i++
	x += b << 49
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 49

	b = uint64(buf[i])
	i++
	x += b << 56
	if b&0x80 == 0 {
		goto done
	}
	x -= 0x80 << 56

	b = uint64(buf[i])
	i++
	x += b << 63
	if b&0x80 == 0 {
		goto done
	}

	return 0, errOverflow

done:
	p.index = i
	return x, nil
}

// DecodeFixed64 reads a 64-bit integer from the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (p *Buffer) DecodeFixed64() (x uint64, err error) {
	// x, err already 0
	i := p.index + 8
	if i < 0 || i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-8])
	x |= uint64(p.buf[i-7]) << 8
	x |= uint64(p.buf[i-6]) << 16
	x |= uint64(p.buf[i-5]) << 24
	x |= uint64(p.buf[i-4]) << 32
	x |= uint64(p.buf[i-3]) << 40
	x |= uint64(p.buf[i-2]) << 48
	x |= uint64(p.buf[i-1]) << 56
	return
}

// DecodeFixed32 reads a 32-bit integer from the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (p *Buffer) DecodeFixed32() (x uint64, err error) {
	// x, err already 0
	i := p.index + 4
	if i < 0 || i > len(p.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	p.index = i

	x = uint64(p.buf[i-4])
	x |= uint64(p.buf[i-3]) << 8
	x |= uint64(p.buf[i-2]) << 16
	x |= uint64(p.buf[i-1]) << 24
	return
}

// DecodeZigzag64 reads a zigzag-encoded 64-bit integer
// from the Buffer.
// This is the format used for the sint64 protocol buffer type.
func (p *Buffer) DecodeZigzag64() (x uint64, err error) {
	x, err = p.DecodeVarint()
	if err != nil {
		return
	}
	x = (x >> 1) ^ uint64((int64(x&1)<<63)>>63)
	return
}

// DecodeZigzag32 reads a zigzag-encoded 32-bit integer
// from  the Buffer.
// This is the format used for the sint32 protocol buffer type.
func (p *Buffer) DecodeZigzag32() (x uint64, err error) {
	x, err = p.DecodeVarint()
	if err != nil {
		return
	}
	x = uint64((uint32(x) >> 1) ^ uint32((int32(x&1)<<31)>>31))
	return
}

// DecodeRawBytes reads a count-delimited byte buffer from the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (p *Buffer) DecodeRawBytes(alloc bool) (buf []byte, err error) {
	n, err := p.DecodeVarint()
	if err != nil {
		return nil, err
	}

	nb := int(n)
	if nb < 0 {
		return nil, fmt.Errorf("proto: bad byte length %d", nb)
	}

	end := p.index + nb
	if end < p.index || end > len(p.buf) {
		return nil, io.ErrUnexpectedEOF
	}

	if !alloc {
		// todo: check if can get more uses of alloc=false
		buf = p.buf[p.index:end]
		p.index += nb
		return
	}

	buf = make([]byte, nb)
	p.index += copy(buf, p.buf[p.index:])
	return
}

// DecodeStringBytes reads an encoded string from the Buffer.
// This is the format used for the proto2 string type.
func (p *Buffer) DecodeStringBytes() (s string, err error) {
	buf, err := p.DecodeRawBytes(false)
	if err != nil {
		return
	}
	return string(buf), nil
}

// InspectWithoutSchema inspects raw protobuf binary data and write to w
func (p *Buffer) InspectWithoutSchema(verbose bool, raw []byte, w io.Writer) error {
	var (
		err    error
		u      = uint64(0)
		obuf   = p.buf
		sindex = p.index
		depth  = 0
		op     uint64
	)

	p.buf = raw
	p.index = 0

out:
	for {
		for i := 0; i < depth; i++ {
			fmt.Fprint(w, "  ")
		}

		index := p.index
		if index == len(p.buf) {
			break
		}

		// Each key in the streamed message is a varint with the value (field_number << 3) | wire_type –
		// in other words, the last three bits of the number store the wire type.
		op, err = p.DecodeVarint()
		if err != nil {
			err = fmt.Errorf("insepctor: [%3d] fetching op err %v", index, err)
			break out
		}
		tag := op >> 3
		wire := op & 7

		switch wire {
		default:
			err = fmt.Errorf("Buffer: [%3d] t=%3d unknown wire=%d", index, tag, wire)
			break out

		case proto.WireBytes:
			var r []byte

			r, err = p.DecodeRawBytes(true)
			if err != nil {
				err = fmt.Errorf("insepctor: [%3d] t=%3d bytes err %v", index, tag, err)
				break out
			}

			fmt.Fprintf(w, "%3d: t=%3d bytes [%d]", index, tag, len(r))

			if verbose {
				for i := 0; i < len(r); i++ {
					fmt.Fprintf(w, " %.2x", r[i])
				}

			} else {
				if len(r) <= 6 {
					for i := 0; i < len(r); i++ {
						fmt.Fprintf(w, " %.2x", r[i])
					}
				} else {
					for i := 0; i < 3; i++ {
						fmt.Fprintf(w, " %.2x", r[i])
					}
					fmt.Printf(" ..")
					for i := len(r) - 3; i < len(r); i++ {
						fmt.Fprintf(w, " %.2x", r[i])
					}
				}
			}
			fmt.Fprintf(w, "\n")

		case proto.WireFixed32:
			u, err = p.DecodeFixed32()
			if err != nil {
				err = fmt.Errorf("insepctor: [%3d] t=%3d fix32 err %v", index, tag, err)
				break out
			}
			fmt.Fprintf(w, "%3d: t=%3d fix32 %d\n", index, tag, u)

		case proto.WireFixed64:
			u, err = p.DecodeFixed64()
			if err != nil {
				err = fmt.Errorf("insepctor: [%3d] t=%3d fix64 err %v", index, tag, err)
				break out
			}
			fmt.Fprintf(w, "%3d: t=%3d fix64 %d\n", index, tag, u)

		case proto.WireVarint:
			u, err = p.DecodeVarint()
			if err != nil {
				err = fmt.Errorf("%3d: t=%3d varint err %v", index, tag, err)
				break out
			}
			fmt.Fprintf(w, "%3d: t=%3d varint %d\n", index, tag, u)

		case proto.WireStartGroup:
			fmt.Fprintf(w, "%3d: t=%3d start\n", index, tag)
			depth++

		case proto.WireEndGroup:
			depth--
			fmt.Fprintf(w, "%3d: t=%3d end\n", index, tag)
		}
	}

	if depth != 0 {
		fmt.Fprintf(w, "%3d: start-end not balanced %d\n", p.index, depth)
	}

	p.buf = obuf
	p.index = sindex

	return err
}
