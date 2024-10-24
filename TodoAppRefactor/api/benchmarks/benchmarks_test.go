package benchmarks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"todoApp/api"
	"todoApp/api/contracts"
	dataService "todoApp/services"
)

var (
	wg               sync.WaitGroup
	requestHandlerwg sync.WaitGroup
	DataService      = dataService.NewDataService()
)

func StartServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /todoapp/item/", api.GetHandler(DataService))
	mux.HandleFunc("POST /todoapp/item/", api.CreateHandler(DataService))
	mux.HandleFunc("PUT /todoapp/item/", api.MarkItemAsCompleteHandler(DataService))
	mux.HandleFunc("DELETE /todoapp/item/", api.DeleteHandler(DataService))
	mux.HandleFunc("/todoapp/items/", api.GetAllHandler(DataService))

	return httptest.NewServer(mux)
}

func BenchmarkRandomApiCalls(b *testing.B) {
	stopCh := make(chan struct{})
	requestHandlerwg.Add(1)
	go api.RequestHandler(DataService, &requestHandlerwg, stopCh)

	server := StartServer()
	defer server.Close()
	fmt.Printf("Server started on: %s\n", server.Listener.Addr().String())

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go CallRandomApiCall(&wg, b)
	}

	wg.Wait()
	close(stopCh)
	requestHandlerwg.Wait()
}

func CallRandomApiCall(wg *sync.WaitGroup, b *testing.B) {
	defer wg.Done()

	request := "/todoapp/item/"
	apiCallId := rand.Intn(50)
	var err error
	var req *http.Request
	httpResponseRec := httptest.NewRecorder()

	switch apiCallId {
	case 0:
		// CREATE API CALL
		newItem := contracts.CreateContract{Name: "Test Item"}
		newItemJson, _ := json.Marshal(newItem)

		req, _ = http.NewRequest(http.MethodPost, request, bytes.NewBuffer(newItemJson))
		handler := api.CreateHandler(DataService)
		handler.ServeHTTP(httpResponseRec, req)
	case 1:
		// READ API CALL
		itemId := rand.Intn(50)
		request = request + strconv.Itoa(itemId)

		req, err = http.NewRequest(http.MethodGet, request, nil)
		handler := api.GetHandler(DataService)
		handler.ServeHTTP(httpResponseRec, req)
	case 2:
		// READ ALL API CALL
		request = "/todoapp/items/"

		req, err = http.NewRequest(http.MethodGet, request, nil)
		handler := api.GetAllHandler(DataService)
		handler.ServeHTTP(httpResponseRec, req)
	case 3:
		// MARK ITEM AS COMPLETE API CALL
		itemId := rand.Intn(50)
		request = request + strconv.Itoa(itemId)

		req, err = http.NewRequest(http.MethodPut, request, nil)
		handler := api.MarkItemAsCompleteHandler(DataService)
		handler.ServeHTTP(httpResponseRec, req)
	case 4:
		// DELETE ITEM API
		itemId := rand.Intn(2)
		request = request + strconv.Itoa(itemId)

		req, err = http.NewRequest(http.MethodDelete, request, nil)
		handler := api.DeleteHandler(DataService)
		handler.ServeHTTP(httpResponseRec, req)
	}

	if err != nil {
		b.Error(err)
		return
	}

	if status := httpResponseRec.Code; status != http.StatusOK && status != http.StatusCreated {
		expectedError := "item at specified index does not exist"
		if strings.TrimSpace(httpResponseRec.Body.String()) != expectedError {
			b.Errorf("handler returned unexpected body. Got: %v Want: %v", strings.TrimSpace(httpResponseRec.Body.String()), expectedError)
		}
	}
}
