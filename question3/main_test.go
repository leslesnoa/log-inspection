package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Expect struct {
	ServerIP     string
	OverloadSpan float64
	ResultLength int
}

func Test_Normal(t *testing.T) {

	exp := Expect{
		ServerIP:     "192.168.1.1/24",
		OverloadSpan: 0.002777777777777778,
		ResultLength: 1,
	}

	filepath = "testlog1.txt"
	main()
	for _, o := range overLoadServers {
		assert.Equal(t, exp.ResultLength, len(overLoadServers))
		assert.Equal(t, exp.ServerIP, o.ServerIP)
		assert.Equal(t, exp.OverloadSpan, o.Span.Hours())
	}
}
