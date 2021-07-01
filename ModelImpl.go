package acccore

import (
	"time"
)

// BaseJournal is the base implementation of Journal
type BaseJournal struct {
	journalID       string
	journalingTime  time.Time
	description     string
	reversal        bool
	reversedJournal Journal
	amount          int64
	transactions    []Transaction
	createTime      time.Time
	createdBy       string
}

// GetJournalID would return the journal unique ID
func (journal *BaseJournal) GetJournalID() string {
	return journal.journalID
}

// SetJournalID will set a new JournalID
func (journal *BaseJournal) SetJournalID(newID string) Journal {
	journal.journalID = newID
	return journal
}

// GetJournalingTime will return the timestamp of when this journal entry is created
func (journal *BaseJournal) GetJournalingTime() time.Time {
	return journal.journalingTime
}

// SetJournalingTime will set new JournalTime
func (journal *BaseJournal) SetJournalingTime(newTime time.Time) Journal {
	journal.journalingTime = newTime
	return journal
}

// GetDescription returns description about this journal entry
func (journal *BaseJournal) GetDescription() string {
	return journal.description
}

// SetDescription will set new description
func (journal *BaseJournal) SetDescription(newDesc string) Journal {
	journal.description = newDesc
	return journal
}

// IsReversal return an indicator if this journal entry is a reversal of other journal
func (journal *BaseJournal) IsReversal() bool {
	return journal.reversal
}

// SetReversal will set new reversal status
func (journal *BaseJournal) SetReversal(rev bool) Journal {
	journal.reversal = rev
	return journal
}

// GetReversedJournal should returned the Journal that is reversed IF `IsReverse()` function returned `true`
func (journal *BaseJournal) GetReversedJournal() Journal {
	return journal.reversedJournal
}

// SetReversedJournal will set the reversed journal
func (journal *BaseJournal) SetReversedJournal(reversedJournal Journal) Journal {
	journal.reversedJournal = reversedJournal
	return journal
}

// GetAmount should return the current amount of total transaction values
func (journal *BaseJournal) GetAmount() int64 {
	return journal.amount
}

// SetAmount will set new total transaction amount
func (journal *BaseJournal) SetAmount(newAmount int64) Journal {
	journal.amount = newAmount
	return journal
}

// GetTransactions should returns all transaction information that being part of this journal entry.
func (journal *BaseJournal) GetTransactions() []Transaction {
	return journal.transactions
}

// SetTransactions will set new list of transaction under this journal
func (journal *BaseJournal) SetTransactions(transactions []Transaction) Journal {
	journal.transactions = transactions
	return journal
}

// GetCreateTime function should return the time when this entry is created/recorded. Logically it the same as `GetTime()` function
// this function serves as audit trail.
func (journal *BaseJournal) GetCreateTime() time.Time {
	return journal.createTime
}

// SetCreateTime will set the creation time
func (journal *BaseJournal) SetCreateTime(newTime time.Time) Journal {
	journal.createTime = newTime
	return journal
}

// GetCreateBy function should return the user accountNumber or some identification of who is creating this journal.
// this function serves as audit trail.
func (journal *BaseJournal) GetCreateBy() string {
	return journal.createdBy
}

// SetCreateBy will set the creator name
func (journal *BaseJournal) SetCreateBy(creator string) Journal {
	journal.createdBy = creator
	return journal
}

// BaseTransaction is the base implementation of Transaction
type BaseTransaction struct {
	transactionID   string
	transactionTime time.Time
	accountNumber   string
	journalID       string
	description     string
	transactionType TransactionType
	amount          int64
	accountBalance  int64
	createTime      time.Time
	createBy        string
}

// GetTransactionID returns the unique ID of this transaction
func (trx *BaseTransaction) GetTransactionID() string {
	return trx.transactionID
}

// SetTransactionID will set new transaction ID
func (trx *BaseTransaction) SetTransactionID(newId string) Transaction {
	trx.transactionID = newId
	return trx
}

// GetTransactionTime returns the timestamp of this transaction
func (trx *BaseTransaction) GetTransactionTime() time.Time {
	return trx.transactionTime
}

// SetTransactionTime will set new transaction time
func (trx *BaseTransaction) SetTransactionTime(newTime time.Time) Transaction {
	trx.transactionTime = newTime
	return trx
}

// GetAccountNumber return the account number of account ID who owns this transaction
func (trx *BaseTransaction) GetAccountNumber() string {
	return trx.accountNumber
}

// SetAccountNumber will set new account number who own this transaction
func (trx *BaseTransaction) SetAccountNumber(number string) Transaction {
	trx.accountNumber = number
	return trx
}

// GetJournalID returns the journal information where this transaction is recorded.
func (trx *BaseTransaction) GetJournalID() string {
	return trx.journalID
}

// SetJournalID will set the journal to which this transaction is recorded
func (trx *BaseTransaction) SetJournalID(journalID string) Transaction {
	trx.journalID = journalID
	return trx
}

// GetDescription return the description of this Transaction.
func (trx *BaseTransaction) GetDescription() string {
	return trx.description
}

// SetDescription will set the transaction description
func (trx *BaseTransaction) SetDescription(desc string) Transaction {
	trx.description = desc
	return trx
}

// GetTransactionType get the transaction type DEBIT or CREDIT
func (trx *BaseTransaction) GetTransactionType() TransactionType {
	return trx.transactionType
}

// SetTransactionType will set the transaction type
func (trx *BaseTransaction) SetTransactionType(txType TransactionType) Transaction {
	trx.transactionType = txType
	return trx
}

// GetAmount return the transaction amount
func (trx *BaseTransaction) GetAmount() int64 {
	return trx.amount
}

// SetAmount will set the amount
func (trx *BaseTransaction) SetAmount(newAmount int64) Transaction {
	trx.amount = newAmount
	return trx
}

// GetAccountBalance return the balance of the account at the time when this transaction has been written.
func (trx *BaseTransaction) GetAccountBalance() int64 {
	return trx.accountBalance
}

// SetAccountBalance will set new account balance
func (trx *BaseTransaction) SetAccountBalance(newBalance int64) Transaction {
	trx.accountBalance = newBalance
	return trx
}

// GetCreateTime function should return the time when this transaction is created/recorded.
// this function serves as audit trail.
func (trx *BaseTransaction) GetCreateTime() time.Time {
	return trx.createTime
}

// SetCreateTime will set new creation time
func (trx *BaseTransaction) SetCreateTime(newTime time.Time) Transaction {
	trx.createTime = newTime
	return trx
}

// GetCreateBy function should return the user accountNumber or some identification of who is creating this transaction.
// this function serves as audit trail.
func (trx *BaseTransaction) GetCreateBy() string {
	return trx.createBy
}

// SetCreateBy will set new creator name
func (trx *BaseTransaction) SetCreateBy(creator string) Transaction {
	trx.createBy = creator
	return trx
}

// BaseAccount is the base implementation of Account
type BaseAccount struct {
	currency            string
	accountNumber       string
	name                string
	description         string
	baseTransactionType TransactionType
	balance             int64
	coa                 string
	createTime          time.Time
	createBy            string
	updateTime          time.Time
	updateBy            string
}

// GetCurrency returns the currency identifier such as `GOLD` or `POINT` or `IDR`
func (acc *BaseAccount) GetCurrency() string {
	return acc.currency
}

// SetCurrency will set the account currency
func (acc *BaseAccount) SetCurrency(newCurrency string) Account {
	acc.currency = newCurrency
	return acc
}

// GetAccountNumber returns the unique account number
func (acc *BaseAccount) GetAccountNumber() string {
	return acc.accountNumber
}

// SetAccountNumber will set new account ID
func (acc *BaseAccount) SetAccountNumber(newNumber string) Account {
	acc.accountNumber = newNumber
	return acc
}

// GetName returns the account name
func (acc *BaseAccount) GetName() string {
	return acc.name
}

// SetName will set the new account name
func (acc *BaseAccount) SetName(newName string) Account {
	acc.name = newName
	return acc
}

// GetDescription returns some description text about this account
func (acc *BaseAccount) GetDescription() string {
	return acc.description
}

// SetDescription will set new description
func (acc *BaseAccount) SetDescription(newDesc string) Account {
	acc.description = newDesc
	return acc
}

// GetBaseTransactionType returns the base transaction type of this account,
// 1. Asset based should be DEBIT
// 2. Equity or Liability based should be CREDIT
func (acc *BaseAccount) GetBaseTransactionType() TransactionType {
	return acc.baseTransactionType
}

// SetBaseTransactionType will set new base transaction type
func (acc *BaseAccount) SetBaseTransactionType(newType TransactionType) Account {
	acc.baseTransactionType = newType
	return acc
}

// GetBalance returns the current balance of this account.
// for each transaction created for this account, this balance MUST BE UPDATED
func (acc *BaseAccount) GetBalance() int64 {
	return acc.balance
}

// SetBalance will set new transaction balance
func (acc *BaseAccount) SetBalance(newBalance int64) Account {
	acc.balance = newBalance
	return acc
}

// GetCOA returns the COA code for this account, used for categorization of account.
func (acc *BaseAccount) GetCOA() string {
	return acc.coa
}

// SetCOA Will set new COA code
func (acc *BaseAccount) SetCOA(newCoa string) Account {
	acc.coa = newCoa
	return acc
}

// GetCreateTime function should return the time when this account is created/recorded.
// this function serves as audit trail.
func (acc *BaseAccount) GetCreateTime() time.Time {
	return acc.createTime
}

// SetCreateTime will set new creation time
func (acc *BaseAccount) SetCreateTime(newTime time.Time) Account {
	acc.createTime = newTime
	return acc
}

// GetCreateBy function should return the user accountNumber or some identification of who is creating this account.
// this function serves as audit trail.
func (acc *BaseAccount) GetCreateBy() string {
	return acc.createBy
}

// SetCreateBy will set the creator name
func (acc *BaseAccount) SetCreateBy(creator string) Account {
	acc.createBy = creator
	return acc
}

// GetUpdateTime function should return the time when this account is updated.
// this function serves as audit trail.
func (acc *BaseAccount) GetUpdateTime() time.Time {
	return acc.updateTime
}

// SetUpdateTime will set the last update time.
func (acc *BaseAccount) SetUpdateTime(newTime time.Time) Account {
	acc.updateTime = newTime
	return acc
}

// GetUpdateBy function should return the user accountNumber or some identification of who is updating this account.
// this function serves as audit trail.
func (acc *BaseAccount) GetUpdateBy() string {
	return acc.updateBy
}

// SetUpdateBy will set the updater name
func (acc *BaseAccount) SetUpdateBy(editor string) Account {
	acc.updateBy = editor
	return acc
}
