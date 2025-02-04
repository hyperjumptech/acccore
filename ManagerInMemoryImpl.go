package acccore

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// RECORD and TABLE simulations ***********************

// InMemoryJournalRecords simulates records in Journal table
type InMemoryJournalRecords struct {
	journalID         string
	journalingTime    time.Time
	description       string
	reversal          bool
	reversedJournalID string
	amount            decimal.Decimal
	createTime        time.Time
	createBy          string
}

// InMemoryAccountRecord is simulating records in Account table
type InMemoryAccountRecord struct {
	currency            string
	id                  string
	name                string
	description         string
	baseTransactionType Alignment
	balance             decimal.Decimal
	coa                 string
	createTime          time.Time
	createBy            string
	updateTime          time.Time
	updateBy            string
}

// InMemoryTransactionRecords is simulating records in Transaction table
type InMemoryTransactionRecords struct {
	transactionID   string
	transactionTime time.Time
	accountNumber   string
	journalID       string
	description     string
	transactionType Alignment
	amount          decimal.Decimal
	accountBalance  decimal.Decimal
	createTime      time.Time
	createBy        string
}

// InMemoryCurrencyRecords is the in memory data structure
type InMemoryCurrencyRecords struct {
	code       string
	name       string
	exchange   decimal.Decimal
	createTime time.Time
	createBy   string
	updateTime time.Time
	updateBy   string
}

var (
	// InMemoryJournalTable the simulated Journal table
	InMemoryJournalTable map[string]*InMemoryJournalRecords

	// InMemoryAccountTable the simulated Account table
	InMemoryAccountTable map[string]*InMemoryAccountRecord

	// InMemoryTransactionTable the simulated Transaction table
	InMemoryTransactionTable map[string]*InMemoryTransactionRecords

	// InMemoryCurrencyTable the simulated Currency table
	InMemoryCurrencyTable map[string]*InMemoryCurrencyRecords
)

func init() {
	ClearInMemoryTables()
}

// ClearInMemoryTables initializes the memory tables
func ClearInMemoryTables() {
	InMemoryJournalTable = make(map[string]*InMemoryJournalRecords, 0)
	InMemoryAccountTable = make(map[string]*InMemoryAccountRecord, 0)
	InMemoryTransactionTable = make(map[string]*InMemoryTransactionRecords, 0)
	InMemoryCurrencyTable = make(map[string]*InMemoryCurrencyRecords, 0)
}

// InMemoryJournalManager implementation of JournalManager using inmemory Journal table map
type InMemoryJournalManager struct {
}

// NewJournal will create new blank un-persisted journal
func (jm *InMemoryJournalManager) NewJournal(context context.Context) Journal {
	return &BaseJournal{}
}

// PersistJournal will record a journal entry into database.
// It requires list of Transactions for which each of the transaction MUST BE :
//
//	1.NOT BE PERSISTED. (the journal AccountNumber is not exist in DB yet)
//	2.Pointing or owned by a PERSISTED Account
//	3.Each of this account must belong to the same Currency
//	4.Balanced. The total sum of DEBIT and total sum of CREDIT is equal.
//	5.No duplicate transaction that belongs to the same Account.
//
// If your database support 2 phased commit, you can make all Balance changes in
// accounts and Transactions. If your db do not support this, you can implement your own 2 phase commits mechanism
// on the CommitJournal and CancelJournal
func (jm *InMemoryJournalManager) PersistJournal(context context.Context, journalToPersist Journal) error {
	// First we have to make sure that the journalToPersist is not yet in our database.
	// 1. Checking if the mandatories is not missing
	if journalToPersist == nil {
		return ErrJournalNil
	}
	if len(journalToPersist.GetJournalID()) == 0 {
		logrus.Errorf("error persisting journal. journal is missing the JournalID")
		return ErrJournalMissingID
	}
	if len(journalToPersist.GetTransactions()) == 0 {
		logrus.Errorf("error persisting journal %s. journal contains no Transactions.", journalToPersist.GetJournalID())
		return ErrJournalNoTransaction
	}
	if len(journalToPersist.GetCreateBy()) == 0 {
		logrus.Errorf("error persisting journal %s. journal author not known.", journalToPersist.GetJournalID())
		return ErrJournalMissingAuthor
	}

	// 2. Checking if the journal ID must not in the Database (already persisted)
	//    SQL HINT : SELECT COUNT(*) FROM JOURNAL WHERE JOURNAL.ID = {journalToPersist.GetJournalID()}
	//    If COUNT(*) is > 0 return error
	if _, exist := InMemoryJournalTable[journalToPersist.GetJournalID()]; exist {
		logrus.Errorf("error persisting journal %s. journal already exist.", journalToPersist.GetJournalID())
		return ErrJournalAlreadyPersisted
	}

	// 3. Make sure all journal Transactions are IDed.
	for idx, trx := range journalToPersist.GetTransactions() {
		if len(trx.GetTransactionID()) == 0 {
			logrus.Errorf("error persisting journal %s. transaction %d is missing TransactionID.", journalToPersist.GetJournalID(), idx)
			return ErrJournalTransactionMissingID
		}
	}

	// 4. Make sure all journal Transactions are not persisted.
	for idx, trx := range journalToPersist.GetTransactions() {
		if _, exist := InMemoryTransactionTable[trx.GetTransactionID()]; exist {
			logrus.Errorf("error persisting journal %s. transaction %d is already exist.", journalToPersist.GetJournalID(), idx)
			return ErrJournalTransactionAlreadyPersisted
		}
	}

	// 5. Make sure Transactions are balanced.
	var creditSum, debitSum decimal.Decimal
	for _, trx := range journalToPersist.GetTransactions() {
		if trx.GetAlignment() == DEBIT {
			debitSum = debitSum.Add(trx.GetAmount())
		}
		if trx.GetAlignment() == CREDIT {
			creditSum = creditSum.Add(trx.GetAmount())
		}
	}
	if !creditSum.Equal(debitSum) {
		logrus.Errorf("error persisting journal %s. debit (%d) != credit (%d). journal not Balance", journalToPersist.GetJournalID(), debitSum, creditSum)
		return ErrJournalNotBalance
	}

	// 6. Make sure Transactions account are not appear twice in the journal
	accountDupCheck := make(map[string]bool)
	for _, trx := range journalToPersist.GetTransactions() {
		if _, exist := accountDupCheck[trx.GetAccountNumber()]; exist {
			logrus.Errorf("error persisting journal %s. multiple transaction belong to the same account (%s)", journalToPersist.GetJournalID(), trx.GetAccountNumber())
			return ErrJournalTransactionAccountDuplicate
		}
		accountDupCheck[trx.GetAccountNumber()] = true
	}

	// 7. Make sure Transactions are all belong to existing accounts
	for _, trx := range journalToPersist.GetTransactions() {
		if _, exist := InMemoryAccountTable[trx.GetAccountNumber()]; !exist {
			logrus.Errorf("error persisting journal %s. theres a transaction belong to non existent account (%s)", journalToPersist.GetJournalID(), trx.GetAccountNumber())
			return ErrJournalTransactionAccountNotPersist
		}
	}

	// 8. Make sure Transactions are all have the same Currency
	var currency string
	for idx, trx := range journalToPersist.GetTransactions() {
		// SELECT CURRENCY FROM ACCOUNT WHERE ACCOUNT_NUMBER = {trx.GetAccountNumber()}
		cur := InMemoryAccountTable[trx.GetAccountNumber()].currency
		if idx == 0 {
			currency = cur
		} else {
			if cur != currency {
				logrus.Errorf("error persisting journal %s. Transactions here uses account with different currencies", journalToPersist.GetJournalID())
				return ErrJournalTransactionMixCurrency
			}
		}
	}

	// 9. If this is a Reversal journal, make sure the journal being reversed have not been reversed before.
	if journalToPersist.GetReversedJournal() != nil {
		reversed, err := jm.IsJournalIDReversed(context, journalToPersist.GetJournalID())
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
		journalID:         journalToPersist.GetJournalID(),
		journalingTime:    time.Now(), // now is set
		description:       journalToPersist.GetDescription(),
		reversal:          false,      // will be set
		reversedJournalID: "",         // will be set
		amount:            creditSum,  // since we know credit sum and debit sum is equal, lets use one of the sum.
		createTime:        time.Now(), // now is set
		createBy:          journalToPersist.GetCreateBy(),
	}
	if journalToPersist.GetReversedJournal() != nil {
		journalToInsert.reversedJournalID = journalToPersist.GetReversedJournal().GetJournalID()
		journalToInsert.reversal = true
	}
	// This is when we insert the record into table.
	InMemoryJournalTable[journalToInsert.journalID] = journalToInsert

	// 2 Save the Transactions
	for _, trx := range journalToPersist.GetTransactions() {
		transactionToInsert := &InMemoryTransactionRecords{
			transactionID:   trx.GetTransactionID(),
			transactionTime: time.Now(), // now is set
			accountNumber:   trx.GetAccountNumber(),
			journalID:       journalToInsert.journalID,
			description:     trx.GetDescription(),
			transactionType: trx.GetAlignment(),
			amount:          trx.GetAmount(),
			accountBalance:  decimal.Zero, // will be updated
			createTime:      time.Now(),   // now is set
			createBy:        trx.GetCreateBy(),
		}
		// get the account current Balance
		// SELECT BALANCE, BASE_TRANSACTION_TYPE FROM ACCOUNT WHERE ACCOUNT_ID = {trx.GetAccountNumber()}
		balance, accountTrxType := InMemoryAccountTable[trx.GetAccountNumber()].balance, InMemoryAccountTable[trx.GetAccountNumber()].baseTransactionType

		var newBalance decimal.Decimal
		if transactionToInsert.transactionType == accountTrxType {
			newBalance = balance.Add(transactionToInsert.amount)
		} else {
			newBalance = balance.Sub(transactionToInsert.amount)
		}
		transactionToInsert.accountBalance = newBalance

		// This is when we insert the record into table.
		InMemoryTransactionTable[transactionToInsert.transactionID] = transactionToInsert

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

// IsJournalIDExist will check if a Journal ID/number is exist in the database.
func (jm *InMemoryJournalManager) IsJournalIDExist(context context.Context, id string) (bool, error) {
	// SELECT COUNT(*) FROM JOURNAL WHERE JOURNAL_ID = <AccountNumber>
	// return true if COUNT > 0
	// return false if COUNT == 0
	_, exist := InMemoryJournalTable[id]
	return exist, nil
}

// GetJournalByID retrieved a Journal information identified by its ID.
// the provided ID must be exactly the same, not uses the LIKE select expression.
func (jm *InMemoryJournalManager) GetJournalByID(context context.Context, journalID string) (Journal, error) {
	journalRecord, exist := InMemoryJournalTable[journalID]
	if !exist {
		return nil, ErrJournalIDNotFound
	}
	journal := jm.NewJournal(context).SetDescription(journalRecord.description).SetCreateTime(journalRecord.createTime).
		SetCreateBy(journalRecord.createBy).SetReversal(journalRecord.reversal).
		SetJournalingTime(journalRecord.journalingTime).SetJournalID(journalRecord.journalID).SetAmount(journalRecord.amount)

	if journalRecord.reversal {
		reversed, err := jm.GetJournalByID(context, journalRecord.reversedJournalID)
		if err != nil {
			return nil, ErrJournalLoadReversalInconsistent
		}
		journal.SetReversedJournal(reversed)
	}

	// Populate all Transactions from DB.
	transactions := make([]Transaction, 0)
	// SELECT * FROM TRANSACTION WHERE JOURNAL_ID = {journalRecord.JournalID}
	for _, trx := range InMemoryTransactionTable {
		if trx.journalID == journalRecord.journalID {
			transaction := &BaseTransaction{
				TransactionID:   trx.transactionID,
				TransactionTime: trx.transactionTime,
				AccountNumber:   trx.accountNumber,
				JournalID:       trx.journalID,
				Description:     trx.description,
				TransactionType: trx.transactionType,
				Amount:          trx.amount,
				AccountBalance:  trx.accountBalance,
				CreateTime:      trx.createTime,
				CreateBy:        trx.createBy,
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
		journal, err := jm.GetJournalByID(context, r.journalID)
		if err != nil {
			return PageResult{}, nil, err
		}
		journals[i] = journal
	}
	return pageResult, journals, nil
}

// GetTotalDebit returns sum of all transaction in the DEBIT Alignment
func GetTotalDebit(journal Journal) decimal.Decimal {
	total := decimal.Zero
	for _, t := range journal.GetTransactions() {
		if t.GetAlignment() == DEBIT {
			total = total.Add(t.GetAmount())
		}
	}
	return total
}

// GetTotalCredit returns sum of all transaction in the CREDIT Alignment
func GetTotalCredit(journal Journal) decimal.Decimal {
	total := decimal.Zero
	for _, t := range journal.GetTransactions() {
		if t.GetAlignment() == CREDIT {
			total = total.Add(t.GetAmount())
		}
	}
	return total
}

// IsJournalIDReversed check if the journal with specified ID has been reversed
func (jm *InMemoryJournalManager) IsJournalIDReversed(context context.Context, journalID string) (bool, error) {
	// SELECT COUNT(*) FROM JOURNAL WHERE REVERSED_JOURNAL_ID = {JournalID}
	// return false if COUNT = 0
	// return true if COUNT > 0
	_, exist := InMemoryJournalTable[journalID]
	if exist {
		for _, j := range InMemoryJournalTable {
			if j.reversedJournalID == journalID {
				return true, nil
			}
		}
		return false, nil
	}
	// todo emit error logs just before returning with errors.
	return false, ErrJournalIDNotFound

}

// RenderJournal will render this journal into string for easy inspection
func (jm *InMemoryJournalManager) RenderJournal(context context.Context, journal Journal) string {

	var buff bytes.Buffer
	table := tablewriter.NewWriter(&buff)
	table.SetHeader([]string{"TRX ID", "Account", "Description", "DEBIT", "CREDIT"})
	table.SetFooter([]string{"", "", "", GetTotalDebit(journal).String(), GetTotalCredit(journal).String()})

	for _, t := range journal.GetTransactions() {
		if t.GetAlignment() == DEBIT {
			table.Append([]string{t.GetTransactionID(), t.GetAccountNumber(), t.GetDescription(), t.GetAmount().String(), ""})
		}
	}
	for _, t := range journal.GetTransactions() {
		if t.GetAlignment() == CREDIT {
			table.Append([]string{t.GetTransactionID(), t.GetAccountNumber(), t.GetDescription(), "", t.GetAmount().String()})
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
	exist, err := am.IsAccountIDExist(context, AccountToPersist.GetAccountNumber())
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
		baseTransactionType: AccountToPersist.GetAlignment(),
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
	exist, err := am.IsAccountIDExist(context, AccountToUpdate.GetAccountNumber())
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
		baseTransactionType: AccountToUpdate.GetAlignment(),
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

// IsAccountIDExist will check if an account ID/number is exist in the database.
func (am *InMemoryAccountManager) IsAccountIDExist(context context.Context, id string) (bool, error) {
	// SELECT COUNT(*) FROM ACCOUNT WHERE ACCOUNT_NUMBER = {AccountNumber}
	_, exist := InMemoryAccountTable[id]
	return exist, nil
}

// GetAccountByID retrieve an account information by specifying the ID/number
func (am *InMemoryAccountManager) GetAccountByID(context context.Context, id string) (Account, error) {
	accountRecord, exist := InMemoryAccountTable[id]
	if !exist {
		return nil, ErrAccountIDNotFound
	}
	return &BaseAccount{
		Currency:      accountRecord.currency,
		AccountNumber: accountRecord.id,
		Name:          accountRecord.name,
		Description:   accountRecord.description,
		Alignment:     accountRecord.baseTransactionType,
		Balance:       accountRecord.balance,
		COA:           accountRecord.coa,
		CreateTime:    accountRecord.createTime,
		CreateBy:      accountRecord.createBy,
		UpdateTime:    accountRecord.updateTime,
		UpdateBy:      accountRecord.updateBy,
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
			Currency:      s.currency,
			AccountNumber: s.id,
			Name:          s.name,
			Description:   s.description,
			Alignment:     s.baseTransactionType,
			Balance:       s.balance,
			COA:           s.coa,
			CreateTime:    s.createTime,
			CreateBy:      s.createBy,
			UpdateTime:    s.updateTime,
			UpdateBy:      s.updateBy,
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
			Currency:      s.currency,
			AccountNumber: s.id,
			Name:          s.name,
			Description:   s.description,
			Alignment:     s.baseTransactionType,
			Balance:       s.balance,
			COA:           s.coa,
			CreateTime:    s.createTime,
			CreateBy:      s.createBy,
			UpdateTime:    s.updateTime,
			UpdateBy:      s.updateBy,
		}
		accounts[i] = bacc
	}

	return pageResult, accounts, nil
}

// FindAccounts returns list of accounts that have their Name contains a substring of specified parameter.
// this search should  be case insensitive.
func (am *InMemoryAccountManager) FindAccounts(context context.Context, nameLike string, request PageRequest) (PageResult, []Account, error) {
	resultSlice := make([]*InMemoryAccountRecord, 0)
	lookup := strings.ToUpper(strings.ReplaceAll(nameLike, "%", ""))
	for _, r := range InMemoryAccountTable {
		if strings.Contains(strings.ToUpper(r.name), lookup) {
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
			Currency:      s.currency,
			AccountNumber: s.id,
			Name:          s.name,
			Description:   s.description,
			Alignment:     s.baseTransactionType,
			Balance:       s.balance,
			COA:           s.coa,
			CreateTime:    s.createTime,
			CreateBy:      s.createBy,
			UpdateTime:    s.updateTime,
			UpdateBy:      s.updateBy,
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

// IsTransactionIDExist will check if an Transaction ID/number is exist in the database.
func (tm *InMemoryTransactionManager) IsTransactionIDExist(context context.Context, id string) (bool, error) {
	_, exist := InMemoryTransactionTable[id]
	return exist, nil
}

// GetTransactionByID will retrieve one single transaction that identified by some ID
func (tm *InMemoryTransactionManager) GetTransactionByID(context context.Context, id string) (Transaction, error) {
	trx, exist := InMemoryTransactionTable[id]
	if !exist {
		return nil, ErrTransactionNotFound
	}
	transaction := &BaseTransaction{
		TransactionID:   trx.transactionID,
		TransactionTime: trx.transactionTime,
		AccountNumber:   trx.accountNumber,
		JournalID:       trx.journalID,
		Description:     trx.description,
		TransactionType: trx.transactionType,
		Amount:          trx.amount,
		AccountBalance:  trx.accountBalance,
		CreateTime:      trx.createTime,
		CreateBy:        trx.createBy,
	}

	return transaction, nil
}

// ListTransactionsOnAccount retrieves list of Transactions that belongs to this account
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
			TransactionID:   trx.transactionID,
			TransactionTime: trx.transactionTime,
			AccountNumber:   trx.accountNumber,
			JournalID:       trx.journalID,
			Description:     trx.description,
			TransactionType: trx.transactionType,
			Amount:          trx.amount,
			AccountBalance:  trx.accountBalance,
			CreateTime:      trx.createTime,
			CreateBy:        trx.createBy,
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
		if t.GetAlignment() == DEBIT {
			table.Append([]string{t.GetTransactionID(), t.GetTransactionTime().String(), t.GetJournalID(), t.GetDescription(), "%s", t.GetAmount().String(), "", "%s", t.GetAccountBalance().String()})
		}
		if t.GetAlignment() == CREDIT {
			table.Append([]string{t.GetTransactionID(), t.GetTransactionTime().String(), t.GetJournalID(), t.GetDescription(), "", t.GetAmount().String(), t.GetAccountBalance().String()})
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

// NewInMemoryExchangeManager initializes a new excahnge manager in memory
func NewInMemoryExchangeManager() ExchangeManager {
	return &InMemoryExchangeManager{
		commonDenominator: decimal.NewFromInt(1),
	}
}

// InMemoryExchangeManager is a base implementation of ExchangeManager.
type InMemoryExchangeManager struct {
	commonDenominator decimal.Decimal
}

// IsCurrencyExist will check in the exchange system for a Currency existance
// non-existent Currency means that the Currency is not supported.
// error should be thrown if only there's an underlying error such as db error.
func (em *InMemoryExchangeManager) IsCurrencyExist(context context.Context, currency string) (bool, error) {
	_, exist := InMemoryCurrencyTable[currency]
	return exist, nil
}

// GetDenom get the current common denominator used in the exchange
func (em *InMemoryExchangeManager) GetDenom(context context.Context) decimal.Decimal {
	return em.commonDenominator
}

// SetDenom set the current common denominator value into the specified value
func (em *InMemoryExchangeManager) SetDenom(context context.Context, denom decimal.Decimal) {
	em.commonDenominator = denom
}

// GetCurrency retrieve currency data indicated by the code argument
func (em *InMemoryExchangeManager) GetCurrency(context context.Context, code string) (Currency, error) {
	if curRec, exist := InMemoryCurrencyTable[code]; exist {
		cur := &BaseCurrency{
			Code:       curRec.code,
			Name:       curRec.name,
			Exchange:   curRec.exchange,
			CreateTime: curRec.createTime,
			CreateBy:   curRec.createBy,
			UpdateTime: curRec.updateTime,
			UpdateBy:   curRec.updateBy,
		}
		return cur, nil
	}
	return nil, ErrCurrencyNotFound

}

// CreateCurrency set the specified value as denominator value for that speciffic Currency.
// This function should return error if the Currency specified is not exist.
func (em *InMemoryExchangeManager) CreateCurrency(context context.Context, code, name string, exchange decimal.Decimal, author string) (Currency, error) {
	if _, exist := InMemoryCurrencyTable[code]; exist {
		return nil, ErrCurrencyAlreadyPersisted
	}
	bc := &InMemoryCurrencyRecords{
		code:       code,
		name:       name,
		exchange:   exchange,
		createTime: time.Now(),
		createBy:   author,
		updateTime: time.Now(),
		updateBy:   author,
	}
	InMemoryCurrencyTable[code] = bc
	return &BaseCurrency{
		Code:       code,
		Name:       name,
		Exchange:   exchange,
		CreateTime: time.Now(),
		CreateBy:   author,
		UpdateTime: time.Now(),
		UpdateBy:   author,
	}, nil
}

// UpdateCurrency updates the currency data
// Error should be returned if the specified Currency is not exist.
func (em *InMemoryExchangeManager) UpdateCurrency(context context.Context, code string, currency Currency, author string) error {
	curr, exist := InMemoryCurrencyTable[code]
	if !exist {
		return ErrCurrencyNotFound
	}
	curr.exchange = currency.GetExchange()
	curr.name = currency.GetName()
	curr.exchange = currency.GetExchange()
	curr.updateBy = author
	curr.updateTime = time.Now()

	currency.SetCode(code)
	return nil

}

// CalculateExchangeRate gets the Currency exchange rate for exchanging between the two Currency.
// if any of the Currency is not exist, an error should be returned.
// if from and to Currency is equal, this must return 1.0
func (em *InMemoryExchangeManager) CalculateExchangeRate(context context.Context, fromCurrency, toCurrency string) (decimal.Decimal, error) {
	from, err := em.GetCurrency(context, fromCurrency)
	if err != nil {
		return decimal.Zero, err
	}
	to, err := em.GetCurrency(context, toCurrency)
	if err != nil {
		return decimal.Zero, err
	}
	m1 := em.GetDenom(context).Div(from.GetExchange())
	m2 := m1.Mul(to.GetExchange())
	m3 := m2.Div(em.GetDenom(context))
	//m1 := new(big.Float).Quo(em.GetDenom(context), big.NewFloat(from.GetExchange()))
	//m2 := new(big.Float).Mul(m1, big.NewFloat(to.GetExchange()))
	//m3 := new(big.Float).Quo(m2, em.GetDenom(context))
	return m3, nil
}

// CalculateExchange gets the Currency exchange value for the Amount of fromCurrency into toCurrency.
// If any of the Currency is not exist, an error should be returned.
// if from and to Currency is equal, the returned Amount must be equal to the Amount in the argument.
func (em *InMemoryExchangeManager) CalculateExchange(context context.Context, fromCurrency, toCurrency string, amount decimal.Decimal) (decimal.Decimal, error) {
	exchange, err := em.CalculateExchangeRate(context, fromCurrency, toCurrency)
	if err != nil {
		return decimal.Zero, err
	}
	m1 := exchange.Mul(amount)
	return m1, nil
}

// ListCurrencies will list all currencies.
func (em *InMemoryExchangeManager) ListCurrencies(context context.Context) ([]Currency, error) {
	ret := make([]Currency, 0)
	for _, cur := range InMemoryCurrencyTable {
		rec := &BaseCurrency{
			Code:       cur.code,
			Name:       cur.name,
			Exchange:   cur.exchange,
			CreateTime: cur.createTime,
			CreateBy:   cur.createBy,
			UpdateTime: cur.updateTime,
			UpdateBy:   cur.updateBy,
		}
		ret = append(ret, rec)
	}
	return ret, nil
}
