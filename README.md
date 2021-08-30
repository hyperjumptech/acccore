# ACCCORE - Accounting Core in Go

---
A core accounting library made in golang. Using this library, you quicly
enable Double Entry Book Keeping which usefull to maintain your internal application
point and reward system. Manage GOLD, POINT, DIAMONDS you distribute to your user
have them accountable, controlled and traced.

```text
=== RUN   TestAccounting_CreateNewJournal
    Accounting_test.go:105: Journal Entry : 5274YX3Y65
        Journal Date  : 2021-07-02 10:59:45.0173345 +0700 +07 m=+0.002096401
        Description   : Creditor Topup Gold
        +------------+------------+--------------------+---------+---------+
        |   TRX ID   |  ACCOUNT   |    DESCRIPTION     |  DEBIT  | CREDIT  |
        +------------+------------+--------------------+---------+---------+
        | ZTUMIZQ565 | T5Z05Z0IX6 | Added Gold Reserve | 1000000 |         |
        | 506O592P3Z | 09TI0ZZZ6U | Added Gold Equity  |         | 1000000 |
        +------------+------------+--------------------+---------+---------+
        |                                                1000000 | 1000000 |
        +------------+------------+--------------------+---------+---------+
        
    Accounting_test.go:126: Journal Entry : 25P8ZVZZOZ
        Journal Date  : 2021-07-02 10:59:45.017858 +0700 +07 m=+0.002619901
        Description   : GOLD purchase transaction
        +------------+------------+-------------------+--------+--------+
        |   TRX ID   |  ACCOUNT   |    DESCRIPTION    | DEBIT  | CREDIT |
        +------------+------------+-------------------+--------+--------+
        | 050LTUIIZX | T0PTZZ70KX | Add debitor AR    | 200000 |        |
        | Q9MUW650YN | T5Z05Z0IX6 | Gold Disbursement |        | 200000 |
        +------------+------------+-------------------+--------+--------+
        |                                               200000 | 200000 |
        +------------+------------+-------------------+--------+--------+
        
    Accounting_test.go:137: Account Number    : T5Z05Z0IX6
        Account Name      : Gold Loan
        Description       : Gold base loan reserve
        Currency          : GOLD
        COA               : 1.1
        Transactions From : 2021-07-02 08:59:45.017858 +0700 +07 m=-7199.997380099
                     To   : 2021-07-02 12:59:45.017858 +0700 +07 m=+7200.002619901
        #Transactions     : 2
        Showing page      : 1/1
        +------------+--------------------------------+------------+--------------------+---------+--------+---------+
        |   TRX ID   |              TIME              | JOURNAL ID |    DESCRIPTION     |  DEBIT  | CREDIT | BALANCE |
        +------------+--------------------------------+------------+--------------------+---------+--------+---------+
        | ZTUMIZQ565 | 2021-07-02 10:59:45.0173345    | 5274YX3Y65 | Added Gold Reserve | 1000000 |        | 1000000 |
        |            | +0700 +07 m=+0.002096401       |            |                    |         |        |         |
        | Q9MUW650YN | 2021-07-02 10:59:45.017858     | 25P8ZVZZOZ | Gold Disbursement  |         | 200000 |  800000 |
        |            | +0700 +07 m=+0.002619901       |            |                    |         |        |         |
        +------------+--------------------------------+------------+--------------------+---------+--------+---------+
        
--- PASS: TestAccounting_CreateNewJournal (0.00s)
PASS
```
