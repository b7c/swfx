package swfx

import (
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io"

	"github.com/b7c/swfx/tagcode"
)

type SwfReader interface {
	io.Reader
	io.ByteReader
	WrapInnerStream(wrap func(inner io.Reader) io.Reader)
	Position() int
	// Reads into the specified byte slice.
	// Panics if the slice is unable to be filled.
	MustRead([]byte)
	// Reads a single bit and advances the bit position.
	ReadBit() byte
	// Reads a single bit as a boolean.
	ReadBool() bool
	// Reads n bits as an unsigned integer.
	ReadUbits(n int) int
	// Reads n bits as a signed integer.
	ReadSbits(n int) int
	// Reads an unsigned 8-bit integer.
	ReadUint8() uint8
	// Reads an unsigned 16-bit integer.
	ReadUint16() uint16
	// Reads an unsigned 32-bit integer.
	ReadUint32() uint32
	// Reads a string.
	ReadString(limit int) string
	// Reads a rectangle.
	ReadRect() Rect
	// Reads a tag code and length.
	ReadTagCodeAndLength() (code tagcode.TagCode, length uint32)
	// Reads a tag.
	ReadTag(code tagcode.TagCode, length uint32) SwfTag
}

type swfReader struct {
	r      io.Reader
	buf    [8]byte
	pos    int
	bitPos int
}

func (r *swfReader) WrapInnerStream(wrap func(inner io.Reader) io.Reader) {
	r.r = wrap(r.r)
}

func (r *swfReader) Position() int {
	return r.pos
}

func (r *swfReader) ReadByte() (byte, error) {
	r.align()
	if byteReader, ok := r.r.(io.ByteReader); ok {
		b, err := byteReader.ReadByte()
		if err != nil {
			r.pos++
		}
		return b, err
	}
	n, err := r.r.Read(r.buf[:1])
	r.pos += n
	return r.buf[0], err
}

func (r *swfReader) Read(p []byte) (int, error) {
	r.align()
	n, err := r.r.Read(p)
	r.pos += n
	return n, err
}

func (r *swfReader) MustRead(p []byte) {
	read := 0
	for read < len(p) {
		n, err := r.Read(p[read:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = io.ErrUnexpectedEOF
			}
			panic(err)
		}
		read += n
	}
}

func (r *swfReader) ReadBool() bool {
	return r.ReadBit() == 1
}

func (r *swfReader) ReadUbits(n int) int {
	var bit byte
	var value int
	for i := 0; i < n; i++ {
		bit = r.ReadBit()
		value = (value << 1) | int(bit)
	}
	return value
}

func (r *swfReader) ReadSbits(n int) int {
	// sign-extend
	shift := 32 - n
	return int(int32(r.ReadUbits(n)) << shift >> shift)
}

func (r *swfReader) ReadFbits(n int) int {
	panic("not implemented")
}

func (r *swfReader) ReadRect() Rect {
	r.align()
	nbits := r.ReadUbits(5)
	xmin := r.ReadSbits(nbits)
	xmax := r.ReadSbits(nbits)
	ymin := r.ReadSbits(nbits)
	ymax := r.ReadSbits(nbits)
	return Rect{
		Xmin: xmin,
		Xmax: xmax,
		Ymin: ymin,
		Ymax: ymax,
	}
}

func (r *swfReader) fill(count int) {
	n, err := r.Read(r.buf[:count])
	if n < count || (err != nil && !errors.Is(err, io.EOF)) {
		panic(err)
	}
}

func (r *swfReader) ReadUint8() uint8 {
	r.fill(1)
	return r.buf[0]
}

func (r *swfReader) ReadUint16() uint16 {
	r.fill(2)
	return binary.LittleEndian.Uint16(r.buf[:2])
}

func (r *swfReader) ReadUint32() uint32 {
	r.fill(4)
	return binary.LittleEndian.Uint32(r.buf[:4])
}

func (r *swfReader) ReadString(limit int) string {
	buf := make([]byte, 0, limit)
	read := 0
	for read < limit {
		b, err := r.ReadByte()
		if err != nil {
			panic(err)
		}
		if b == 0 {
			break
		}
		buf = append(buf, b)
		read++
	}
	return string(buf)
}

func (r *swfReader) ReadBit() byte {
	var bit byte
	var err error
	if r.bitPos == 0 {
		bit, err = r.ReadByte()
		if err != nil && !errors.Is(err, io.EOF) {
			panic(err)
		}
		r.buf[0] = bit
	} else {
		bit = r.buf[0]
	}
	r.bitPos++
	bit = (bit >> (8 - r.bitPos)) & 1
	r.bitPos %= 8
	return bit
}

func (r *swfReader) ReadTagCodeAndLength() (code tagcode.TagCode, length uint32) {
	tagCodeAndLength := r.ReadUint16()
	code = tagcode.TagCode(tagCodeAndLength >> 6)
	length = uint32(tagCodeAndLength & 0x3f)
	if length == 0x3f {
		length = r.ReadUint32()
	}
	return
}

func (r *swfReader) ReadTag(tagCode tagcode.TagCode, length uint32) SwfTag {
	t := MakeTag(tagCode)
	t.readData(r, int(length))
	return t
}

func (r *swfReader) align() {
	r.bitPos = 0
}

func NewReader(r io.Reader) SwfReader {
	return &swfReader{
		r: r,
	}
}

// Reads the SWF header and applies compression if necessary as specified by the file signature.
func ReadHeader(r SwfReader) (header SwfHeader, err error) {
	buf := make([]byte, 8)
	_, err = r.Read(buf)
	if err != nil {
		return
	}
	switch buf[0] {
	case 'F':
	case 'C':
		r.WrapInnerStream(func(inner io.Reader) io.Reader {
			inner, err = zlib.NewReader(inner)
			if err != nil {
				panic(err)
			}
			return inner
		})
	default:
		err = errors.ErrUnsupported
		return
	}

	header = SwfHeader{
		Version:    int(buf[3]),
		FrameSize:  r.ReadRect(),
		FrameRate:  r.ReadUint16(),
		FrameCount: r.ReadUint16(),
	}
	return
}
