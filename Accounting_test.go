package acccore

import (
	"testing"
	"time"
)

func TestAccounting_CreateNewAccount(t *testing.T) {
	ClearInMemoryTables()

	acc := &Accounting{
		accountManager:     &InMemoryAccountManager{},
		transactionManager: &InMemoryTransactionManager{},
		journalManager:     &InMemoryJournalManager{},
		uniqueIDGenerator: &RandomGenUniqueIDGenerator{
			Length:        10,
			LowerAlpha:    false,
			UpperAlpha:    true,
			Numeric:       true,
			CharSetBuffer: nil,
		},
	}

	account, err := acc.CreateNewAccount("Test Account", "Gold base test user account", "1.1", "GOLD", CREDIT, "aCreator")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	exist, err := acc.accountManager.IsAccountIdExist(account.GetAccountNumber())
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	if !exist {
		t.Error("account should exist after creation")
		t.FailNow()
	}
	render, err := acc.transactionManager.RenderTransactionsOnAccount(time.Now().Add(-2*time.Hour), time.Now().Add(2*time.Hour), account, PageRequest{
		PageNo:   1,
		ItemSize: 10,
		Sorts:    nil,
	})
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	} else {
		t.Log(render)
	}

}

func TestAccounting_CreateNewJournal(t *testing.T) {
	ClearInMemoryTables()

	acc := &Accounting{
		accountManager:     &InMemoryAccountManager{},
		transactionManager: &InMemoryTransactionManager{},
		journalManager:     &InMemoryJournalManager{},
		uniqueIDGenerator: &RandomGenUniqueIDGenerator{
			Length:        10,
			LowerAlpha:    false,
			UpperAlpha:    true,
			Numeric:       true,
			CharSetBuffer: nil,
		},
	}

	goldLoan, err := acc.CreateNewAccount("Gold Loan", "Gold base loan reserve", "1.1", "GOLD", DEBIT, "aCreator")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	alphaCreditor, err := acc.CreateNewAccount("Gold Creditor Alpha", "Gold base debitor alpha", "2.1", "GOLD", CREDIT, "aCreator")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	betaDebitor, err := acc.CreateNewAccount("Gold Debitor Alpha", "Gold base creditor beta", "3.1", "GOLD", DEBIT, "aCreator")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	topupTransactions := []TransactionInfo{
		{
			AccountNumber: goldLoan.GetAccountNumber(),
			Description:   "Added Gold Reserve",
			TxType:        DEBIT,
			Amount:        1000000,
		},
		{
			AccountNumber: alphaCreditor.GetAccountNumber(),
			Description:   "Added Gold Equity",
			TxType:        CREDIT,
			Amount:        1000000,
		},
	}
	journal, err := acc.CreateNewJournal("Creditor Topup Gold", topupTransactions, "aCreator")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	t.Log(acc.journalManager.RenderJournal(journal))

	goldPurchaseTransaction := []TransactionInfo{
		{
			AccountNumber: betaDebitor.GetAccountNumber(),
			Description:   "Add debitor AR",
			TxType:        DEBIT,
			Amount:        200000,
		},
		{
			AccountNumber: goldLoan.GetAccountNumber(),
			Description:   "Gold Disbursement",
			TxType:        CREDIT,
			Amount:        200000,
		},
	}
	journal, err = acc.CreateNewJournal("GOLD purchase transaction", goldPurchaseTransaction, "aCreator")
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	}
	t.Log(acc.journalManager.RenderJournal(journal))

	render, err := acc.transactionManager.RenderTransactionsOnAccount(time.Now().Add(-2*time.Hour), time.Now().Add(2*time.Hour), goldLoan, PageRequest{
		PageNo:   1,
		ItemSize: 10,
		Sorts:    nil,
	})
	if err != nil {
		t.Error(err.Error())
		t.FailNow()
	} else {
		t.Log(render)
	}
}
