package acccore

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	// DEBIT is enum transaction type DEBIT
	DEBIT Alignment = iota
	// CREDIT is enum transaction type CREDIT
	CREDIT
)

// Alignment is the enum type of transaction type, DEBIT and CREDIT
type Alignment int

// Journal interface define a base Journal structure.
// A journal depict an event where Transactions is happening.
// Important to understand, that Journal don't have update or delete function, its due to accountability reason.
// To delete a journal, one should create a Reversal journal.
// To update a journal, one should create a Reversal journal and then followed with a correction journal.
// If your implementation database do not support 2 phased commit, you should maintain your own committed flag in
// this journal table. When you want to select those journal, you only select those  that have committed flag status on.
// Committing this journal, will propagate to commit the child Transactions
type Journal interface {
	// GetJournalID would return the journal unique ID
	GetJournalID() string
	// SetJournalID will set a new JournalID
	SetJournalID(newID string) Journal

	// GetJournalingTime will return the timestamp of when this journal entry is created
	GetJournalingTime() time.Time
	// SetJournalingTime will set new JournalTime
	SetJournalingTime(newTime time.Time) Journal

	// GetDescription returns Description about this journal entry
	GetDescription() string
	// SetDescription will set new Description
	SetDescription(newDesc string) Journal

	// IsReversal return an indicator if this journal entry is a Reversal of other journal
	IsReversal() bool
	// SetReversal will set new Reversal status
	SetReversal(rev bool) Journal

	// GetReversedJournal should returned the Journal that is reversed IF `IsReverse()` function returned `true`
	GetReversedJournal() Journal
	// SetReversedJournal will set the reversed journal
	SetReversedJournal(journal Journal) Journal

	// GetAmount should return the current Amount of total transaction values
	GetAmount() decimal.Decimal
	// SetAmount will set new total transaction Amount
	SetAmount(newAmount decimal.Decimal) Journal

	// GetTransactions should returns all transaction information that being part of this journal entry.
	GetTransactions() []Transaction
	// SetTransactions will set new list of transaction under this journal
	SetTransactions(transactions []Transaction) Journal

	// GetCreateTime function should return the time when this entry is created/recorded. Logically it the same as `GetTime()` function
	// this function serves as audit trail.
	GetCreateTime() time.Time
	// SetCreateTime will set the creation time
	SetCreateTime(newTime time.Time) Journal

	// GetCreateBy function should return the user AccountNumber or some identification of who is creating this journal.
	// this function serves as audit trail.
	GetCreateBy() string
	// SetCreateBy will set the creator Name
	SetCreateBy(creator string) Journal
}

// Transaction interface define a base Transaction structure
// A transaction is a unit of transaction element that involved within a journal.
// A transaction must include reference to the journal that binds the transaction with other transaction and
// also must state the Account tha doing the transaction
// If your implementation database do not support 2 phased commit, you should maintain your own committed flag in
// this transaction table. When you want to select those transaction, you only select those  that have committed flag status on.
type Transaction interface {
	// GetTransactionID returns the unique ID of this transaction
	GetTransactionID() string
	// SetTransactionID will set new transaction ID
	SetTransactionID(newID string) Transaction

	// GetTransactionTime returns the timestamp of this transaction
	GetTransactionTime() time.Time
	// SetTransactionTime will set new transaction time
	SetTransactionTime(newTime time.Time) Transaction

	// GetAccountNumber return the account number of account ID who owns this transaction
	GetAccountNumber() string
	// SetAccountNumber will set new account number who own this transaction
	SetAccountNumber(number string) Transaction

	// GetJournal returns the journal information where this transaction is recorded.
	GetJournalID() string
	// SetJournal will set the journal to which this transaction is recorded
	SetJournalID(journalID string) Transaction

	// GetDescription return the Description of this Transaction.
	GetDescription() string
	// SetDescription will set the transaction Description
	SetDescription(desc string) Transaction

	// GetAlignment get the transaction type DEBIT or CREDIT
	GetAlignment() Alignment
	// SetAlignment will set the transaction type
	SetAlignment(txType Alignment) Transaction

	// GetAmount return the transaction Amount
	GetAmount() decimal.Decimal
	// SetAmount will set the Amount
	SetAmount(newAmount decimal.Decimal) Transaction

	// GetBookBalance return the Balance of the account at the time when this transaction has been written.
	GetAccountBalance() decimal.Decimal
	// SetAccountBalance will set new account Balance
	SetAccountBalance(newBalance decimal.Decimal) Transaction

	// GetCreateTime function should return the time when this transaction is created/recorded.
	// this function serves as audit trail.
	GetCreateTime() time.Time
	// SetCreateTime will set new creation time
	SetCreateTime(newTime time.Time) Transaction

	// GetCreateBy function should return the user AccountNumber or some identification of who is creating this transaction.
	// this function serves as audit trail.
	GetCreateBy() string
	// SetCreateBy will set new creator Name
	SetCreateBy(creator string) Transaction
}

// Account interface provides base structure of Account
type Account interface {
	// GetCurrency returns the Currency identifier such as `GOLD` or `POINT` or `IDR`
	GetCurrency() string
	// SetCurrency will set the account Currency
	SetCurrency(newCurrency string) Account

	// GetAccountNumber returns the unique account number
	GetAccountNumber() string
	// SetAccountNumber will set new account ID
	SetAccountNumber(newNumber string) Account

	// GetName returns the account Name
	GetName() string
	// SetName will set the new account Name
	SetName(newName string) Account

	// GetDescription returns some Description text about this account
	GetDescription() string
	// SetDescription will set new Description
	SetDescription(newDesc string) Account

	// GetAlignment returns the base transaction type of this account,
	// 1. Asset based should be DEBIT
	// 2. Equity or Liability based should be CREDIT
	GetAlignment() Alignment
	// SetAlignment will set new base transaction type
	SetAlignment(newType Alignment) Account

	// GetBalance returns the current Balance of this account.
	// for each transaction created for this account, this Balance MUST BE UPDATED
	GetBalance() decimal.Decimal
	// SetBalance will set new transaction Balance
	SetBalance(newBalance decimal.Decimal) Account

	// GetCOA returns the COA code for this account, used for categorization of account.
	GetCOA() string
	// SetCOA Will set new COA code
	SetCOA(newCoa string) Account

	// GetCreateTime function should return the time when this account is created/recorded.
	// this function serves as audit trail.
	GetCreateTime() time.Time
	// SetCreateTime will set new creation time
	SetCreateTime(newTime time.Time) Account

	// GetCreateBy function should return the user AccountNumber or some identification of who is creating this account.
	// this function serves as audit trail.
	GetCreateBy() string
	// SetCreateBy will set the creator Name
	SetCreateBy(creator string) Account

	// GetUpdateTime function should return the time when this account is updated.
	// this function serves as audit trail.
	GetUpdateTime() time.Time
	// SetUpdateTime will set the last update time.
	SetUpdateTime(newTime time.Time) Account

	// GetUpdateBy function should return the user AccountNumber or some identification of who is updating this account.
	// this function serves as audit trail.
	GetUpdateBy() string
	// SetUpdateBy will set the updater Name
	SetUpdateBy(editor string) Account
}

// Currency interface provides base structure of Currency
type Currency interface {
	// GetCode get the currency short code. e.g. USD
	GetCode() string
	// SetCode set the currency short code. e.g. USD
	SetCode(code string) Currency

	// GetName get the textual name of the currency. e.g. United States Dollar
	GetName() string
	// SetName set the currency textual name of the currency. e.g. United States Dollar
	SetName(name string) Currency

	// GetExchange get the exchange unit of this currency toward the denominator value
	GetExchange() decimal.Decimal
	// SetExchange set the exchange unit of this currency toward the denominator value
	SetExchange(exchange decimal.Decimal) Currency

	// GetCreateTime function should return the time when this account is created/recorded.
	// this function serves as audit trail.
	GetCreateTime() time.Time
	// SetCreateTime will set new creation time
	SetCreateTime(newTime time.Time) Currency

	// GetCreateBy function should return the user AccountNumber or some identification of who is creating this account.
	// this function serves as audit trail.
	GetCreateBy() string
	// SetCreateBy will set the creator Name
	SetCreateBy(creator string) Currency

	// GetUpdateTime function should return the time when this account is updated.
	// this function serves as audit trail.
	GetUpdateTime() time.Time
	// SetUpdateTime will set the last update time.
	SetUpdateTime(newTime time.Time) Currency

	// GetUpdateBy function should return the user AccountNumber or some identification of who is updating this account.
	// this function serves as audit trail.
	GetUpdateBy() string
	// SetUpdateBy will set the updater Name
	SetUpdateBy(editor string) Currency
}
