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
	ctx := context.Background()
	exchangeManager := NewInMemoryExchangeManager()
	_, _ = exchangeManager.CreateCurrency(ctx, "PLATINUM", "Platinum", big.NewFloat(0.001), "superman")
	_, _ = exchangeManager.CreateCurrency(ctx, "GOLD", "Gold", big.NewFloat(0.01), "superman")
	_, _ = exchangeManager.CreateCurrency(ctx, "SILVER", "Silver", big.NewFloat(0.1), "superman")
	_, _ = exchangeManager.CreateCurrency(ctx, "COPPER", "Copper", big.NewFloat(1.0), "superman")
	testData := []*ExchangeTest{
		&ExchangeTest{
			from:       "GOLD",
			fromAmount: 1000,
			to:         "PLATINUM",
			toAmount:   100,
		},
		&ExchangeTest{
			from:       "GOLD",
			fromAmount: 1000,
			to:         "SILVER",
			toAmount:   10000,
		},
		&ExchangeTest{
			from:       "GOLD",
			fromAmount: 1000,
			to:         "GOLD",
			toAmount:   1000,
		},
	}
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
