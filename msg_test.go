package dsomessage

import (
	"testing"

	"github.com/miekg/dns"
)

func TestMsgPackWithCompressionAndPadding(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name     string
		blockLen uint16 // 0 means m.Len()
	}{
		{"short buffer", 0},
		{"short block", 2},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := NewUniMsg(&Push{[]dns.RR{
				aRRf("a. IN A 192.0.2.1"),
				aRRf("a. IN A 192.0.2.2"),
			}})
			maxLen := m.Len()
			blockLen := tc.blockLen
			if blockLen == 0 {
				blockLen = uint16(maxLen)
			}
			msg, err := m.PackWith(PackOptions{
				BlockLen:    blockLen,
				Compression: true,
			})
			if err != nil {
				t.Fatalf("Got PackWith() = %v, want <nil>", err)
			}
			if d := maxLen - len(msg); d == 0 || d >= TLVHeaderLen {
				t.Fatalf("Got %v extra bytes, want 1-3", d)
			}

			got, err := UnpackMsg(msg, OriginServer)
			if err != nil {
				t.Fatalf("Got UnpackMsg() = %v, want <nil>", err)
			}
			if len(got.TLV) != 2 || !got.TLV[0].Equal(m.TLV[0]) || got.TLV[1].Type() != TypeEncryptionPadding {
				t.Errorf("Got TLVs = %v, want EncryptionPadding", got.TLV)
			}
		})
	}
}
