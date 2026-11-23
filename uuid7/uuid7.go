package uuid7

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

func New() string {
	var b [16]byte

	// 48-bit timestamp milidetik (Unix epoch)
	ms := uint64(time.Now().UnixMilli())
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)

	// 10 byte acak sisanya
	_, _ = rand.Read(b[6:])

	// Set versi 7 (0111) dan varian RFC 4122 (10xx)
	b[6] = (b[6] & 0x0f) | 0x70
	b[8] = (b[8] & 0x3f) | 0x80

	_ = binary.BigEndian

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}