package refid

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
)

var (
	// crockford base32
	// ref: https://en.wikipedia.org/wiki/Base32#Crockford's_Base32
	Alphabet         = "0123456789abcdefghjkmnpqrstvwxyz"
	WordSafeEncoding = base32.NewEncoding(Alphabet).WithPadding(base32.NoPadding)
	// Nil is the nil RefId, that has all 128 bits set to zero.
	Nil = RefId{}
)

const size = 18
const tagIndex = 0
const uuidOffset = 2

type RefId struct {
	UUID uuid.UUID
	tag  byte `db:"-"`
}

func New() (RefId, error) {
	var refId RefId
	id, err := uuid.NewV7()
	if err != nil {
		return refId, err
	}
	refId.UUID = id
	return refId, nil
}

func MustNew() RefId {
	refId, err := New()
	if err != nil {
		panic(err)
	}
	return refId
}

func NewTagged(tag byte) (RefId, error) {
	refId, err := New()
	if err != nil {
		return refId, err
	}
	refId.SetTag(tag)
	return refId, nil
}

func MustNewTagged(tag byte) RefId {
	refId := MustNew()
	refId.SetTag(tag)
	return refId
}

func Parse(s string) (RefId, error) {
	var refId RefId
	err := refId.UnmarshalText([]byte(s))
	if err != nil {
		return refId, err
	}
	return refId, err
}

func MustParse(s string) RefId {
	refId, err := Parse(s)
	if err != nil {
		panic(`RefId: Parse(` + s + `): ` + err.Error())
	}
	return refId
}

func ParseTagged(tag byte, s string) (RefId, error) {
	refId, err := Parse(s)
	if err != nil {
		return refId, err
	}

	if !refId.HasTag(tag) {
		return refId, fmt.Errorf("RefId tag mismatch: %d != %d", refId.tag, tag)
	}
	return refId, nil
}

func MustParseTagged(tag byte, s string) RefId {
	refId, err := ParseTagged(tag, s)
	if err != nil {
		panic(`RefId: ExpectParse(` + s + `): ` + "RefId tag mismatch")
	}
	return refId
}

func FromBytes(input []byte) (RefId, error) {
	var refId RefId
	err := refId.UnmarshalBinary(input)
	if err != nil {
		return refId, err
	}
	return refId, nil
}

func FromString(input string) (RefId, error) {
	var refId RefId
	err := refId.UnmarshalText([]byte(input))
	if err != nil {
		return refId, err
	}
	return refId, nil
}

func FromBase64String(input string) (RefId, error) {
	var refId RefId
	bx, err := base64.RawURLEncoding.DecodeString(input)
	if err != nil {
		return refId, err
	}
	if len(bx) != size {
		return refId, fmt.Errorf("wrong unmarshal size")
	}
	refId.UUID, err = uuid.FromBytes(bx[uuidOffset:])
	if err != nil {
		return refId, err
	}
	refId.tag = bx[tagIndex]
	return refId, nil
}

func FromHexString(input string) (RefId, error) {
	var refId RefId
	bx, err := hex.DecodeString(input)
	if err != nil {
		return refId, err
	}
	if len(bx) != size {
		return refId, fmt.Errorf("wrong unmarshal size")
	}
	refId.UUID, err = uuid.FromBytes(bx[uuidOffset:])
	if err != nil {
		return refId, err
	}
	refId.tag = bx[tagIndex]
	return refId, nil
}

func (refId *RefId) SetTag(tag byte) *RefId {
	refId.tag = tag
	return refId
}

func (refId *RefId) ClearTag() *RefId {
	refId.tag = 0
	return refId
}

func (refId RefId) IsTagged() bool {
	return refId.tag != 0
}

func (refId RefId) HasTag(tag byte) bool {
	return (refId.IsTagged() && refId.tag == tag)
}

func (refId RefId) IsNil() bool {
	return refId == Nil
}

func (refId RefId) Equal(other RefId) bool {
	return refId.String() == other.String()
}

func (refId RefId) MarshalText() ([]byte, error) {
	return []byte(refId.String()), nil
}

func (refId RefId) Time() time.Time {
	u := refId.UUID[:]

	t := 0 |
		(int64(u[0]) << 40) |
		(int64(u[1]) << 32) |
		(int64(u[2]) << 24) |
		(int64(u[3]) << 16) |
		(int64(u[4]) << 8) |
		int64(u[5])
	return time.UnixMilli(t)
}

func (refId *RefId) UnmarshalText(b []byte) error {
	// lowercase, then replace ambigious chars
	b = bytes.ToLower(b)
	for i := range b {
		switch b[i] {
		case 'i', 'l':
			b[i] = '1'
		case 'o', 'O':
			b[i] = '0'
		}
	}
	bx := make([]byte, size)
	n, err := WordSafeEncoding.Decode(bx, b)
	if err != nil {
		return err
	}
	if n != size {
		return fmt.Errorf("wrong unmarshal size")
	}
	refId.UUID, err = uuid.FromBytes(bx[uuidOffset:])
	if err != nil {
		return err
	}
	refId.tag = bx[tagIndex]
	return nil
}

func (refId RefId) Bytes() []byte {
	b := make([]byte, size)
	b[tagIndex] = refId.tag
	copy(b[uuidOffset:], refId.UUID[:])
	return b
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (refId RefId) MarshalBinary() ([]byte, error) {
	return refId.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// It will return an error if the slice isn't 16 bytes long.
func (refId *RefId) UnmarshalBinary(data []byte) error {
	if len(data) != size {
		return fmt.Errorf("refid: RefId must be exactly 16 bytes long, got %d bytes", len(data))
	}
	copy(refId.UUID[:], data[uuidOffset:])
	refId.tag = data[tagIndex]
	return nil
}

func (refId RefId) String() string {
	data := make([]byte, size)
	data[0] = refId.tag
	copy(data[uuidOffset:], refId.UUID[:])
	return WordSafeEncoding.EncodeToString(data)
}

func (refId RefId) ToBase64String() string {
	data := make([]byte, size)
	data[0] = refId.tag
	copy(data[uuidOffset:], refId.UUID[:])
	return base64.RawURLEncoding.EncodeToString(data)
}

func (refId RefId) ToHexString() string {
	data := make([]byte, size)
	data[0] = refId.tag
	copy(data[uuidOffset:], refId.UUID[:])
	return hex.EncodeToString(data)
}

func (refId RefId) Format(f fmt.State, c rune) {
	if c == 'v' && f.Flag('#') {
		fmt.Fprintf(f, "%#v", refId.Bytes())
		return
	}
	switch c {
	case 'x', 'X':
		b := make([]byte, size*2)
		hex.Encode(b, refId.Bytes())
		if c == 'X' {
			bytes.ToUpper(b)
		}
		_, _ = f.Write(b)
	case 'v', 's', 'S':
		b, _ := refId.MarshalText()
		if c == 'S' {
			bytes.ToUpper(b)
		}
		_, _ = f.Write(b)
	case 'q':
		_, _ = f.Write([]byte{'"'})
		_, _ = f.Write(refId.Bytes())
		_, _ = f.Write([]byte{'"'})
	default:
		// invalid/unsupported format verb
		fmt.Fprintf(f, "%%!%c(refid.RefId=%s)", c, refId.String())
	}
}
