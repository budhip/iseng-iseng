# BILLING ENGINE #
## Loan Billing System ##

A billing engine service for managing loan schedules, payments, and delinquency tracking. The system provides loan billing schedules, tracks outstanding amounts, and identifies delinquent borrowers.

### Problem Statement: ###

We are building a billing system for our Loan Engine that provides:
* **Loan Schedule**: When and how much customers need to pay for their loans
* **Outstanding Amount**: Current pending amount for any given loan
* **Delinquency Status**: Whether a customer has missed 2+ consecutive payments

### Loan Terms: ###
* **Loan Duration**: 50 weeks
* **Loan Amount**: Rp 5,000,000
* **Interest Rate**: 10% per annum (flat rate)
* **Weekly Payment**: Rp 110,000 (equal installments)
* **Payment Rule**: Borrowers can only pay the exact weekly amount or not pay at all
* **Status**: Always ACTIVE

### Data Models: ###

**Loan ID Format**: String identifier (e.g., "100", "ABC123")

### API Endpoints: ###

* **POST /bills** - Create loan with billing schedule
  * Input: Loan details (customer_id, amount, period, interest_rate)
  * Output: Loan with generated weekly bills
  
* **GET /bills/:loan_id** - Get loan billing schedule
  * Input: loan_id (string parameter)
  * Output: Loan details with all bills --Outstanding should equal sum of unpaid bills' TotalAmount

* **POST /bills/status** - Check delinquency status
  * Input: Billing request with loan_id
  * Output: Delinquency status (true if 2+ consecutive missed payments)
  
* **POST /bills/:loan_id/payment** - Make loan payment
  * Input: loan_id (string parameter) + payment details
  * Output: Payment confirmation

### Business Logic: ###

**Outstanding Balance**: 
* Starts at total loan amount (principal + interest)
* Decreases with each payment
* Should reach 0 when loan is fully paid

**Delinquency Rules**:
* Customer becomes delinquent after missing 2 consecutive weekly payments
* Must pay exact weekly amounts to catch up on missed payments
* System tracks delinquent_at date (date of second consecutive missed payment)

**Payment Schedule Example** (50-week loan):
```
Week 1:  Rp 110,000
Week 2:  Rp 110,000  
Week 3:  Rp 110,000
...
Week 50: Rp 110,000
```

# GETTING STARTED #
## GUIDE ##
* If you need to initialize the database, you can use init function in `pkg/db` and create your DDL / DML scripts in `db/*.sql` (we are using dbmate)
* Don't change the request & response body, it will cause the test to fail
* Don't change the API endpoint, it will cause the test to fail

## TODO ##
* Complete the API handler implementation
* Implement bill generation logic in CreateBills
* Implement delinquency checking in GetBillStatus
* Create service/usecase layer for business logic
* Add payment processing in MakePayment endpoint