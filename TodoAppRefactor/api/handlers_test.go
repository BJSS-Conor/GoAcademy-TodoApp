package api

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"todoApp/api/contracts"
	apiMocks "todoApp/api/mocks"
	"todoApp/data"
)

var (
	mockDataService = apiMocks.NewMockDataService()
	stopCh          = make(chan struct{})
)

func RequestHandlerSetup() {
	var wg sync.WaitGroup
	wg.Add(1)
	go RequestHandler(mockDataService, &wg, stopCh)
}

func TestCreateHandler_ValidName(t *testing.T) {
	defer close(stopCh)
	request := "todoapp/item/"
	newItem := contracts.CreateContract{Name: "Test Item"}
	newItemJson, _ := json.Marshal(newItem)
	expectedRes, _ := json.Marshal(contracts.GetContract{Name: newItem.Name, Complete: false})
	RequestHandlerSetup()

	req, err := http.NewRequest(http.MethodPost, request, bytes.NewBuffer(newItemJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := CreateHandler(mockDataService)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, http.StatusCreated)
	} else if strings.TrimSpace(rr.Body.String()) != string(expectedRes) {
		t.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(rr.Body.String()), string(expectedRes))
	}
}

func TestCreateHandler_InvalidName(t *testing.T) {
	defer close(stopCh)
	request := "todoapp/item/"
	newItem := contracts.CreateContract{Name: ""}
	newItemJson, _ := json.Marshal(newItem)
	expectedRes := "name cannot be empty"
	RequestHandlerSetup()

	req, err := http.NewRequest(http.MethodPost, request, bytes.NewBuffer(newItemJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := CreateHandler(mockDataService)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, http.StatusInternalServerError)
	} else if strings.TrimSpace(rr.Body.String()) != string(expectedRes) {
		t.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(rr.Body.String()), string(expectedRes))
	}
}

func TestGetHandler_ValidRequest(t *testing.T) {
	defer close(stopCh)
	request := "/todoapp/item/0"
	expectedValue := contracts.GetContract{Name: "MockItem", Complete: false}
	expectedJson, _ := json.Marshal(expectedValue)
	RequestHandlerSetup()

	req, err := http.NewRequest(http.MethodGet, request, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := GetHandler(mockDataService)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, http.StatusOK)
	} else if strings.TrimSpace(rr.Body.String()) != string(expectedJson) {
		t.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(rr.Body.String()), string(expectedJson))
	}
}

func TestGetHandler_InvalidRequest(t *testing.T) {
	defer close(stopCh)
	testcases := []struct {
		testName       string
		request        string
		expectedStatus int
		expectedRes    string
	}{
		{"Testing invalid index", "/todoapp/item/10", 404, "item at specified index does not exist"},
		{"Testing invalid type", "todoapp/item/index", 400, "invalid request parameter type"},
	}
	RequestHandlerSetup()

	for _, test := range testcases {
		t.Run(test.testName, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, test.request, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := GetHandler(mockDataService)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedStatus {
				t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, test.expectedStatus)
			} else if strings.TrimSpace(rr.Body.String()) != string(test.expectedRes) {
				t.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(rr.Body.String()), string(test.expectedRes))
			}
		})
	}
}

func TestGetAllHandler(t *testing.T) {
	defer close(stopCh)
	request := "/todoapp/items/"
	expectedValue := contracts.GetAllContract{TodoItems: []data.TodoItem{
		{Name: "TodoItem1", Complete: true},
		{Name: "TodoItem2", Complete: false},
		{Name: "TodoItem3", Complete: true},
	}}
	expectedJson, _ := json.Marshal(expectedValue.TodoItems)
	RequestHandlerSetup()

	req, err := http.NewRequest(http.MethodGet, request, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := GetAllHandler(mockDataService)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, http.StatusOK)
	} else if strings.TrimSpace(rr.Body.String()) != string(expectedJson) {
		t.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(rr.Body.String()), string(expectedJson))
	}
}

func TestMarkItemAsCompleteHandler_ValidRequest(t *testing.T) {
	defer close(stopCh)
	request := "/todoapp/item/0"
	RequestHandlerSetup()

	req, err := http.NewRequest(http.MethodPut, request, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := MarkItemAsCompleteHandler(mockDataService)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, http.StatusOK)
	}
}

func TestMarkItemAsCompleteHandler_InvalidRequest(t *testing.T) {
	defer close(stopCh)
	testCases := []struct {
		testName       string
		request        string
		expectedStatus int
		expectedRes    string
	}{
		{"Testing invalid index", "/todoapp/item/100", 500, "item at specified index does not exist"},
		{"Testing invalid request type", "/todoapp/item/index", 400, "invalid request parameter type"},
	}
	RequestHandlerSetup()

	for _, test := range testCases {
		req, err := http.NewRequest(http.MethodPut, test.request, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := MarkItemAsCompleteHandler(mockDataService)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.expectedStatus {
			t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, test.expectedStatus)
		} else if strings.TrimSpace(rr.Body.String()) != test.expectedRes {
			t.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(rr.Body.String()), test.expectedRes)
		}
	}
}

func TestDeleteHandler_ValidRequest(t *testing.T) {
	defer close(stopCh)
	request := "/todoapp/item/0"
	RequestHandlerSetup()

	req, err := http.NewRequest(http.MethodDelete, request, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := DeleteHandler(mockDataService)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, http.StatusOK)
	}
}

func TestDeleteHandler_InvalidRequest(t *testing.T) {
	defer close(stopCh)
	testCases := []struct {
		testName       string
		request        string
		expectedStatus int
		expectedRes    string
	}{
		{"Testing invalid index", "/todoapp/item/100", 500, "item at specified index does not exist"},
		{"Testing invalid request type", "/todoapp/item/index", 400, "invalid request parameter type"},
	}
	RequestHandlerSetup()

	for _, test := range testCases {
		req, err := http.NewRequest(http.MethodDelete, test.request, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := DeleteHandler(mockDataService)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != test.expectedStatus {
			t.Errorf("handler returned wrong status code. Got: %v Want: %v", status, test.expectedStatus)
		} else if strings.TrimSpace(rr.Body.String()) != test.expectedRes {
			t.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(rr.Body.String()), test.expectedRes)
		}
	}
}

func TestConcurrentAPICalls(t *testing.T) {
	defer close(stopCh)
	var wg sync.WaitGroup
	concurrentRequests := 100
	RequestHandlerSetup()

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go CallRandomApiCall(&wg, t)
	}

	wg.Wait()
}

func CallRandomApiCall(wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()

	request := "/todoapp/item/"
	apiCallId := rand.Intn(5)
	var err error
	var req *http.Request
	httpResponseRec := httptest.NewRecorder()

	switch apiCallId {
	case 0:
		// CREATE API CALL
		newItem := contracts.CreateContract{Name: "Test Item"}
		newItemJson, _ := json.Marshal(newItem)

		req, err = http.NewRequest(http.MethodPost, request, bytes.NewBuffer(newItemJson))
		handler := CreateHandler(mockDataService)
		handler.ServeHTTP(httpResponseRec, req)
	case 1:
		// READ API CALL
		itemId := rand.Intn(50)
		request = request + strconv.Itoa(itemId)

		req, err = http.NewRequest(http.MethodGet, request, nil)
		handler := GetHandler(mockDataService)
		handler.ServeHTTP(httpResponseRec, req)
	case 2:
		// READ ALL API CALL
		request = "/todoapp/items/"

		req, err = http.NewRequest(http.MethodGet, request, nil)
		handler := GetAllHandler(mockDataService)
		handler.ServeHTTP(httpResponseRec, req)
	case 3:
		// MARK ITEM AS COMPLETE API CALL
		itemId := rand.Intn(50)
		request = request + strconv.Itoa(itemId)

		req, err = http.NewRequest(http.MethodPut, request, nil)
		handler := MarkItemAsCompleteHandler(mockDataService)
		handler.ServeHTTP(httpResponseRec, req)
	case 4:
		// DELETE ITEM API
		itemId := rand.Intn(50)
		request = request + strconv.Itoa(itemId)

		req, err = http.NewRequest(http.MethodDelete, request, nil)
		handler := DeleteHandler(mockDataService)
		handler.ServeHTTP(httpResponseRec, req)
	}

	if err != nil {
		t.Error(err)
		return
	}

	if status := httpResponseRec.Code; status != http.StatusOK && status != http.StatusCreated {
		expectedError := "item at specified index does not exist"
		if strings.TrimSpace(httpResponseRec.Body.String()) != expectedError {
			t.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(httpResponseRec.Body.String()), expectedError)
		}
	}
}
