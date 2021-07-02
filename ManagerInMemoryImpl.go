package acccore

import (
	"bytes"
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"sort"
	"strings"
	"time"
)

// RECORD and TABLE simulations.
// InMemoryJournalRecords is simulating records in Journal table
type InMemoryJournalRecords struct {
	journalId         string
	journalingTime    time.Time
	description       string
	reversal          bool
	reversedJournalId string
	amount            int64
	createTime        time.Time
	createBy          string
}

// InMemoryAccountRecord is simulating records in Account table
type InMemoryAccountRecord struct {
	currency            string
	id                  string
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

// InMemoryAccountRecord is simulating records in Transaction table
type InMemoryTransactionRecords struct {
	transactionId   string
	transactionTime time.Time
	accountNumber   string
	journalId       string
	description     string
	transactionType TransactionType
	amount          int64
	accountBalance  int64
	createTime      time.Time
	createBy        string
}

var (
	// InMemoryJournalTable the simulated Journal table
	InMemoryJournalTable map[string]*InMemoryJournalRecords

	// InMemoryAccountTable the simulated Account table
	InMemoryAccountTable map[string]*InMemoryAccountRecord

	// InMemoryTransactionTable the simulated Transaction table
	InMemoryTransactionTable map[string]*InMemoryTransactionRecords
)

func init() {
	ClearInMemoryTables()
}

func ClearInMemoryTables() {
	InMemoryJournalTable = make(map[string]*InMemoryJournalRecords, 0)
	InMemoryAccountTable = make(map[string]*InMemoryAccountRecord, 0)
	InMemoryTransactionTable = make(map[string]*InMemoryTransactionRecords, 0)
}

// InMemoryJournalManager implementation of JournalManager using inmemory Journal table map
type InMemoryJournalManager struct {
}

// NewJournal will create new blank un-persisted journal
func (jm *InMemoryJournalManager) NewJournal(context context.Context) Journal {
	return &BaseJournal{}
}

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
func (jm *InMemoryJournalManager) PersistJournal(context context.Context, journalToPersist Journal) error {
	// First we have to make sure that the journalToPersist is not yet in our database.
	// 1. Checking if the mandatories is not missing
	if journalToPersist == nil {
		return ErrJournalNil
	}
	if len(journalToPersist.GetJournalID()) == 0 {
		logrus.Errorf("error persisting journal. journal is missing the journalID")
		return ErrJournalMissingId
	}
	if len(journalToPersist.GetTransactions()) == 0 {
		logrus.Errorf("error persisting journal %s. journal contains no transactions.", journalToPersist.GetJournalID())
		return ErrJournalNoTransaction
	}
	if len(journalToPersist.GetCreateBy()) == 0 {
		logrus.Errorf("error persisting journal %s. journal author not known.", journalToPersist.GetJournalID())
		return ErrJournalMissingAuthor
	}

	// 2. Checking if the journal ID must not in the Database (already persisted)
	//    SQL HINT : SELECT COUNT(*) FROM JOURNAL WHERE JOURNAL.ID = {journalToPersist.GetJournalID()}
	//    If COUNT(*) is > 0 return error
	if _, exist := InMemoryJournalTable[journalToPersist.GetJournalID()]; exist == true {
		logrus.Errorf("error persisting journal %s. journal already exist.", journalToPersist.GetJournalID())
		return ErrJournalAlreadyPersisted
	}

	// 3. Make sure all journal transactions are IDed.
	for idx, trx := range journalToPersist.GetTransactions() {
		if len(trx.GetTransactionID()) == 0 {
			logrus.Errorf("error persisting journal %s. transaction %d is missing transactionID.", journalToPersist.GetJournalID(), idx)
			return ErrJournalTransactionMissingID
		}
	}

	// 4. Make sure all journal transactions are not persisted.
	for idx, trx := range journalToPersist.GetTransactions() {
		if _, exist := InMemoryTransactionTable[trx.GetTransactionID()]; exist {
			logrus.Errorf("error persisting journal %s. transaction %d is already exist.", journalToPersist.GetJournalID(), idx)
			return ErrJournalTransactionAlreadyPersisted
		}
	}

	// 5. Make sure transactions are balanced.
	var creditSum, debitSum int64
	for _, trx := range journalToPersist.GetTransactions() {
		if trx.GetTransactionType() == DEBIT {
			debitSum += trx.GetAmount()
		}
		if trx.GetTransactionType() == CREDIT {
			creditSum += trx.GetAmount()
		}
	}
	if creditSum != debitSum {
		logrus.Errorf("error persisting journal %s. debit (%d) != credit (%d). journal not balance", journalToPersist.GetJournalID(), debitSum, creditSum)
		return ErrJournalNotBalance
	}

	// 6. Make sure transactions account are not appear twice in the journal
	accountDupCheck := make(map[string]bool)
	for _, trx := range journalToPersist.GetTransactions() {
		if _, exist := accountDupCheck[trx.GetAccountNumber()]; exist {
			logrus.Errorf("error persisting journal %s. multiple transaction belong to the same account (%s)", journalToPersist.GetJournalID(), trx.GetAccountNumber())
			return ErrJournalTransactionAccountDuplicate
		}
		accountDupCheck[trx.GetAccountNumber()] = true
	}

	// 7. Make sure transactions are all belong to existing accounts
	for _, trx := range journalToPersist.GetTransactions() {
		if _, exist := InMemoryAccountTable[trx.GetAccountNumber()]; !exist {
			logrus.Errorf("error persisting journal %s. theres a transaction belong to non existent account (%s)", journalToPersist.GetJournalID(), trx.GetAccountNumber())
			return ErrJournalTransactionAccountNotPersist
		}
	}

	// 8. Make sure transactions are all have the same currency
	var currency string
	for idx, trx := range journalToPersist.GetTransactions() {
		// SELECT CURRENCY FROM ACCOUNT WHERE ACCOUNT_NUMBER = {trx.GetAccountNumber()}
		cur := InMemoryAccountTable[trx.GetAccountNumber()].currency
		if idx == 0 {
			currency = cur
		} else {
			if cur != currency {
				logrus.Errorf("error persisting journal %s. transactions here uses account with different currencies", journalToPersist.GetJournalID())
				return ErrJournalTransactionMixCurrency
			}
		}
	}

	// 9. If this is a reversal journal, make sure the journal being reversed have not been reversed before.
	if journalToPersist.GetReversedJournal() != nil {
		reversed, err := jm.IsJournalIdReversed(context, journalToPersist.GetJournalID())
		if err != nil {
			return err
		}
		if reversed {
			logrus.Errorf("error persisting journal %s. this journal try to make reverse transaction on journals thats already reversed %s", journalToPersist.GetJournalID(), journalToPersist.GetJournalID())
			return ErrJournalCanNotDoubleReverse
		}
	}

	// ALL is OK. So lets start persisting.

	// BEGIN transaction

	// 1. Save the Journal
	journalToInsert := &InMemoryJournalRecords{
		journalId:         journalToPersist.GetJournalID(),
		journalingTime:    time.Now(), // now is set
		description:       journalToPersist.GetDescription(),
		reversal:          false,      // will be set
		reversedJournalId: "",         // will be set
		amount:            creditSum,  // since we know credit sum and debit sum is equal, lets use one of the sum.
		createTime:        time.Now(), // now is set
		createBy:          journalToPersist.GetCreateBy(),
	}
	if journalToPersist.GetReversedJournal() != nil {
		journalToInsert.reversedJournalId = journalToPersist.GetReversedJournal().GetJournalID()
		journalToInsert.reversal = true
	}
	// This is when we insert the record into table.
	InMemoryJournalTable[journalToInsert.journalId] = journalToInsert

	// 2 Save the Transactions
	for _, trx := range journalToPersist.GetTransactions() {
		transactionToInsert := &InMemoryTransactionRecords{
			transactionId:   trx.GetTransactionID(),
			transactionTime: time.Now(), // now is set
			accountNumber:   trx.GetAccountNumber(),
			journalId:       journalToInsert.journalId,
			description:     trx.GetDescription(),
			transactionType: trx.GetTransactionType(),
			amount:          trx.GetAmount(),
			accountBalance:  0,          // will be updated
			createTime:      time.Now(), // now is set
			createBy:        trx.GetCreateBy(),
		}
		// get the account current balance
		// SELECT BALANCE, BASE_TRANSACTION_TYPE FROM ACCOUNT WHERE ACCOUNT_ID = {trx.GetAccountNumber()}
		balance, accountTrxType := InMemoryAccountTable[trx.GetAccountNumber()].balance, InMemoryAccountTable[trx.GetAccountNumber()].baseTransactionType

		newBalance := int64(0)
		if transactionToInsert.transactionType == accountTrxType {
			newBalance = balance + transactionToInsert.amount
		} else {
			newBalance = balance - transactionToInsert.amount
		}
		transactionToInsert.accountBalance = newBalance

		// This is when we insert the record into table.
		InMemoryTransactionTable[transactionToInsert.transactionId] = transactionToInsert

		// Update Account Balance.
		// UPDATE ACCOUNT SET BALANCE = {newBalance},  UPDATEBY = {trx.GetCreateBy()}, UPDATE_TIME = {time.Now()} WHERE ACCOUNT_ID = {trx.GetAccountNumber()}
		InMemoryAccountTable[trx.GetAccountNumber()].balance = newBalance
		InMemoryAccountTable[trx.GetAccountNumber()].updateTime = time.Now()
		InMemoryAccountTable[trx.GetAccountNumber()].updateBy = trx.GetCreateBy()
	}

	// COMMIT transaction

	return nil
}

// CommitJournal will commit the journal into the system
// Only non committed journal can be committed.
// use this if the implementation database do not support 2 phased commit.
// if your database support 2 phased commit, you should do all commit in the PersistJournal function
// and this function should simply return nil.
func (jm *InMemoryJournalManager) CommitJournal(context context.Context, journalToCommit Journal) error {
	return nil
}

// CancelJournal Cancel a journal
// Only non committed journal can be committed.
// use this if the implementation database do not support 2 phased commit.
// if your database do not support 2 phased commit, you should do all roll back in the PersistJournal function
// and this function should simply return nil.
func (jm *InMemoryJournalManager) CancelJournal(context context.Context, journalToCancel Journal) error {
	return nil
}

// IsTransactionIdExist will check if an Transaction ID/number is exist in the database.
func (jm *InMemoryJournalManager) IsJournalIdExist(context context.Context, id string) (bool, error) {
	// SELECT COUNT(*) FROM JOURNAL WHERE JOURNAL_ID = <accountNumber>
	// return true if COUNT > 0
	// return false if COUNT == 0
	_, exist := InMemoryJournalTable[id]
	return exist, nil
}

// GetJournalById retrieved a Journal information identified by its ID.
// the provided ID must be exactly the same, not uses the LIKE select expression.
func (jm *InMemoryJournalManager) GetJournalById(context context.Context, journalId string) (Journal, error) {
	journalRecord, exist := InMemoryJournalTable[journalId]
	if !exist {
		return nil, ErrJournalIdNotFound
	}
	journal := jm.NewJournal(context).SetDescription(journalRecord.description).SetCreateTime(journalRecord.createTime).
		SetCreateBy(journalRecord.createBy).SetReversal(journalRecord.reversal).
		SetJournalingTime(journalRecord.journalingTime).SetJournalID(journalRecord.journalId).SetAmount(journalRecord.amount)

	if journalRecord.reversal == true {
		reversed, err := jm.GetJournalById(context, journalRecord.reversedJournalId)
		if err != nil {
			return nil, ErrJournalLoadReversalInconsistent
		}
		journal.SetReversedJournal(reversed)
	}

	// Populate all transactions from DB.
	transactions := make([]Transaction, 0)
	// SELECT * FROM TRANSACTION WHERE JOURNAL_ID = {journalRecord.journalID}
	for _, trx := range InMemoryTransactionTable {
		if trx.journalId == journalRecord.journalId {
			transaction := &BaseTransaction{
				transactionID:   trx.transactionId,
				transactionTime: trx.transactionTime,
				accountNumber:   trx.accountNumber,
				journalID:       trx.journalId,
				description:     trx.description,
				transactionType: trx.transactionType,
				amount:          trx.amount,
				accountBalance:  trx.accountBalance,
				createTime:      trx.createTime,
				createBy:        trx.createBy,
			}
			transactions = append(transactions, transaction)
		}
	}

	journal.SetTransactions(transactions)

	return journal, nil
}

// ListJournals retrieve list of journals with transaction date between the `from` and `until` time range inclusive.
// This function uses pagination.
func (jm *InMemoryJournalManager) ListJournals(context context.Context, from time.Time, until time.Time, request PageRequest) (PageResult, []Journal, error) {
	// SELECT COUNT(*) FROM JOURNAL WHERE JOURNALING_TIME < {until} AND JOURNALING_TIME > {from}
	allResult := make([]*InMemoryJournalRecords, 0)
	for _, j := range InMemoryJournalTable {
		if j.journalingTime.After(from) && j.journalingTime.Before(until) {
			allResult = append(allResult, j)
		}
	}
	count := len(allResult)
	pageResult := PageResultFor(request, count)

	// SELECT COUNT(*) FROM JOURNAL WHERE JOURNALING_TIME < {until} AND JOURNALING_TIME > {from} ORDER BY JOURNALING TIME LIMIT {pageResult.offset}, {pageResult.pageSize}
	sort.SliceStable(allResult, func(i, j int) bool {
		return allResult[i].journalingTime.Before(allResult[j].journalingTime)
	})

	journals := make([]Journal, pageResult.PageSize)
	for i, r := range allResult[pageResult.Offset : pageResult.Offset+pageResult.PageSize] {
		journal, err := jm.GetJournalById(context, r.journalId)
		if err != nil {
			return PageResult{}, nil, err
		}
		journals[i] = journal
	}
	return pageResult, journals, nil
}

// GetTotalDebit returns sum of all transaction in the DEBIT alignment
func GetTotalDebit(journal Journal) int64 {
	total := int64(0)
	for _, t := range journal.GetTransactions() {
		if t.GetTransactionType() == DEBIT {
			total += t.GetAmount()
		}
	}
	return total
}

// GetTotalCredit returns sum of all transaction in the CREDIT alignment
func GetTotalCredit(journal Journal) int64 {
	total := int64(0)
	for _, t := range journal.GetTransactions() {
		if t.GetTransactionType() == CREDIT {
			total += t.GetAmount()
		}
	}
	return total
}

// IsJournalIdReversed check if the journal with specified ID has been reversed
func (jm *InMemoryJournalManager) IsJournalIdReversed(context context.Context, journalId string) (bool, error) {
	// SELECT COUNT(*) FROM JOURNAL WHERE REVERSED_JOURNAL_ID = {journalID}
	// return false if COUNT = 0
	// return true if COUNT > 0
	_, exist := InMemoryJournalTable[journalId]
	if exist {
		for _, j := range InMemoryJournalTable {
			if j.reversedJournalId == journalId {
				return true, nil
			}
		}
		return false, nil
	} else {
		// todo emit error logs just before returning with errors.
		return false, ErrJournalIdNotFound
	}
}

// Render this journal into string for easy inspection
func (jm *InMemoryJournalManager) RenderJournal(context context.Context, journal Journal) string {
	var buff bytes.Buffer
	table := tablewriter.NewWriter(&buff)
	table.SetHeader([]string{"TRX ID", "Account", "Description", "DEBIT", "CREDIT"})
	table.SetFooter([]string{"", "", "", fmt.Sprintf("%d", GetTotalDebit(journal)), fmt.Sprintf("%d", GetTotalCredit(journal))})

	for _, t := range journal.GetTransactions() {
		if t.GetTransactionType() == DEBIT {
			table.Append([]string{t.GetTransactionID(), t.GetAccountNumber(), t.GetDescription(), fmt.Sprintf("%d", t.GetAmount()), ""})
		}
	}
	for _, t := range journal.GetTransactions() {
		if t.GetTransactionType() == CREDIT {
			table.Append([]string{t.GetTransactionID(), t.GetAccountNumber(), t.GetDescription(), "", fmt.Sprintf("%d", t.GetAmount())})
		}
	}
	buff.WriteString(fmt.Sprintf("Journal Entry : %s\n", journal.GetJournalID()))
	buff.WriteString(fmt.Sprintf("Journal Date  : %s\n", journal.GetJournalingTime().String()))
	buff.WriteString(fmt.Sprintf("Description   : %s\n", journal.GetDescription()))
	table.Render()
	return buff.String()
}

// InMemoryAccountManager implementation of AccountManager using inmemory Account table map
type InMemoryAccountManager struct {
}

// NewAccount will create a new blank un-persisted account.
func (am *InMemoryAccountManager) NewAccount(context context.Context) Account {
	return &BaseAccount{}
}

// PersistAccount will save the account into database.
// will throw error if the account already persisted
func (am *InMemoryAccountManager) PersistAccount(context context.Context, AccountToPersist Account) error {
	if len(AccountToPersist.GetAccountNumber()) == 0 {
		return ErrAccountMissingID
	}
	if len(AccountToPersist.GetName()) == 0 {
		return ErrAccountMissingName
	}
	if len(AccountToPersist.GetDescription()) == 0 {
		return ErrAccountMissingDescription
	}
	if len(AccountToPersist.GetCreateBy()) == 0 {
		return ErrAccountMissingCreator
	}

	// First make sure that The account have never been created in DB.
	exist, err := am.IsAccountIdExist(context, AccountToPersist.GetAccountNumber())
	if err != nil {
		return err
	}
	if exist {
		return ErrAccountAlreadyPersisted
	}

	accountRecord := &InMemoryAccountRecord{
		currency:            AccountToPersist.GetCurrency(),
		id:                  AccountToPersist.GetAccountNumber(),
		name:                AccountToPersist.GetName(),
		description:         AccountToPersist.GetDescription(),
		baseTransactionType: AccountToPersist.GetBaseTransactionType(),
		balance:             AccountToPersist.GetBalance(),
		coa:                 AccountToPersist.GetCOA(),
		createTime:          time.Now(),
		createBy:            AccountToPersist.GetCreateBy(),
		updateTime:          time.Now(),
		updateBy:            AccountToPersist.GetUpdateBy(),
	}

	InMemoryAccountTable[accountRecord.id] = accountRecord

	return nil
}

// UpdateAccount will update the account database to reflect to the provided account information.
// This update account function will fail if the account ID/number is not existing in the database.
func (am *InMemoryAccountManager) UpdateAccount(context context.Context, AccountToUpdate Account) error {
	if len(AccountToUpdate.GetAccountNumber()) == 0 {
		return ErrAccountMissingID
	}
	if len(AccountToUpdate.GetName()) == 0 {
		return ErrAccountMissingName
	}
	if len(AccountToUpdate.GetDescription()) == 0 {
		return ErrAccountMissingDescription
	}
	if len(AccountToUpdate.GetCreateBy()) == 0 {
		return ErrAccountMissingCreator
	}

	// First make sure that The account have never been created in DB.
	exist, err := am.IsAccountIdExist(context, AccountToUpdate.GetAccountNumber())
	if err != nil {
		return err
	}
	if !exist {
		return ErrAccountIsNotPersisted
	}

	accountRecord := &InMemoryAccountRecord{
		currency:            AccountToUpdate.GetCurrency(),
		id:                  AccountToUpdate.GetAccountNumber(),
		name:                AccountToUpdate.GetName(),
		description:         AccountToUpdate.GetDescription(),
		baseTransactionType: AccountToUpdate.GetBaseTransactionType(),
		balance:             AccountToUpdate.GetBalance(),
		coa:                 AccountToUpdate.GetCOA(),
		createTime:          time.Now(),
		createBy:            AccountToUpdate.GetCreateBy(),
		updateTime:          time.Now(),
		updateBy:            AccountToUpdate.GetUpdateBy(),
	}

	InMemoryAccountTable[accountRecord.id] = accountRecord

	return nil
}

// IsAccountIdExist will check if an account ID/number is exist in the database.
func (am *InMemoryAccountManager) IsAccountIdExist(context context.Context, id string) (bool, error) {
	// SELECT COUNT(*) FROM ACCOUNT WHERE ACCOUNT_NUMBER = {accountNumber}
	_, exist := InMemoryAccountTable[id]
	return exist, nil
}

// GetAccountById retrieve an account information by specifying the ID/number
func (am *InMemoryAccountManager) GetAccountById(context context.Context, id string) (Account, error) {
	accountRecord, exist := InMemoryAccountTable[id]
	if !exist {
		return nil, ErrAccountIdNotFound
	}
	return &BaseAccount{
		currency:            accountRecord.currency,
		accountNumber:       accountRecord.id,
		name:                accountRecord.name,
		description:         accountRecord.description,
		baseTransactionType: accountRecord.baseTransactionType,
		balance:             accountRecord.balance,
		coa:                 accountRecord.coa,
		createTime:          accountRecord.createTime,
		createBy:            accountRecord.createBy,
		updateTime:          accountRecord.updateTime,
		updateBy:            accountRecord.updateBy,
	}, nil
}

// ListAccounts list all account in the database.
// This function uses pagination
func (am *InMemoryAccountManager) ListAccounts(context context.Context, request PageRequest) (PageResult, []Account, error) {
	resultSlice := make([]*InMemoryAccountRecord, 0)
	for _, r := range InMemoryAccountTable {
		resultSlice = append(resultSlice, r)
	}
	sort.SliceStable(resultSlice, func(i, j int) bool {
		return resultSlice[i].createTime.Before(resultSlice[j].createTime)
	})

	pageResult := PageResultFor(request, len(resultSlice))
	accounts := make([]Account, pageResult.PageSize)

	for i, s := range resultSlice[pageResult.Offset : pageResult.Offset+pageResult.PageSize] {
		bacc := &BaseAccount{
			currency:            s.currency,
			accountNumber:       s.id,
			name:                s.name,
			description:         s.description,
			baseTransactionType: s.baseTransactionType,
			balance:             s.balance,
			coa:                 s.coa,
			createTime:          s.createTime,
			createBy:            s.createBy,
			updateTime:          s.updateTime,
			updateBy:            s.updateBy,
		}
		accounts[i] = bacc
	}

	return pageResult, accounts, nil
}

// ListAccountByCOA returns list of accounts that have the same COA number.
// This function uses pagination
func (am *InMemoryAccountManager) ListAccountByCOA(context context.Context, coa string, request PageRequest) (PageResult, []Account, error) {
	resultSlice := make([]*InMemoryAccountRecord, 0)
	for _, r := range InMemoryAccountTable {
		if r.coa == coa {
			resultSlice = append(resultSlice, r)
		}
	}
	sort.SliceStable(resultSlice, func(i, j int) bool {
		return resultSlice[i].createTime.Before(resultSlice[j].createTime)
	})

	pageResult := PageResultFor(request, len(resultSlice))
	accounts := make([]Account, pageResult.PageSize)

	for i, s := range resultSlice[pageResult.Offset : pageResult.Offset+pageResult.PageSize] {
		bacc := &BaseAccount{
			currency:            s.currency,
			accountNumber:       s.id,
			name:                s.name,
			description:         s.description,
			baseTransactionType: s.baseTransactionType,
			balance:             s.balance,
			coa:                 s.coa,
			createTime:          s.createTime,
			createBy:            s.createBy,
			updateTime:          s.updateTime,
			updateBy:            s.updateBy,
		}
		accounts[i] = bacc
	}

	return pageResult, accounts, nil
}

// FindAccounts returns list of accounts that have their name contains a substring of specified parameter.
// this search should  be case insensitive.
func (am *InMemoryAccountManager) FindAccounts(context context.Context, nameLike string, request PageRequest) (PageResult, []Account, error) {
	resultSlice := make([]*InMemoryAccountRecord, 0)
	for _, r := range InMemoryAccountTable {
		if strings.Contains(strings.ToUpper(r.name), strings.ToUpper(nameLike)) {
			resultSlice = append(resultSlice, r)
		}
	}
	sort.SliceStable(resultSlice, func(i, j int) bool {
		return resultSlice[i].createTime.Before(resultSlice[j].createTime)
	})

	pageResult := PageResultFor(request, len(resultSlice))
	accounts := make([]Account, pageResult.PageSize)

	for i, s := range resultSlice[pageResult.Offset : pageResult.Offset+pageResult.PageSize] {
		bacc := &BaseAccount{
			currency:            s.currency,
			accountNumber:       s.id,
			name:                s.name,
			description:         s.description,
			baseTransactionType: s.baseTransactionType,
			balance:             s.balance,
			coa:                 s.coa,
			createTime:          s.createTime,
			createBy:            s.createBy,
			updateTime:          s.updateTime,
			updateBy:            s.updateBy,
		}
		accounts[i] = bacc
	}

	return pageResult, accounts, nil
}

// InMemoryTransactionManager implementation of TransactionManager using inmemory Account table map
type InMemoryTransactionManager struct {
}

// NewTransaction will create new blank un-persisted Transaction
func (tm *InMemoryTransactionManager) NewTransaction(context context.Context) Transaction {
	return &BaseTransaction{}
}

// IsTransactionIdExist will check if an Transaction ID/number is exist in the database.
func (tm *InMemoryTransactionManager) IsTransactionIdExist(context context.Context, id string) (bool, error) {
	_, exist := InMemoryTransactionTable[id]
	return exist, nil
}

// GetTransactionById will retrieve one single transaction that identified by some ID
func (tm *InMemoryTransactionManager) GetTransactionById(context context.Context, id string) (Transaction, error) {
	trx, exist := InMemoryTransactionTable[id]
	if !exist {
		return nil, ErrTransactionNotFound
	}
	transaction := &BaseTransaction{
		transactionID:   trx.transactionId,
		transactionTime: trx.transactionTime,
		accountNumber:   trx.accountNumber,
		journalID:       trx.journalId,
		description:     trx.description,
		transactionType: trx.transactionType,
		amount:          trx.amount,
		accountBalance:  trx.accountBalance,
		createTime:      trx.createTime,
		createBy:        trx.createBy,
	}

	return transaction, nil
}

// ListTransactionsWithAccount retrieves list of transactions that belongs to this account
// that transaction happens between the `from` and `until` time range.
// This function uses pagination
func (tm *InMemoryTransactionManager) ListTransactionsOnAccount(context context.Context, from time.Time, until time.Time, account Account, request PageRequest) (PageResult, []Transaction, error) {
	resultRecord := make([]*InMemoryTransactionRecords, 0)
	for _, trx := range InMemoryTransactionTable {
		if trx.accountNumber == account.GetAccountNumber() {
			resultRecord = append(resultRecord, trx)
		}
	}
	sort.SliceStable(resultRecord, func(i, j int) bool {
		return resultRecord[i].createTime.Before(resultRecord[j].createTime)
	})

	pageResult := PageResultFor(request, len(resultRecord))

	transactions := make([]Transaction, len(resultRecord))
	for idx, trx := range resultRecord {
		transaction := &BaseTransaction{
			transactionID:   trx.transactionId,
			transactionTime: trx.transactionTime,
			accountNumber:   trx.accountNumber,
			journalID:       trx.journalId,
			description:     trx.description,
			transactionType: trx.transactionType,
			amount:          trx.amount,
			accountBalance:  trx.accountBalance,
			createTime:      trx.createTime,
			createBy:        trx.createBy,
		}
		transactions[idx] = transaction
	}
	return pageResult, transactions, nil
}

// RenderTransactionsOnAccount Render list of transaction been down on an account in a time span
func (tm *InMemoryTransactionManager) RenderTransactionsOnAccount(context context.Context, from time.Time, until time.Time, account Account, request PageRequest) (string, error) {

	result, transactions, err := tm.ListTransactionsOnAccount(context, from, until, account, request)
	if err != nil {
		return "Error rendering", err
	}

	var buff bytes.Buffer
	table := tablewriter.NewWriter(&buff)
	table.SetHeader([]string{"TRX ID", "TIME", "JOURNAL ID", "Description", "DEBIT", "CREDIT", "BALANCE"})

	for _, t := range transactions {
		if t.GetTransactionType() == DEBIT {
			table.Append([]string{t.GetTransactionID(), t.GetTransactionTime().String(), t.GetJournalID(), t.GetDescription(), fmt.Sprintf("%d", t.GetAmount()), "", fmt.Sprintf("%d", t.GetAccountBalance())})
		}
		if t.GetTransactionType() == CREDIT {
			table.Append([]string{t.GetTransactionID(), t.GetTransactionTime().String(), t.GetJournalID(), t.GetDescription(), "", fmt.Sprintf("%d", t.GetAmount()), fmt.Sprintf("%d", t.GetAccountBalance())})
		}
	}

	buff.WriteString(fmt.Sprintf("Account Number    : %s\n", account.GetAccountNumber()))
	buff.WriteString(fmt.Sprintf("Account Name      : %s\n", account.GetName()))
	buff.WriteString(fmt.Sprintf("Description       : %s\n", account.GetDescription()))
	buff.WriteString(fmt.Sprintf("Currency          : %s\n", account.GetCurrency()))
	buff.WriteString(fmt.Sprintf("COA               : %s\n", account.GetCOA()))
	buff.WriteString(fmt.Sprintf("Transactions From : %s\n", from.String()))
	buff.WriteString(fmt.Sprintf("             To   : %s\n", until.String()))
	buff.WriteString(fmt.Sprintf("#Transactions     : %d\n", result.TotalEntries))
	buff.WriteString(fmt.Sprintf("Showing page      : %d/%d\n", result.Page, result.TotalPages))
	table.Render()
	return buff.String(), err
}
