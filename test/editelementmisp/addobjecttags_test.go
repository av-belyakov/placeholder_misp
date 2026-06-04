package editelementmisp_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/av-belyakov/objectsmispformat"
	"github.com/av-belyakov/placeholder_misp/internal/mispapi"
	"github.com/stretchr/testify/assert"
	"github.com/subosito/gotenv"
)

// const Event_Id = "44066"
const Event_Id = "142706"

func TestAddObjectTags(t *testing.T) {
	if err := gotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"Host:'%s'\nUser Auth:'%s'\nMaster Auth:'%s'\n",
		os.Getenv("GO_PHMISP_MHOST"),
		os.Getenv("GO_PHMISP_MAUTH"),
		os.Getenv("GO_PHMISP_MAUTH"),
	)

	requestMisp, err := mispapi.NewMispRequest(
		mispapi.WithHost(os.Getenv("GO_PHMISP_MHOST")),
		mispapi.WithUserAuthKey(os.Getenv("GO_PHMISP_MAUTH")),
		mispapi.WithMasterAuthKey(os.Getenv("GO_PHMISP_MAUTH")),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// отменяем публикацию события
	//_, err = requestMisp.SendRequestUnpublishEvent_ForTest(t.Context(), Event_Id)
	//assert.NoError(t, err)

	err = requestMisp.AddTagToEvent_ForTest(
		t.Context(),
		objectsmispformat.EventObjectTagsMispFormat{
			Event: Event_Id,
			//Tag:   "misp-galaxy:Sector=\"Оборонная промышленность\"",
			//Tag: "ATs:geoip=\"Россия\"",
			Tag: "misp-galaxy:Sector=\"Государственные учреждения\"",
		},
	)
	assert.NoError(t, err)

	// публикуем событие повторно
	_, err = requestMisp.SendRequestPublishEvent_ForTest(t.Context(), Event_Id)
	assert.NoError(t, err)
}
