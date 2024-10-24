package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"todoApp/api/contracts"
	"todoApp/api/responses"
	dataService "todoApp/services"
)

type CreateCommand struct {
	Item contracts.CreateContract
	Resp chan responses.CreateRes
}

type GetCommand struct {
	Id   int
	Resp chan responses.GetRes
}

type GetAllCommand struct {
	Resp chan responses.GetAllRes
}

type MarkAsCompleteCommand struct {
	Id   int
	Resp chan responses.MarkAsCompleteRes
}

type DeleteCommand struct {
	Id   int
	Resp chan responses.DeleteRes
}

var (
	createCh         = make(chan CreateCommand)
	getCh            = make(chan GetCommand)
	getAllCh         = make(chan GetAllCommand)
	markAsCompleteCh = make(chan MarkAsCompleteCommand)
	deleteCh         = make(chan DeleteCommand)
)

func RootHanlder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Server Successfully launched")
}

func RequestHandler(dataService dataService.IDataService, wg *sync.WaitGroup, stopCh <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case cmd := <-createCh:
			err := dataService.CreateTodoItem(cmd.Item.Name)
			cmd.Resp <- responses.CreateRes{Error: err}
		case cmd := <-getCh:
			item, err := dataService.GetTodoItem(cmd.Id)
			cmd.Resp <- responses.GetRes{Item: item, Error: err}
		case cmd := <-getAllCh:
			items := dataService.GetAllTodoItems()
			cmd.Resp <- responses.GetAllRes{Items: items}
		case cmd := <-markAsCompleteCh:
			err := dataService.MarkItemAsComplete(cmd.Id)
			cmd.Resp <- responses.MarkAsCompleteRes{Error: err}
		case cmd := <-deleteCh:
			err := dataService.DeleteTodoItem(cmd.Id)
			cmd.Resp <- responses.DeleteRes{Error: err}
		case <-stopCh:
			return
		}
	}
}

func CreateHandler(dataService dataService.IDataService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var todoItemName contracts.CreateContract
		json.NewDecoder(r.Body).Decode(&todoItemName)
		respCh := make(chan responses.CreateRes)
		createCh <- CreateCommand{Item: todoItemName, Resp: respCh}
		resp := <-respCh
		if resp.Error != nil {
			http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contracts.GetContract{Name: todoItemName.Name, Complete: false})
	}
}

func GetHandler(dataService dataService.IDataService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexStr := strings.TrimPrefix(r.URL.Path, "/todoapp/item/")
		if index, convErr := strconv.Atoi(indexStr); convErr == nil {
			respCh := make(chan responses.GetRes)
			getCh <- GetCommand{Id: index, Resp: respCh}
			resp := <-respCh
			if resp.Error != nil {
				http.Error(w, resp.Error.Error(), http.StatusNotFound)
				return
			} else {
				jsonRes := resp.Item
				json.NewEncoder(w).Encode(jsonRes)
			}
		} else {
			http.Error(w, "invalid request parameter type", http.StatusBadRequest)
		}
	}
}

func GetAllHandler(dataService dataService.IDataService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respCh := make(chan responses.GetAllRes)
		getAllCh <- GetAllCommand{Resp: respCh}
		resp := <-respCh

		jsonRes := resp.Items
		json.NewEncoder(w).Encode(jsonRes)
	}
}

func MarkItemAsCompleteHandler(dataService dataService.IDataService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexStr := strings.TrimPrefix(r.URL.Path, "/todoapp/item/")
		if index, convErr := strconv.Atoi(indexStr); convErr == nil {
			respCh := make(chan responses.MarkAsCompleteRes)
			markAsCompleteCh <- MarkAsCompleteCommand{Id: index, Resp: respCh}
			resp := <-respCh

			if resp.Error != nil {
				http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Item Successfully marked as completed")
		} else {
			http.Error(w, "invalid request parameter type", http.StatusBadRequest)
		}
	}
}

func DeleteHandler(dataService dataService.IDataService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexStr := strings.TrimPrefix(r.URL.Path, "/todoapp/item/")
		if index, convErr := strconv.Atoi(indexStr); convErr == nil {
			respCh := make(chan responses.DeleteRes)
			deleteCh <- DeleteCommand{Id: index, Resp: respCh}
			resp := <-respCh
			if resp.Error != nil {
				http.Error(w, resp.Error.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode("Item successfully deleted")
		} else {
			http.Error(w, "invalid request parameter type", http.StatusBadRequest)
		}
	}
}
