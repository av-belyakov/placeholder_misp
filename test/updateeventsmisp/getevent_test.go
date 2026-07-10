package updateeventsmisp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_misp/v2/internal/mispapi"
)

func TestGetEventElementMISP(t *testing.T) {
	var (
		eventId string = "43940" // = case id 39100
	)

	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Host:'%s', auth token:'%s'", os.Getenv("GO_PHMISP_MHOST"), os.Getenv("GO_PHMISP_MAUTH"))

	client, err := mispapi.NewClientMISP(os.Getenv("GO_PHMISP_MHOST"), os.Getenv("GO_PHMISP_MAUTH"), false)
	if err != nil {
		t.Fatal(err)
	}

	res, raw, err := client.Get(t.Context(), fmt.Sprintf("/events/view/%s", eventId), []byte{})
	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	//fmt.Println("Get event response:", string(raw))

	oldEvents := struct {
		Event struct {
			UUID string `json:"uuid"`
		} `json:"Event"`
	}{}

	assert.NoError(t, json.Unmarshal(raw, &oldEvents))

	fmt.Println("Events UUID:", oldEvents.Event.UUID)
}
