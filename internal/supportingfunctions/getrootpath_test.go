package supportingfunctions_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_misp/v2/internal/supportingfunctions"
)

func TestGetRootPath(t *testing.T) {
	rootDir := "placeholder_misp"

	str, err := supportingfunctions.GetRootPath(rootDir)
	assert.NoError(t, err)

	fmt.Println("ROOT DIR:", str)

	assert.Equal(t, str, "/home/artemij/go/src/placeholder_misp")
}
