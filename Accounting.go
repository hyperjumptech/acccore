package acccore

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
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

// Accounting is the account detail structure
type Accounting struct {
	accountManager     AccountManager
	transactionManager TransactionManager
	journalManager     JournalManager
	uniqueIDGenerator  UniqueIDGenerator
}

// GetAccountManager returns account manager
func (acc *Accounting) GetAccountManager() AccountManager {
	return acc.accountManager
}

// GetTransactionManager returns transaction manager
func (acc *Accounting) GetTransactionManager() TransactionManager {
	return acc.transactionManager
}

// GetJournalManager returns journal manager
func (acc *Accounting) GetJournalManager() JournalManager {
	return acc.journalManager
}

// GetUniqueIDGenerator returns id generator
func (acc *Accounting) GetUniqueIDGenerator() UniqueIDGenerator {
	return acc.uniqueIDGenerator
}

// CreateNewAccount creates a new account
func (acc *Accounting) CreateNewAccount(context context.Context, accountNumber, name, description, coa string, currency string, alignment Alignment, creator string) (Account, error) {
	account := acc.GetAccountManager().NewAccount(context).
		SetName(name).SetDescription(description).SetCOA(coa).
		SetCurrency(currency).SetAlignment(alignment).
		SetCreateBy(creator).SetCreateTime(time.Now())
	if len(accountNumber) == 0 {
		account.SetAccountNumber(acc.GetUniqueIDGenerator().NewUniqueID())
	} else {
		account.SetAccountNumber(accountNumber)
	}
	err := acc.GetAccountManager().PersistAccount(context, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// TransactionInfo transaction info details
type TransactionInfo struct {
	AccountNumber string
	Description   string
	TxType        Alignment
	Amount        decimal.Decimal
}

// CreateNewJournal creates a new journal
func (acc *Accounting) CreateNewJournal(context context.Context, description string, transactions []TransactionInfo, creator string) (Journal, error) {
	journal := acc.GetJournalManager().NewJournal(context).SetDescription(description)

	journal.SetJournalID(acc.GetUniqueIDGenerator().NewUniqueID()).SetCreateBy(creator).
		SetCreateTime(time.Now()).SetJournalingTime(time.Now()).
		SetReversal(false).SetReversedJournal(nil)

	transacs := make([]Transaction, 0)

	// make sure all Transactions have accounts of the same Currency
	for _, txinfo := range transactions {
		newTransaction := acc.GetTransactionManager().NewTransaction(context).SetCreateBy(creator).SetCreateTime(time.Now()).
			SetDescription(txinfo.Description).SetAccountNumber(txinfo.AccountNumber).SetAmount(txinfo.Amount).
			SetTransactionTime(time.Now()).SetAlignment(txinfo.TxType).SetTransactionID(acc.GetUniqueIDGenerator().NewUniqueID())

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
		fmt.Println("No error but no content")
		return nil, fmt.Errorf("commit journal raised no error but no content")
	}
	return journal, nil
}

// CreateReversal creats a reversal
func (acc *Accounting) CreateReversal(context context.Context, description string, reversed Journal, creator string) (Journal, error) {
	journal := acc.GetJournalManager().NewJournal(context).SetDescription(description)
	journal.SetJournalID(acc.GetUniqueIDGenerator().NewUniqueID()).SetCreateBy(creator).SetCreateTime(time.Now()).SetJournalingTime(time.Now()).
		SetReversal(true).SetReversedJournal(reversed)

	transacs := make([]Transaction, 0)

	// make sure all Transactions have accounts of the same Currency
	for _, txinfo := range reversed.GetTransactions() {
		tx := DEBIT
		if txinfo.GetAlignment() == DEBIT {
			tx = CREDIT
		}

		newTransaction := acc.GetTransactionManager().NewTransaction(context).SetCreateBy(creator).SetCreateTime(time.Now()).
			SetDescription(fmt.Sprintf("%s - reversed", txinfo.GetDescription())).SetAccountNumber(txinfo.GetAccountNumber()).
			SetTransactionTime(time.Now()).SetAlignment(tx).SetTransactionID(acc.GetUniqueIDGenerator().NewUniqueID())

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
