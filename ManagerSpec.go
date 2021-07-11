package acccore

import (
	"context"
	"fmt"
	"math/big"
	"time"
)

var (
	ErrJournalNil                          = fmt.Errorf("journal is nil")
	ErrJournalMissingId                    = fmt.Errorf("journal is missing accountNumber")
	ErrJournalNoTransaction                = fmt.Errorf("journal contains no transactions")
	ErrJournalMissingAuthor                = fmt.Errorf("journal author is not known")
	ErrJournalAlreadyPersisted             = fmt.Errorf("journal is already persisted")
	ErrJournalTransactionAlreadyPersisted  = fmt.Errorf("journal transaction is already persisted")
	ErrJournalTransactionMissingID         = fmt.Errorf("journal transactions missing accountNumber")
	ErrJournalNotBalance                   = fmt.Errorf("journal's sum of debit and sum of credit do not balance")
	ErrJournalTransactionMixCurrency       = fmt.Errorf("journal transactions contains mixed currencies, all transaction in a journal must belong to the same currency")
	ErrJournalTransactionAccountNotPersist = fmt.Errorf("journal transactions revering to non-existent account")
	ErrJournalTransactionAccountDuplicate  = fmt.Errorf("multiple journal transactions belongs to the same account")
	ErrJournalIdNotFound                   = fmt.Errorf("journal with specified ID not in database")
	ErrJournalLoadReversalInconsistent     = fmt.Errorf("reversed journal reverence to unexistent journal")
	ErrJournalCanNotDoubleReverse          = fmt.Errorf("journal can only reversed once")

	ErrAccountAlreadyPersisted   = fmt.Errorf("account is already persisted")
	ErrAccountIsNotPersisted     = fmt.Errorf("account is not persisted")
	ErrAccountIdNotFound         = fmt.Errorf("account accountNumber not in database")
	ErrAccountMissingID          = fmt.Errorf("account accountNumber or number is not provided")
	ErrAccountMissingName        = fmt.Errorf("account name is not provided")
	ErrAccountMissingDescription = fmt.Errorf("account description is not provided")
	ErrAccountMissingCreator     = fmt.Errorf("account creator is not provided")

	ErrTransactionNotFound = fmt.Errorf("transaction accountNumber not in database")

	ErrCurrencyNotFound = fmt.Errorf("currency not found")
)

// JournalManager is interface used of managing journals
type JournalManager interface {
	// NewJournal will create new blank un-persisted journal
	NewJournal(context context.Context) Journal

	// PersistJournal will record a journal entry into database.
	// It requires list of transactions for which each of the transaction MUST BE :
	//    1.NOT BE PERSISTED. (the journal accountNumber is not exist in DB yet)
	//    2.Pointing or owned by a PERSISTED Account
	//    3.Each of this account must belong to the same Currency
	//    4.Balanced. The total sum of DEBIT and total sum of CREDIT is equal.
	//    5.No duplicate transaction that belongs to the same Account.
	// If your database support 2 phased commit, you can make all balance changes in
	// accounts and transactions. If your db do not support this, you can implement your own 2 phase commits mechanism
	// on the CommitJournal and CancelJournal
	PersistJournal(context context.Context, journalToPersist Journal) error

	// CommitJournal will commit the journal into the system
	// Only non committed journal can be committed.
	// use this if the implementation database do not support 2 phased commit.
	// if your database support 2 phased commit, you should do all commit in the PersistJournal function
	// and this function should simply return nil.
	CommitJournal(context context.Context, journalToCommit Journal) error

	// CancelJournal Cancel a journal
	// Only non committed journal can be committed.
	// use this if the implementation database do not support 2 phased commit.
	// if your database do not support 2 phased commit, you should do all roll back in the PersistJournal function
	// and this function should simply return nil.
	CancelJournal(context context.Context, journalToCancel Journal) error

	// IsJournalIdReversed check if the journal with specified ID has been reversed
	IsJournalIdReversed(context context.Context, journalId string) (bool, error)

	// IsTransactionIdExist will check if an Transaction ID/number is exist in the database.
	IsJournalIdExist(context context.Context, journalId string) (bool, error)

	// GetJournalById retrieved a Journal information identified by its ID.
	// the provided ID must be exactly the same, not uses the LIKE select expression.
	GetJournalById(context context.Context, journalId string) (Journal, error)

	// ListJournals retrieve list of journals with transaction date between the `from` and `until` time range inclusive.
	// This function uses pagination.
	ListJournals(context context.Context, from time.Time, until time.Time, request PageRequest) (PageResult, []Journal, error)

	// RenderJournal Render this journal into string for easy inspection
	RenderJournal(context context.Context, journal Journal) string
}

// TransactionManager is interface used for managing transaction data/table
type TransactionManager interface {
	// NewTransaction will create new blank un-persisted Transaction
	NewTransaction(context context.Context) Transaction

	// IsTransactionIdExist will check if an Transaction ID/number is exist in the database.
	IsTransactionIdExist(context context.Context, id string) (bool, error)

	// GetTransactionById will retrieve one single transaction that identified by some ID
	GetTransactionById(context context.Context, id string) (Transaction, error)

	// ListTransactionsWithAccount retrieves list of transactions that belongs to this account
	// that transaction happens between the `from` and `until` time range.
	// This function uses pagination
	ListTransactionsOnAccount(context context.Context, from time.Time, until time.Time, account Account, request PageRequest) (PageResult, []Transaction, error)

	// RenderTransactionsOnAccount Render list of transaction been down on an account in a time span
	RenderTransactionsOnAccount(context context.Context, from time.Time, until time.Time, account Account, request PageRequest) (string, error)
}

// AccountManager interface is used for managing Accounts
type AccountManager interface {
	// NewAccount will create a new blank un-persisted account.
	NewAccount(context context.Context) Account

	// PersistAccount will save the account into database.
	// will throw error if the account already persisted
	PersistAccount(context context.Context, AccountToPersist Account) error

	// UpdateAccount will update the account database to reflect to the provided account information.
	// This update account function will fail if the account ID/number is not existing in the database.
	UpdateAccount(context context.Context, AccountToUpdate Account) error

	// IsAccountIdExist will check if an account ID/number is exist in the database.
	IsAccountIdExist(context context.Context, id string) (bool, error)

	// GetAccountById retrieve an account information by specifying the ID/number
	GetAccountById(context context.Context, id string) (Account, error)

	// ListAccounts list all account in the database.
	// This function uses pagination
	ListAccounts(context context.Context, request PageRequest) (PageResult, []Account, error)

	// ListAccountByCOA returns list of accounts that have the same COA number.
	// This function uses pagination
	ListAccountByCOA(context context.Context, coa string, request PageRequest) (PageResult, []Account, error)

	// FindAccounts returns list of accounts that have their name contains a substring of specified parameter.
	// this search should  be case insensitive.
	FindAccounts(context context.Context, nameLike string, request PageRequest) (PageResult, []Account, error)
}

// ExchangeManager will define functions to be implemented for currency exchanges.
// this interface follows the exchange mechanism using a common denominator.
type ExchangeManager interface {
	// IsCurrencyExist will check in the exchange system for a currency existance
	// non-existent currency means that the currency is not supported.
	// error should be thrown if only there's an underlying error such as db error.
	IsCurrencyExist(currency string) (bool, error)
	// GetDenom get the current common denominator used in the exchange
	GetDenom(context context.Context) *big.Float
	// SetDenom set the current common denominator value into the specified value
	SetDenom(context context.Context, denom *big.Float)

	// SetExchangeValueOf set the specified value as denominator value for that speciffic currency.
	// This function should return error if the currency specified is not exist.
	SetExchangeValueOf(context context.Context, currency string, exchange *big.Float) error
	// GetExchangeValueOf get the denominator value of the specified currency.
	// Error should be returned if the specified currency is not exist.
	GetExchangeValueOf(context context.Context, currency string) (*big.Float, error)

	// Get the currency exchange rate for exchanging between the two currency.
	// if any of the currency is not exist, an error should be returned.
	// if from and to currency is equal, this must return 1.0
	CalculateExchangeRate(context context.Context, fromCurrency, toCurrency string) (*big.Float, error)
	// Get the currency exchange value for the amount of fromCurrency into toCurrency.
	// If any of the currency is not exist, an error should be returned.
	// if from and to currency is equal, the returned amount must be equal to the amount in the argument.
	CalculateExchange(context context.Context, fromCurrency, toCurrency string, amount int64) (int64, error)
}
