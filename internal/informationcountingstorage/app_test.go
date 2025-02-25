package informationcountingstorage_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	infocountstor "github.com/av-belyakov/placeholder_misp/internal/informationcountingstorage"
)

var (
	imcs   *infocountstor.InformationCountingStorage
	myTime time.Time = time.Date(2025, time.June, 21, 12, 01, 0, 0, time.FixedZone("Moscow", 0))
)

func TestMain(m *testing.M) {
	imcs = infocountstor.NewInformationMessageCountingStorage()

	os.Exit(m.Run())
}

func TestInformationCountingStorage(t *testing.T) {
	imcs.SetStartTime(myTime)
	assert.Equal(t, imcs.GetStartTime(), myTime)

	tmo := "test message one"
	imcs.Increase(tmo)
	imcs.Increase(tmo)
	imcs.Increase(tmo)
	count, err := imcs.GetCount(tmo)
	assert.NoError(t, err)
	assert.Equal(t, count, uint(3))

	tmt := "test message two"
	imcs.Increase(tmt)
	count, err = imcs.GetCount(tmt)
	assert.NoError(t, err)
	assert.Equal(t, count, uint(1))

	assert.Equal(t, len(imcs.GetAllCount()), 2)

	tsm := "some message"
	count, err = imcs.GetCount(tsm)
	assert.Error(t, err)
	assert.Equal(t, count, uint(0))

	imcs.SetCount(tsm, uint(13))
	count, err = imcs.GetCount(tsm)
	assert.NoError(t, err)
	assert.Equal(t, count, uint(13))
}
