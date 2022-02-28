package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	filePath = "test_log1"
	main()
	assert.Equal(t, 1, len(res))
}
