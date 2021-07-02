package acccore

import (
	"fmt"
	"time"
)

type Accounting struct {
	accountManager     AccountManager
	transactionManager TransactionManager
	journalManager     JournalManager
	uniqueIDGenerator  UniqueIDGenerator
}

func (acc *Accounting) CreateNewAccount(name, description, coa string, currency string, alignment TransactionType, creator string) (Account, error) {
	account := acc.accountManager.NewAccount().
		SetName(name).SetDescription(description).SetCOA(coa).
		SetCurrency(currency).SetBaseTransactionType(alignment).
		SetAccountNumber(acc.uniqueIDGenerator.NewUniqueID()).
		SetCreateBy(creator).SetCreateTime(time.Now())
	err := acc.accountManager.PersistAccount(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

type TransactionInfo struct {
	AccountNumber string
	Description   string
	TxType        TransactionType
	Amount        int64
}

func (acc *Accounting) CreateNewJournal(description string, transactions []TransactionInfo, creator string) (Journal, error) {
	journal := acc.journalManager.NewJournal().SetDescription(description)

	journal.SetJournalID(acc.uniqueIDGenerator.NewUniqueID()).SetCreateBy(creator).
		SetCreateTime(time.Now()).SetJournalingTime(time.Now()).
		SetReversal(false).SetReversedJournal(nil)

	transacs := make([]Transaction, 0)

	// make sure all transactions have accounts of the same Currency
	for _, txinfo := range transactions {
		newTransaction := acc.transactionManager.NewTransaction().SetCreateBy(creator).SetCreateTime(time.Now()).
			SetDescription(txinfo.Description).SetAccountNumber(txinfo.AccountNumber).SetAmount(txinfo.Amount).
			SetTransactionTime(time.Now()).SetTransactionType(txinfo.TxType).SetTransactionID(acc.uniqueIDGenerator.NewUniqueID())

		transacs = append(transacs, newTransaction)
	}

	journal.SetTransactions(transacs)

	err := acc.journalManager.PersistJournal(journal)
	if err != nil {
		err = acc.journalManager.CommitJournal(journal)
		if err != nil {
			err = acc.journalManager.CancelJournal(journal)
			return nil, err
		}
		return nil, err
	}
	return journal, nil
}

func (acc *Accounting) CreateReversal(description string, reversed Journal, creator string) (Journal, error) {
	journal := acc.journalManager.NewJournal().SetDescription(description)
	journal.SetJournalID(acc.uniqueIDGenerator.NewUniqueID()).SetCreateBy(creator).SetCreateTime(time.Now()).SetJournalingTime(time.Now()).
		SetReversal(true).SetReversedJournal(reversed)

	transacs := make([]Transaction, 0)

	// make sure all transactions have accounts of the same Currency
	for _, txinfo := range reversed.GetTransactions() {
		tx := DEBIT
		if txinfo.GetTransactionType() == DEBIT {
			tx = CREDIT
		}

		newTransaction := acc.transactionManager.NewTransaction().SetCreateBy(creator).SetCreateTime(time.Now()).
			SetDescription(fmt.Sprintf("%s - reversed", txinfo.GetDescription())).SetAccountNumber(txinfo.GetAccountNumber()).
			SetTransactionTime(time.Now()).SetTransactionType(tx).SetTransactionID(acc.uniqueIDGenerator.NewUniqueID())

		transacs = append(transacs, newTransaction)
	}

	journal.SetTransactions(transacs)

	err := acc.journalManager.PersistJournal(journal)
	if err != nil {
		err = acc.journalManager.CommitJournal(journal)
		if err != nil {
			err = acc.journalManager.CancelJournal(journal)
			return nil, err
		}
		return nil, err
	}
	return journal, nil
}
