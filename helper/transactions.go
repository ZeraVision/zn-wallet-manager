package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/ZeraVision/zn-wallet-manager/functions"
)

type Transactions struct {
	Found       int64               `json:"found"`
	Transaction []TransactionReturn `json:"transaction"`
}

type TransactionReturn struct {
	BlockHeight           int64         `json:"blockHeight"`
	Symbol                string        `json:"symbol"`
	TokenName             string        `json:"tokenName"`
	ItemName              *string       `json:"itemName,omitempty"`
	TokenType             string        `json:"tokenType"`
	TransactionHash       string        `json:"hash"`
	RelatingHash          *string       `json:"relatingHash,omitempty"`
	Timestamp             int64         `json:"timestamp"`
	FromAddress           string        `json:"from"`
	ToAddress             *string       `json:"to,omitempty"`
	Amount                *Amount       `json:"amount,omitempty"`
	CurType               string        `json:"curType"`
	CurEquiv              *CurEquiv     `json:"curEquiv,omitempty"`
	HighLowVal            *float64      `json:"value"`
	ExpenseRatio          *ExpenseRatio `json:"expenseRatio,omitempty"`
	Memo                  *string       `json:"memo"`
	TransactionType       string        `json:"type"`
	Icon                  *string       `json:"icon"`
	Status                string        `json:"status"`
	Governance            *Governance   `json:"governance,omitempty"`
	ValidatorRegistration *bool         `json:"validatorRegistration,omitempty"`
	Fees                  *FeesStruct   `json:"fees,omitempty"`
}

type Governance struct {
	ProposalHash     *string `json:"hash,omitempty"`
	ProposalTitle    *string `json:"title,omitempty"`
	ProposalSynopsis *string `json:"synopsis,omitempty"`
	ProposalVote     *string `json:"vote,omitempty"`

	// Result
	ResultStatus      *string    `json:"status,omitempty"`
	ResultOptions     *[]string  `json:"options,omitempty"`
	ResultPercentages *[]float64 `json:"votePercent,omitempty"`
}

type ExpenseRatio struct {
	TotalCalled   *big.Int `json:"totalCalled"`
	PartsPerToken *big.Int `json:"partsPerToken"`
	NumWallets    int32    `json:"numWalletsCalled"`
	Rate          *float64 `json:"rate"`
}

type CurEquiv struct {
	LastRate float64 `json:"lastRate"`
	CurRate  float64 `json:"curRate"`

	FeeAuthorized bool `json:"authorized"`

	LastMaxStake float64 `json:"lastMaxStake"`
	CurMaxStake  float64 `json:"curMaxStake"`
}

type Amount struct {
	NonItemized *NonItemizedAmount `json:"nonItemized,omitempty"`
	Itemized    *ItemizedAmount    `json:"itemized,omitempty"`
}

type NonItemizedAmount struct {
	Amount   *big.Int `json:"amount"`
	Parts    *big.Int `json:"parts"`
	CurEquiv float64  `json:"rate"`
}

type ItemizedAmount struct {
	ItemID   *big.Int `json:"id"`
	CurEquiv *float64 `json:"rate"`
}

type FeesStruct struct {
	Base     *BaseFeesStruct     `json:"base,omitempty"`
	Contract *ContractFeesStruct `json:"contract,omitempty"`
}

type BaseFeesStruct struct {
	BaseFees           *big.Int `json:"total"`
	BaseValidator      *big.Int `json:"validator"`
	BaseBurn           *big.Int `json:"burn"`
	BaseOther          *big.Int `json:"treasury"`
	BaseFeesSymbol     *string  `json:"feeSymbol"`
	BaseFeeSymbolParts *big.Int `json:"parts"`
	Rate               *float64 `json:"rate"`
	Type               *string  `json:"type"`
}

type ContractFeesStruct struct {
	TotalFees         *big.Int `json:"total"`
	TknValidator      *big.Int `json:"validator"`
	TknBurn           *big.Int `json:"burn"`
	TknOther          *big.Int `json:"project"`
	TknFeesSymbol     *string  `json:"feeSymbol"`
	TknFeeSymbolParts *big.Int `json:"parts"`
	Rate              *float64 `json:"rate"`
	Type              *string  `json:"type"`
}

// Utilizing getWalletTxnsAdvanced indexer function
// GetTransactions retrieves transactions for a given address (without pagination)
// * this function may not return and parse all data for all transactions correctly. Template for your customization if nessasary. Please ensure you test for your case.
func GetTransactions(
	address string,
	symbol *string,
	itemID *big.Int,
	sortType *string,
	limit *int,
	showFees bool,
	filter *[]string,
	specialInstructions *string,
) (*Transactions, error) {

	baseUrl := fmt.Sprintf("%s/store?requestType=getWalletTxnsAdvanced&address=%s", os.Getenv("INDEXER_URL"), address)

	// Append optional query parameters
	if symbol != nil {
		baseUrl += fmt.Sprintf("&symbol=%s", *symbol)
	}
	if itemID != nil {
		baseUrl += fmt.Sprintf("&itemID=%d", itemID)
	}
	if sortType != nil {
		baseUrl += fmt.Sprintf("&sortType=%s", *sortType)
	}
	if limit != nil {
		baseUrl += fmt.Sprintf("&limit=%d", *limit)
	} else { // default 10 to limit, change as needed
		baseUrl += "&limit=10"
	}

	if showFees {
		baseUrl += "&showFees=true"
	}
	if filter != nil {
		baseUrl += fmt.Sprintf("&filter=%s", strings.Join(*filter, ","))
	}
	if specialInstructions != nil {
		baseUrl += fmt.Sprintf("&specialInstructions=%s", *specialInstructions)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Target", "explorer")
	req.Header.Set("Authorization", "Api-Key "+os.Getenv("INDEXER_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed with status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal raw JSON response
	var txnResponse struct {
		Found       int64                    `json:"found"`
		Transaction []map[string]interface{} `json:"transaction"` // Raw transactions
	}

	if err := json.Unmarshal(body, &txnResponse); err != nil {
		return nil, err
	}

	// Convert raw transactions into structured format
	var transactions Transactions
	transactions.Found = txnResponse.Found
	for _, rawTxn := range txnResponse.Transaction {
		transaction, err := parseTransaction(rawTxn)

		if err != nil {
			return nil, err
		}

		transactions.Transaction = append(transactions.Transaction, transaction)
	}

	return &transactions, nil
}

// Safe transaction parsing with error handling and logging
func parseTransaction(rawTxn map[string]interface{}) (TransactionReturn, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] parseTransaction failed: %v, rawTxn: %+v\n", r, rawTxn)
		}
	}()

	// Extract required fields safely
	txnBlockHeight, ok := functions.SafeInt64(rawTxn["blockHeight"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid blockHeight: %v", rawTxn["blockHeight"])
	}
	txnTokenSymbol, ok := functions.SafeString(rawTxn["symbol"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid symbol: %v", rawTxn["symbol"])
	}
	tokenName, ok := functions.SafeString(rawTxn["tokenName"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid tokenName: %v", rawTxn["tokenName"])
	}
	tokenType, ok := functions.SafeString(rawTxn["tokenType"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid tokenType: %v", rawTxn["tokenType"])
	}
	txnTransactionHash, ok := functions.SafeString(rawTxn["hash"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid hash: %v", rawTxn["hash"])
	}
	txnUnixTime, ok := functions.SafeInt64(rawTxn["timestamp"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid timestamp: %v", rawTxn["timestamp"])
	}
	txnFromAddress, ok := functions.SafeString(rawTxn["from"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid from address: %v", rawTxn["from"])
	}
	txnTransactionType, ok := functions.SafeString(rawTxn["type"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid transaction type: %v", rawTxn["type"])
	}
	txnTransactionStatus, ok := functions.SafeString(rawTxn["status"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid transaction status: %v", rawTxn["status"])
	}
	curType, ok := functions.SafeString(rawTxn["curType"])
	if !ok {
		return TransactionReturn{}, fmt.Errorf("invalid curType: %v", rawTxn["curType"])
	}

	// Optional fields
	toAddress, _ := functions.SafeOptionalString(rawTxn["to"])
	relatingHash, _ := functions.SafeOptionalString(rawTxn["relatingHash"])
	icon, _ := functions.SafeOptionalString(rawTxn["icon"])
	memo, _ := functions.SafeOptionalString(rawTxn["memo"])
	highLow, _ := functions.SafeOptionalFloat(rawTxn["value"])

	// Parse Nested Structures
	amount, _ := parseAmount(rawTxn["amount"])
	expenseRatio, _ := parseExpenseRatio(rawTxn["expenseRatio"])
	fees, _ := parseFees(rawTxn["fees"])
	governance, _ := parseGovernance(rawTxn["governance"])
	validatorReg, _ := functions.SafeOptionalBool(rawTxn["validatorRegistration"])
	curEquiv, _ := parseCurEquiv(rawTxn["curEquiv"]) // <--- FIX: Parse Currency Equivalent

	// Build and return the transaction struct
	return TransactionReturn{
		BlockHeight:           txnBlockHeight,
		Symbol:                txnTokenSymbol,
		TokenName:             tokenName,
		TokenType:             tokenType,
		TransactionHash:       txnTransactionHash,
		Timestamp:             txnUnixTime,
		FromAddress:           txnFromAddress,
		ToAddress:             toAddress,
		RelatingHash:          relatingHash,
		Amount:                amount,
		CurType:               curType,
		CurEquiv:              curEquiv, // <--- FIX: Assign properly parsed Currency Equivalent
		HighLowVal:            highLow,
		ExpenseRatio:          expenseRatio,
		Memo:                  memo,
		TransactionType:       txnTransactionType,
		Icon:                  icon,
		Status:                txnTransactionStatus,
		Governance:            governance,
		ValidatorRegistration: validatorReg,
		Fees:                  fees,
	}, nil
}

func parseAmount(data interface{}) (*Amount, error) {
	if data == nil {
		return nil, nil
	}

	amountMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("amount data is not a valid map: %+v", data)
	}

	return &Amount{
		NonItemized: &NonItemizedAmount{
			Amount:   functions.ToBigInt(functions.SafeMapFloatToString(amountMap, "amount")),
			Parts:    functions.ToBigInt(functions.SafeMapFloatToString(amountMap, "parts")),
			CurEquiv: functions.SafeMapFloatDefault(amountMap, "rate", 0),
		},
	}, nil
}

func parseExpenseRatio(data interface{}) (*ExpenseRatio, error) {
	if data == nil {
		return nil, nil
	}

	expRatioMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expenseRatio data is not a valid map: %+v", data)
	}

	return &ExpenseRatio{
		TotalCalled:   functions.ToBigInt(functions.SafeMapFloatToString(expRatioMap, "totalCalled")),
		PartsPerToken: functions.ToBigInt(functions.SafeMapFloatToString(expRatioMap, "partsPerToken")),
		NumWallets:    int32(functions.SafeMapFloatDefault(expRatioMap, "numWalletsCalled", 0)),
		Rate:          functions.SafeMapFloat(expRatioMap, "rate"),
	}, nil
}

func parseFees(data interface{}) (*FeesStruct, error) {
	if data == nil {
		return nil, nil
	}

	feesMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("fees data is not a valid map: %+v", data)
	}

	// Parse base fees
	var baseFees *BaseFeesStruct
	if baseData, exists := feesMap["base"].(map[string]interface{}); exists {
		baseFees = &BaseFeesStruct{
			BaseFees:           functions.ToBigInt(functions.SafeMapFloatToString(baseData, "total")),
			BaseValidator:      functions.ToBigInt(functions.SafeMapFloatToString(baseData, "validator")),
			BaseBurn:           functions.ToBigInt(functions.SafeMapFloatToString(baseData, "burn")),
			BaseOther:          functions.ToBigInt(functions.SafeMapFloatToString(baseData, "treasury")),
			BaseFeesSymbol:     functions.SafeMapString(baseData, "feeSymbol"),
			BaseFeeSymbolParts: functions.ToBigInt(functions.SafeMapFloatToString(baseData, "parts")),
			Rate:               functions.SafeMapFloat(baseData, "rate"),
			Type:               functions.SafeMapString(baseData, "type"),
		}
	}

	// Parse contract fees (optional)
	var contractFees *ContractFeesStruct
	if contractData, exists := feesMap["contract"].(map[string]interface{}); exists {
		contractFees = &ContractFeesStruct{
			TotalFees:         functions.ToBigInt(functions.SafeMapFloatToString(contractData, "total")),
			TknValidator:      functions.ToBigInt(functions.SafeMapFloatToString(contractData, "validator")),
			TknBurn:           functions.ToBigInt(functions.SafeMapFloatToString(contractData, "burn")),
			TknOther:          functions.ToBigInt(functions.SafeMapFloatToString(contractData, "project")),
			TknFeesSymbol:     functions.SafeMapString(contractData, "feeSymbol"),
			TknFeeSymbolParts: functions.ToBigInt(functions.SafeMapFloatToString(contractData, "parts")),
			Rate:              functions.SafeMapFloat(contractData, "rate"),
			Type:              functions.SafeMapString(contractData, "type"),
		}
	}

	return &FeesStruct{
		Base:     baseFees,
		Contract: contractFees,
	}, nil
}

func parseGovernance(data interface{}) (*Governance, error) {
	if data == nil {
		return nil, nil
	}

	govMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("governance data is not a valid map: %+v", data)
	}

	return &Governance{
		ProposalHash:      functions.SafeMapString(govMap, "hash"),
		ProposalTitle:     functions.SafeMapString(govMap, "title"),
		ProposalSynopsis:  functions.SafeMapString(govMap, "synopsis"),
		ProposalVote:      functions.SafeMapString(govMap, "vote"),
		ResultStatus:      functions.SafeMapString(govMap, "status"),
		ResultOptions:     functions.SafeMapStringArray(govMap, "options"),
		ResultPercentages: functions.SafeMapFloatArray(govMap, "votePercent"),
	}, nil
}

func parseCurEquiv(data interface{}) (*CurEquiv, error) {
	if data == nil {
		return nil, nil
	}

	curEquivMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("curEquiv data is not a valid map: %+v", data)
	}

	return &CurEquiv{
		LastRate:      functions.SafeMapFloatDefault(curEquivMap, "lastRate", 0),
		CurRate:       functions.SafeMapFloatDefault(curEquivMap, "curRate", 0),
		FeeAuthorized: functions.SafeMapBoolDefault(curEquivMap, "authorized", false),
		LastMaxStake:  functions.SafeMapFloatDefault(curEquivMap, "lastMaxStake", 0),
		CurMaxStake:   functions.SafeMapFloatDefault(curEquivMap, "curMaxStake", 0),
	}, nil
}
