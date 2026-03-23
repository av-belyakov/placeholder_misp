package testaddnewelements_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"

	"github.com/av-belyakov/placeholder_misp/cmd/mispapi"
)

func TestHandlerTags(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	os.Unsetenv("GO_PHMISP_MHOST")
	os.Unsetenv("GO_PHMISP_MAUTH")

	if err := gotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	requestMisp, err := mispapi.NewMispRequest(
		mispapi.WithHost(os.Getenv("GO_PHMISP_MHOST")),
		mispapi.WithUserAuthKey(os.Getenv("GO_PHMISP_MAUTH")),
		mispapi.WithMasterAuthKey(os.Getenv("GO_PHMISP_MAUTH")),
	)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Host:", os.Getenv("GO_PHMISP_MHOST"))
	fmt.Println("Pass:", os.Getenv("GO_PHMISP_MAUTH"))

	anyTag := "Sensor:ID=\"1200000\""

	t.Run("Test 1. Add any tag", func(t *testing.T) {
		res, err := requestMisp.AddTag_ForTest(ctx, anyTag, "#FF0000")
		assert.NoError(t, err)

		fmt.Printf("Add tag response:'%+v'\n", res)
	})

	t.Run("Test 2. Search tag is success", func(t *testing.T) {
		res, err := requestMisp.SearchTag_ForTest(ctx, anyTag)
		assert.NoError(t, err)

		fmt.Printf("Search tag response:'%+v'\n", res)
	})

	t.Run("Test 3. Search tag is failure", func(t *testing.T) {
		res, err := requestMisp.SearchTag_ForTest(ctx, "Sensor:ID=\"00000\"")
		assert.NoError(t, err)

		fmt.Printf("Search tag response:'%+v'\n", res)
	})

	t.Cleanup(func() {
		cancel()

		os.Unsetenv("GO_PHMISP_MHOST")
		os.Unsetenv("GO_PHMISP_MAUTH")
	})
}
