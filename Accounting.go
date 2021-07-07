package acccore

import (
	"context"
	"fmt"
	"time"
)

// NewAccounting instantiate new Accounting logic modules.
func NewAccounting(accountManager AccountManager, transactionManager TransactionManager, journalManager JournalManager, uniqueIDGenerator UniqueIDGenerator) *Accounting {
	return &Accounting{
		accountManager:     accountManager,
		transactionManager: transactionManager,
		journalManager:     journalManager,
		uniqueIDGenerator:  uniqueIDGenerator,
	}
}

type Accounting struct {
	accountManager     AccountManager
	transactionManager TransactionManager
	journalManager     JournalManager
	uniqueIDGenerator  UniqueIDGenerator
}

func (acc *Accounting) GetAccountManager() AccountManager {
	return acc.accountManager
}
func (acc *Accounting) GetTransactionManager() TransactionManager {
	return acc.transactionManager
}
func (acc *Accounting) GetJournalManager() JournalManager {
	return acc.journalManager
}
func (acc *Accounting) GetUniqueIDGenerator() UniqueIDGenerator {
	return acc.uniqueIDGenerator
}

func (acc *Accounting) CreateNewAccount(context context.Context, name, description, coa string, currency string, alignment TransactionType, creator string) (Account, error) {
	account := acc.GetAccountManager().NewAccount(context).
		SetName(name).SetDescription(description).SetCOA(coa).
		SetCurrency(currency).SetBaseTransactionType(alignment).
		SetAccountNumber(acc.GetUniqueIDGenerator().NewUniqueID()).
		SetCreateBy(creator).SetCreateTime(time.Now())
	err := acc.GetAccountManager().PersistAccount(context, account)
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

func (acc *Accounting) CreateNewJournal(context context.Context, description string, transactions []TransactionInfo, creator string) (Journal, error) {
	journal := acc.GetJournalManager().NewJournal(context).SetDescription(description)

	journal.SetJournalID(acc.GetUniqueIDGenerator().NewUniqueID()).SetCreateBy(creator).
		SetCreateTime(time.Now()).SetJournalingTime(time.Now()).
		SetReversal(false).SetReversedJournal(nil)

	transacs := make([]Transaction, 0)

	// make sure all transactions have accounts of the same Currency
	for _, txinfo := range transactions {
		newTransaction := acc.GetTransactionManager().NewTransaction(context).SetCreateBy(creator).SetCreateTime(time.Now()).
			SetDescription(txinfo.Description).SetAccountNumber(txinfo.AccountNumber).SetAmount(txinfo.Amount).
			SetTransactionTime(time.Now()).SetTransactionType(txinfo.TxType).SetTransactionID(acc.GetUniqueIDGenerator().NewUniqueID())

		transacs = append(transacs, newTransaction)
	}

	journal.SetTransactions(transacs)

	err := acc.GetJournalManager().PersistJournal(context, journal)
	if err != nil {
		err = acc.GetJournalManager().CommitJournal(context, journal)
		if err != nil {
			err = acc.GetJournalManager().CancelJournal(context, journal)
			return nil, err
		}
		return nil, err
	}
	return journal, nil
}

func (acc *Accounting) CreateReversal(context context.Context, description string, reversed Journal, creator string) (Journal, error) {
	journal := acc.GetJournalManager().NewJournal(context).SetDescription(description)
	journal.SetJournalID(acc.GetUniqueIDGenerator().NewUniqueID()).SetCreateBy(creator).SetCreateTime(time.Now()).SetJournalingTime(time.Now()).
		SetReversal(true).SetReversedJournal(reversed)

	transacs := make([]Transaction, 0)

	// make sure all transactions have accounts of the same Currency
	for _, txinfo := range reversed.GetTransactions() {
		tx := DEBIT
		if txinfo.GetTransactionType() == DEBIT {
			tx = CREDIT
		}

		newTransaction := acc.GetTransactionManager().NewTransaction(context).SetCreateBy(creator).SetCreateTime(time.Now()).
			SetDescription(fmt.Sprintf("%s - reversed", txinfo.GetDescription())).SetAccountNumber(txinfo.GetAccountNumber()).
			SetTransactionTime(time.Now()).SetTransactionType(tx).SetTransactionID(acc.GetUniqueIDGenerator().NewUniqueID())

		transacs = append(transacs, newTransaction)
	}

	journal.SetTransactions(transacs)

	err := acc.GetJournalManager().PersistJournal(context, journal)
	if err != nil {
		err = acc.GetJournalManager().CommitJournal(context, journal)
		if err != nil {
			err = acc.GetJournalManager().CancelJournal(context, journal)
			return nil, err
		}
		return nil, err
	}
	return journal, nil
}
