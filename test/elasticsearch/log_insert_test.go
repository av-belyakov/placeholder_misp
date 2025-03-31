package testelasticsearch_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/placeholder_misp/cmd/elasticsearchapi"
)

func TestInserMsgLog(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	esc, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:               9200,
		Host:               "datahook.cloud.gcm",
		User:               "log_writer",
		Passwd:             os.Getenv("GO_PHMISP_DBWLOGPASSWD"),
		IndexDB:            "placeholder_misp",
		NameRegionalObject: "test-region",
	})
	assert.NoError(t, err)

	testStr := `message from MISP: status '403 Forbidden', error - {
    'saved': false,
    'name': 'Could not delete Event',
    'message': 'Could not delete Event',
    'url': '\events\delete',
    'errors': 'Event was not deleted.'
}`

	err = esc.Write("error", testStr)
	assert.NoError(t, err)
}
