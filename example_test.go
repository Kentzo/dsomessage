package dsomessage

import (
	"errors"
	"fmt"

	"github.com/miekg/dns"
)

func ExampleBuilder() {
	b := NewBuilder(make([]byte, 128)).
		SetHeader(MsgHeader{42, false, dns.RcodeSuccess})
	b.WriteKeepAlive(&KeepAlive{InactivityTimeoutDefault, KeepAliveIntervalDefault})
	msg, err := b.Message()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("% x\n", msg)
	// Output: 00 2a 30 00 00 00 00 00 00 00 00 00 00 01 00 08 00 00 3a 98 00 00 3a 98
}

func ExampleParser() {
	var (
		msg = []byte("\x00\x2a\x30\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x01\x00\x08\x00\x00\x3a\x98\x00\x00\x3a\x98")

		p Parser

		header    MsgHeader
		tlvHeader TLVHeader
		usage     Usage
		tlv       TLV
		err       error
	)
	header, err = p.Start(msg, OriginClient)
	fmt.Printf("%+v\n", header)
	for err == nil {
		tlvHeader, err = p.TLVHeader()
		if err != nil {
			break
		}
		fmt.Printf("%+v\n", tlvHeader)

		usage = p.TLVUsage()
		fmt.Printf("Usage:%d\n", usage)
		tlv, err = p.TLV()
		if err != nil {
			break
		}
		fmt.Printf("%v\n", tlv)

		err = tlv.Verify(usage)
		if err != nil {
			break
		}
	}
	if err != nil && !errors.Is(err, ErrDone) {
		panic(err.Error())
	}
	// Output:
	// {ID:42 Response:false Rcode:0}
	// {Type:KeepAlive Length:8}
	// Usage:32
	// timeout 15s, interval 15s
}
