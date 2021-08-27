package acccore

import (
	"time"
)

// BaseJournal is the base implementation of Journal
type BaseJournal struct {
	JournalID       string        `json:"journal_id"`
	JournalingTime  time.Time     `json:"journaling_time"`
	Description     string        `json:"description"`
	Reversal        bool          `json:"reversal"`
	ReversedJournal Journal       `json:"reversed_journal"`
	Amount          int64         `json:"amount"`
	Transactions    []Transaction `json:"transactions"`
	CreateTime      time.Time     `json:"create_time"`
	CreatedBy       string        `json:"created_by"`
}

// GetJournalID would return the journal unique ID
func (journal *BaseJournal) GetJournalID() string {
	return journal.JournalID
}

// SetJournalID will set a new JournalID
func (journal *BaseJournal) SetJournalID(newID string) Journal {
	journal.JournalID = newID
	return journal
}

// GetJournalingTime will return the timestamp of when this journal entry is created
func (journal *BaseJournal) GetJournalingTime() time.Time {
	return journal.JournalingTime
}

// SetJournalingTime will set new JournalTime
func (journal *BaseJournal) SetJournalingTime(newTime time.Time) Journal {
	journal.JournalingTime = newTime
	return journal
}

// GetDescription returns Description about this journal entry
func (journal *BaseJournal) GetDescription() string {
	return journal.Description
}

// SetDescription will set new Description
func (journal *BaseJournal) SetDescription(newDesc string) Journal {
	journal.Description = newDesc
	return journal
}

// IsReversal return an indicator if this journal entry is a Reversal of other journal
func (journal *BaseJournal) IsReversal() bool {
	return journal.Reversal
}

// SetReversal will set new Reversal status
func (journal *BaseJournal) SetReversal(rev bool) Journal {
	journal.Reversal = rev
	return journal
}

// GetReversedJournal should returned the Journal that is reversed IF `IsReverse()` function returned `true`
func (journal *BaseJournal) GetReversedJournal() Journal {
	return journal.ReversedJournal
}

// SetReversedJournal will set the reversed journal
func (journal *BaseJournal) SetReversedJournal(reversedJournal Journal) Journal {
	journal.ReversedJournal = reversedJournal
	return journal
}

// GetAmount should return the current Amount of total transaction values
func (journal *BaseJournal) GetAmount() int64 {
	return journal.Amount
}

// SetAmount will set new total transaction Amount
func (journal *BaseJournal) SetAmount(newAmount int64) Journal {
	journal.Amount = newAmount
	return journal
}

// GetTransactions should returns all transaction information that being part of this journal entry.
func (journal *BaseJournal) GetTransactions() []Transaction {
	return journal.Transactions
}

// SetTransactions will set new list of transaction under this journal
func (journal *BaseJournal) SetTransactions(transactions []Transaction) Journal {
	journal.Transactions = transactions
	return journal
}

// GetCreateTime function should return the time when this entry is created/recorded. Logically it the same as `GetTime()` function
// this function serves as audit trail.
func (journal *BaseJournal) GetCreateTime() time.Time {
	return journal.CreateTime
}

// SetCreateTime will set the creation time
func (journal *BaseJournal) SetCreateTime(newTime time.Time) Journal {
	journal.CreateTime = newTime
	return journal
}

// GetCreateBy function should return the user AccountNumber or some identification of who is creating this journal.
// this function serves as audit trail.
func (journal *BaseJournal) GetCreateBy() string {
	return journal.CreatedBy
}

// SetCreateBy will set the creator Name
func (journal *BaseJournal) SetCreateBy(creator string) Journal {
	journal.CreatedBy = creator
	return journal
}

// BaseTransaction is the base implementation of Transaction
type BaseTransaction struct {
	TransactionID   string    `json:"transaction_id"`
	TransactionTime time.Time `json:"transaction_time"`
	AccountNumber   string    `json:"account_number"`
	JournalID       string    `json:"journal_id"`
	Description     string    `json:"description"`
	TransactionType Alignment `json:"transaction_type"`
	Amount          int64     `json:"amount"`
	AccountBalance  int64     `json:"account_balance"`
	CreateTime      time.Time `json:"create_time"`
	CreateBy        string    `json:"create_by"`
}

// GetTransactionID returns the unique ID of this transaction
func (trx *BaseTransaction) GetTransactionID() string {
	return trx.TransactionID
}

// SetTransactionID will set new transaction ID
func (trx *BaseTransaction) SetTransactionID(newID string) Transaction {
	trx.TransactionID = newID
	return trx
}

// GetTransactionTime returns the timestamp of this transaction
func (trx *BaseTransaction) GetTransactionTime() time.Time {
	return trx.TransactionTime
}

// SetTransactionTime will set new transaction time
func (trx *BaseTransaction) SetTransactionTime(newTime time.Time) Transaction {
	trx.TransactionTime = newTime
	return trx
}

// GetAccountNumber return the account number of account ID who owns this transaction
func (trx *BaseTransaction) GetAccountNumber() string {
	return trx.AccountNumber
}

// SetAccountNumber will set new account number who own this transaction
func (trx *BaseTransaction) SetAccountNumber(number string) Transaction {
	trx.AccountNumber = number
	return trx
}

// GetJournalID returns the journal information where this transaction is recorded.
func (trx *BaseTransaction) GetJournalID() string {
	return trx.JournalID
}

// SetJournalID will set the journal to which this transaction is recorded
func (trx *BaseTransaction) SetJournalID(journalID string) Transaction {
	trx.JournalID = journalID
	return trx
}

// GetDescription return the Description of this Transaction.
func (trx *BaseTransaction) GetDescription() string {
	return trx.Description
}

// SetDescription will set the transaction Description
func (trx *BaseTransaction) SetDescription(desc string) Transaction {
	trx.Description = desc
	return trx
}

// GetAlignment get the transaction type DEBIT or CREDIT
func (trx *BaseTransaction) GetAlignment() Alignment {
	return trx.TransactionType
}

// SetAlignment will set the transaction type
func (trx *BaseTransaction) SetAlignment(txType Alignment) Transaction {
	trx.TransactionType = txType
	return trx
}

// GetAmount return the transaction Amount
func (trx *BaseTransaction) GetAmount() int64 {
	return trx.Amount
}

// SetAmount will set the Amount
func (trx *BaseTransaction) SetAmount(newAmount int64) Transaction {
	trx.Amount = newAmount
	return trx
}

// GetAccountBalance return the Balance of the account at the time when this transaction has been written.
func (trx *BaseTransaction) GetAccountBalance() int64 {
	return trx.AccountBalance
}

// SetAccountBalance will set new account Balance
func (trx *BaseTransaction) SetAccountBalance(newBalance int64) Transaction {
	trx.AccountBalance = newBalance
	return trx
}

// GetCreateTime function should return the time when this transaction is created/recorded.
// this function serves as audit trail.
func (trx *BaseTransaction) GetCreateTime() time.Time {
	return trx.CreateTime
}

// SetCreateTime will set new creation time
func (trx *BaseTransaction) SetCreateTime(newTime time.Time) Transaction {
	trx.CreateTime = newTime
	return trx
}

// GetCreateBy function should return the user AccountNumber or some identification of who is creating this transaction.
// this function serves as audit trail.
func (trx *BaseTransaction) GetCreateBy() string {
	return trx.CreateBy
}

// SetCreateBy will set new creator Name
func (trx *BaseTransaction) SetCreateBy(creator string) Transaction {
	trx.CreateBy = creator
	return trx
}

// BaseAccount is the base implementation of Account
type BaseAccount struct {
	Currency      string    `json:"currency"`
	AccountNumber string    `json:"account_number"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Alignment     Alignment `json:"alignment"`
	Balance       int64     `json:"balance"`
	COA           string    `json:"coa"`
	CreateTime    time.Time `json:"create_time"`
	CreateBy      string    `json:"create_by"`
	UpdateTime    time.Time `json:"update_time"`
	UpdateBy      string    `json:"update_by"`
}

// GetCurrency returns the Currency identifier such as `GOLD` or `POINT` or `IDR`
func (acc *BaseAccount) GetCurrency() string {
	return acc.Currency
}

// SetCurrency will set the account Currency
func (acc *BaseAccount) SetCurrency(newCurrency string) Account {
	acc.Currency = newCurrency
	return acc
}

// GetAccountNumber returns the unique account number
func (acc *BaseAccount) GetAccountNumber() string {
	return acc.AccountNumber
}

// SetAccountNumber will set new account ID
func (acc *BaseAccount) SetAccountNumber(newNumber string) Account {
	acc.AccountNumber = newNumber
	return acc
}

// GetName returns the account Name
func (acc *BaseAccount) GetName() string {
	return acc.Name
}

// SetName will set the new account Name
func (acc *BaseAccount) SetName(newName string) Account {
	acc.Name = newName
	return acc
}

// GetDescription returns some Description text about this account
func (acc *BaseAccount) GetDescription() string {
	return acc.Description
}

// SetDescription will set new Description
func (acc *BaseAccount) SetDescription(newDesc string) Account {
	acc.Description = newDesc
	return acc
}

// GetAlignment returns the base transaction type of this account,
// 1. Asset based should be DEBIT
// 2. Equity or Liability based should be CREDIT
func (acc *BaseAccount) GetAlignment() Alignment {
	return acc.Alignment
}

// SetAlignment will set new base transaction type
func (acc *BaseAccount) SetAlignment(newType Alignment) Account {
	acc.Alignment = newType
	return acc
}

// GetBalance returns the current Balance of this account.
// for each transaction created for this account, this Balance MUST BE UPDATED
func (acc *BaseAccount) GetBalance() int64 {
	return acc.Balance
}

// SetBalance will set new transaction Balance
func (acc *BaseAccount) SetBalance(newBalance int64) Account {
	acc.Balance = newBalance
	return acc
}

// GetCOA returns the COA code for this account, used for categorization of account.
func (acc *BaseAccount) GetCOA() string {
	return acc.COA
}

// SetCOA Will set new COA code
func (acc *BaseAccount) SetCOA(newCoa string) Account {
	acc.COA = newCoa
	return acc
}

// GetCreateTime function should return the time when this account is created/recorded.
// this function serves as audit trail.
func (acc *BaseAccount) GetCreateTime() time.Time {
	return acc.CreateTime
}

// SetCreateTime will set new creation time
func (acc *BaseAccount) SetCreateTime(newTime time.Time) Account {
	acc.CreateTime = newTime
	return acc
}

// GetCreateBy function should return the user AccountNumber or some identification of who is creating this account.
// this function serves as audit trail.
func (acc *BaseAccount) GetCreateBy() string {
	return acc.CreateBy
}

// SetCreateBy will set the creator Name
func (acc *BaseAccount) SetCreateBy(creator string) Account {
	acc.CreateBy = creator
	return acc
}

// GetUpdateTime function should return the time when this account is updated.
// this function serves as audit trail.
func (acc *BaseAccount) GetUpdateTime() time.Time {
	return acc.UpdateTime
}

// SetUpdateTime will set the last update time.
func (acc *BaseAccount) SetUpdateTime(newTime time.Time) Account {
	acc.UpdateTime = newTime
	return acc
}

// GetUpdateBy function should return the user AccountNumber or some identification of who is updating this account.
// this function serves as audit trail.
func (acc *BaseAccount) GetUpdateBy() string {
	return acc.UpdateBy
}

// SetUpdateBy will set the updater Name
func (acc *BaseAccount) SetUpdateBy(editor string) Account {
	acc.UpdateBy = editor
	return acc
}

type BaseCurrency struct {
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Exchange   float64   `json:"exchange"`
	CreateTime time.Time `json:"create_time"`
	CreateBy   string    `json:"create_by"`
	UpdateTime time.Time `json:"update_time"`
	UpdateBy   string    `json:"update_by"`
}

// GetCode get the currency short code. e.g. USD
func (bc *BaseCurrency) GetCode() string {
	return bc.Code
}

// SetCode set the currency short code. e.g. USD
func (bc *BaseCurrency) SetCode(code string) Currency {
	bc.Code = code
	return bc
}

// GetName get the textual name of the currency. e.g. United States Dollar
func (bc *BaseCurrency) GetName() string {
	return bc.Name
}

// SetName set the currency textual name of the currency. e.g. United States Dollar
func (bc *BaseCurrency) SetName(name string) Currency {
	bc.Name = name
	return bc
}

// GetExchange get the exchange unit of this currency toward the denominator value
func (bc *BaseCurrency) GetExchange() float64 {
	return bc.Exchange
}

// SetExchange set the exchange unit of this currency toward the denominator value
func (bc *BaseCurrency) SetExchange(exchange float64) Currency {
	bc.Exchange = exchange
	return bc
}

// GetCreateTime function should return the time when this account is created/recorded.
// this function serves as audit trail.
func (bc *BaseCurrency) GetCreateTime() time.Time {
	return bc.CreateTime
}

// SetCreateTime will set new creation time
func (bc *BaseCurrency) SetCreateTime(newTime time.Time) Currency {
	bc.CreateTime = newTime
	return bc
}

// GetCreateBy function should return the user AccountNumber or some identification of who is creating this account.
// this function serves as audit trail.
func (bc *BaseCurrency) GetCreateBy() string {
	return bc.CreateBy
}

// SetCreateBy will set the creator Name
func (bc *BaseCurrency) SetCreateBy(creator string) Currency {
	bc.CreateBy = creator
	return bc
}

// GetUpdateTime function should return the time when this account is updated.
// this function serves as audit trail.
func (bc *BaseCurrency) GetUpdateTime() time.Time {
	return bc.UpdateTime
}

// SetUpdateTime will set the last update time.
func (bc *BaseCurrency) SetUpdateTime(newTime time.Time) Currency {
	bc.UpdateTime = newTime
	return bc
}

// GetUpdateBy function should return the user AccountNumber or some identification of who is updating this account.
// this function serves as audit trail.
func (bc *BaseCurrency) GetUpdateBy() string {
	return bc.UpdateBy
}

// SetUpdateBy will set the updater Name
func (bc *BaseCurrency) SetUpdateBy(editor string) Currency {
	bc.UpdateBy = editor
	return bc
}
