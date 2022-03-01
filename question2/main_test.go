package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Expect struct {
	ServerIP     string
	FailedSpan   float64
	ResultLength int
}

func Test_Normal(t *testing.T) {

	exp := Expect{
		ServerIP:     "10.20.31.1/16",
		FailedSpan:   0.034722222222222224,
		ResultLength: 1,
	}

	filepath = "testlog1.txt"
	main()
	for _, r := range res {
		assert.Equal(t, exp.ResultLength, len(res))
		assert.Equal(t, exp.ServerIP, r.FailedHost)
		assert.Equal(t, exp.FailedSpan, r.FailedSpan.Hours())
	}
}
