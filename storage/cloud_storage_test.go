package cloud_storage

import (
	"testing"

	"go.uber.org/mock/gomock"
)

type PriceData struct {
	Symbol   string
	Datetime string
	Exchange string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
}

func TestCloudStorage_LoadJSON(t *testing.T) {
	ctr := gomock.NewController(t)
	defer ctr.Finish()
}
