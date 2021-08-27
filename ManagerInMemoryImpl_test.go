package acccore

import (
	"context"
	"math/big"
	"testing"
)

type ExchangeTest struct {
	from       string
	fromAmount int64
	to         string
	toAmount   int64
}

func TestInMemoryExchangeManager_CalculateExchange(t *testing.T) {
	exchangeManager := NewInMemoryExchangeManager(map[string]*big.Float{
		"PLATINUM": big.NewFloat(0.001),
		"GOLD":     big.NewFloat(0.01),
		"SILVER":   big.NewFloat(0.1),
		"COPPER":   big.NewFloat(1.0),
	})
	testData := []*ExchangeTest{
		{
			from:       "GOLD",
			fromAmount: 1000,
			to:         "PLATINUM",
			toAmount:   100,
		},
		{
			from:       "GOLD",
			fromAmount: 1000,
			to:         "SILVER",
			toAmount:   10000,
		},
		{
			from:       "GOLD",
			fromAmount: 1000,
			to:         "GOLD",
			toAmount:   1000,
		},
	}
	ctx := context.Background()
	for idx, data := range testData {
		result, err := exchangeManager.CalculateExchange(ctx, data.from, data.to, data.fromAmount)
		if err != nil {
			t.Errorf(err.Error())
			t.Fail()
		}
		if result != data.toAmount {
			t.Errorf("#%d. Expect %d but %d", idx, data.toAmount, result)
			t.Fail()
		}
	}

	exchangeManager.SetDenom(ctx, big.NewFloat(123.456))
	for idx, data := range testData {
		result, err := exchangeManager.CalculateExchange(ctx, data.from, data.to, data.fromAmount)
		if err != nil {
			t.Errorf(err.Error())
			t.Fail()
		}
		if result != data.toAmount {
			t.Errorf("#%d. Expect %d but %d", idx, data.toAmount, result)
			t.Fail()
		}
	}
}
