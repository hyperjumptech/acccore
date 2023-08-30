package acccore

import (
	"context"
	"github.com/shopspring/decimal"
	"testing"
)

type ExchangeTest struct {
	from       string
	fromAmount decimal.Decimal
	to         string
	toAmount   decimal.Decimal
}

func TestInMemoryExchangeManager_CalculateExchange(t *testing.T) {
	ctx := context.Background()
	exchangeManager := NewInMemoryExchangeManager()
	_, _ = exchangeManager.CreateCurrency(ctx, "PLATINUM", "Platinum", decimal.NewFromFloat(0.001), "superman")
	_, _ = exchangeManager.CreateCurrency(ctx, "GOLD", "Gold", decimal.NewFromFloat(0.01), "superman")
	_, _ = exchangeManager.CreateCurrency(ctx, "SILVER", "Silver", decimal.NewFromFloat(0.1), "superman")
	_, _ = exchangeManager.CreateCurrency(ctx, "COPPER", "Copper", decimal.NewFromFloat(1.0), "superman")
	testData := []*ExchangeTest{
		{
			from:       "GOLD",
			fromAmount: decimal.NewFromInt(1000),
			to:         "PLATINUM",
			toAmount:   decimal.NewFromInt(100),
		},
		{
			from:       "GOLD",
			fromAmount: decimal.NewFromInt(1000),
			to:         "SILVER",
			toAmount:   decimal.NewFromInt(10000),
		},
		{
			from:       "GOLD",
			fromAmount: decimal.NewFromInt(1000),
			to:         "GOLD",
			toAmount:   decimal.NewFromInt(1000),
		},
	}
	for idx, data := range testData {
		result, err := exchangeManager.CalculateExchange(ctx, data.from, data.to, data.fromAmount)
		if err != nil {
			t.Errorf(err.Error())
			t.Fail()
		}
		if result.InexactFloat64() != data.toAmount.InexactFloat64() {
			t.Errorf("#%d. Expect %f but %f", idx, data.toAmount.InexactFloat64(), result.InexactFloat64())
			t.Fail()
		}
	}

	exchangeManager.SetDenom(ctx, decimal.NewFromFloat(123.456))
	for idx, data := range testData {
		result, err := exchangeManager.CalculateExchange(ctx, data.from, data.to, data.fromAmount)
		if err != nil {
			t.Errorf(err.Error())
			t.Fail()
		}
		if result.InexactFloat64() != data.toAmount.InexactFloat64() {
			t.Errorf("#%d. Expect %f but %f", idx, data.toAmount.InexactFloat64(), result.InexactFloat64())
			t.Fail()
		}
	}
}
