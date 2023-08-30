package acccore

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBaseAccount_MarshalJSON(t *testing.T) {
	sample := &BaseAccount{
		Currency:      "CUR",
		AccountNumber: "1234",
		Name:          "Sample",
		Description:   "Sample Description",
		Alignment:     DEBIT,
		Balance:       decimal.NewFromInt(1234),
		COA:           "COA",
		CreateTime:    time.Date(2000, time.January, 1, 1, 1, 1, 1, time.UTC),
		CreateBy:      "Creator",
		UpdateTime:    time.Date(2000, time.January, 1, 1, 1, 1, 1, time.UTC),
		UpdateBy:      "Updater",
	}

	bytes, err := json.Marshal(sample)
	assert.NoError(t, err)
	t.Log(string(bytes))

	result := &BaseAccount{}
	err = json.Unmarshal(bytes, &result)
	assert.NoError(t, err)

	assert.Equal(t, sample.Name, result.Name)
}
