package scoring

import (
	"billing/api"
	"billing/api/response"
	"billing/internal/model"
	"billing/internal/util"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type APIRequest struct {
	Method string
	Path   string
	Body   any
	Param  map[string]string
}

const (
	APIGetBill = iota
	APIGetBillStatus
	APICreatedBill
	APIMakePayment
)

var mapAPI = map[int]APIRequest{
	APIGetBill: {
		Method: http.MethodGet,
		Path:   "/bills/:loan_id",
	},
	APIGetBillStatus: {
		Method: http.MethodGet,
		Path:   "/bills/:loan_id/status",
	},
	APICreatedBill: {
		Method: http.MethodPost,
		Path:   "/bills",
	},
	APIMakePayment: {
		Method: http.MethodPost,
		Path:   "/bills/:loan_id/payments",
	},
}

func callAPI(req APIRequest) *httptest.ResponseRecorder {
	// Setup
	e := api.Init()

	jsonBody, _ := json.Marshal(req.Body)

	// Add path parameter
	if req.Param != nil {
		for k, v := range req.Param {
			req.Path = strings.ReplaceAll(req.Path, ":"+k, v)
		}
	}

	httpReq := httptest.NewRequest(req.Method, req.Path, bytes.NewReader(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")

	// Create recorder to capture response
	rec := httptest.NewRecorder()

	// Serve the request
	e.ServeHTTP(rec, httpReq)

	return rec
}

func unmarshalResponse[T any](rec *httptest.ResponseRecorder) (T, error) {
	var resp T
	var respAPI response.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &respAPI); err != nil {
		return resp, err
	}
	if respAPI.Data != nil {
		jsonData, _ := json.Marshal(respAPI.Data)
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func seedData() (model.LoanWithBills, error) {
	reqBody := model.Loan{
		ID:           fmt.Sprintf("loan-%d", randomNumber()),
		CustomerID:   fmt.Sprintf("cust-%d", randomNumber()),
		Name:         "Test Loan",
		Period:       50,
		Amount:       5000000,
		InterestRate: 10,
	}

	req := mapAPI[APICreatedBill]
	req.Body = reqBody

	rec := callAPI(req)

	respData, err := unmarshalResponse[model.LoanWithBills](rec)
	if err != nil {
		return respData, err
	}

	return respData, nil
}

func randomNumber() int {
	return rand.Intn(1000000)
}

func addTimeNow(week int) func() {
	testTime := time.Now().Truncate(24*time.Hour).AddDate(0, 0, week*7)
	oldTimeNow := util.TimeNow
	util.TimeNow = func() time.Time {
		return testTime
	}
	return func() {
		util.TimeNow = oldTimeNow
	}
}

// assertDateEqual compares only the date part of two times, ignoring hours/minutes/seconds
func assertDateEqual(t *testing.T, expected, actual time.Time, msgAndArgs ...interface{}) {
	expectedDate := expected.Truncate(24 * time.Hour)
	actualDate := actual.Truncate(24 * time.Hour)
	assert.Equal(t, expectedDate, actualDate, msgAndArgs...)
}
