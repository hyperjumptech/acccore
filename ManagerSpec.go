package acccore

import (
	"fmt"
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
)

// JournalManager is interface used of managing journals
type JournalManager interface {
	// NewJournal will create new blank un-persisted journal
	NewJournal() Journal

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
	PersistJournal(journalToPersist Journal) error

	// CommitJournal will commit the journal into the system
	// Only non committed journal can be committed.
	// use this if the implementation database do not support 2 phased commit.
	// if your database support 2 phased commit, you should do all commit in the PersistJournal function
	// and this function should simply return nil.
	CommitJournal(journalToCommit Journal) error

	// CancelJournal Cancel a journal
	// Only non committed journal can be committed.
	// use this if the implementation database do not support 2 phased commit.
	// if your database do not support 2 phased commit, you should do all roll back in the PersistJournal function
	// and this function should simply return nil.
	CancelJournal(journalToCancel Journal) error

	// IsJournalIdReversed check if the journal with specified ID has been reversed
	IsJournalIdReversed(journalId string) (bool, error)

	// IsTransactionIdExist will check if an Transaction ID/number is exist in the database.
	IsJournalIdExist(journalId string) (bool, error)

	// GetJournalById retrieved a Journal information identified by its ID.
	// the provided ID must be exactly the same, not uses the LIKE select expression.
	GetJournalById(journalId string) (Journal, error)

	// ListJournals retrieve list of journals with transaction date between the `from` and `until` time range inclusive.
	// This function uses pagination.
	ListJournals(from time.Time, until time.Time, request PageRequest) (PageResult, []Journal, error)

	// RenderJournal Render this journal into string for easy inspection
	RenderJournal(journal Journal) string
}

// TransactionManager is interface used for managing transaction data/table
type TransactionManager interface {
	// NewTransaction will create new blank un-persisted Transaction
	NewTransaction() Transaction

	// IsTransactionIdExist will check if an Transaction ID/number is exist in the database.
	IsTransactionIdExist(id string) (bool, error)

	// GetTransactionById will retrieve one single transaction that identified by some ID
	GetTransactionById(id string) (Transaction, error)

	// ListTransactionsWithAccount retrieves list of transactions that belongs to this account
	// that transaction happens between the `from` and `until` time range.
	// This function uses pagination
	ListTransactionsOnAccount(from time.Time, until time.Time, account Account, request PageRequest) (PageResult, []Transaction, error)

	// RenderTransactionsOnAccount Render list of transaction been down on an account in a time span
	RenderTransactionsOnAccount(from time.Time, until time.Time, account Account, request PageRequest) (string, error)
}

// AccountManager interface is used for managing Accounts
type AccountManager interface {
	// NewAccount will create a new blank un-persisted account.
	NewAccount() Account

	// PersistAccount will save the account into database.
	// will throw error if the account already persisted
	PersistAccount(AccountToPersist Account) error

	// UpdateAccount will update the account database to reflect to the provided account information.
	// This update account function will fail if the account ID/number is not existing in the database.
	UpdateAccount(AccountToUpdate Account) error

	// IsAccountIdExist will check if an account ID/number is exist in the database.
	IsAccountIdExist(id string) (bool, error)

	// GetAccountById retrieve an account information by specifying the ID/number
	GetAccountById(id string) (Account, error)

	// ListAccounts list all account in the database.
	// This function uses pagination
	ListAccounts(request PageRequest) (PageResult, []Account, error)

	// ListAccountByCOA returns list of accounts that have the same COA number.
	// This function uses pagination
	ListAccountByCOA(coa string, request PageRequest) (PageResult, []Account, error)

	// FindAccounts returns list of accounts that have their name contains a substring of specified parameter.
	// this search should  be case insensitive.
	FindAccounts(nameLike string, request PageRequest) (PageResult, []Account, error)
}
