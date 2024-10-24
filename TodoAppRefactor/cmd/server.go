package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"todoApp/api"
	"todoApp/api/responses"
	dataService "todoApp/services"
)

var (
	wg          sync.WaitGroup
	DataService = dataService.NewDataService()
)

func StartServer() {
	http.Handle("/stylesheets/", http.StripPrefix("/stylesheets/", http.FileServer(http.Dir("cmd/web/stylesheets"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("cmd/web/images"))))

	stopCh := make(chan struct{})
	wg.Add(1)
	go api.RequestHandler(DataService, &wg, stopCh)

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("GET /todoapp/item/", api.GetHandler(DataService))
	http.HandleFunc("POST /todoapp/item/", api.CreateHandler(DataService))
	http.HandleFunc("PUT /todoapp/item/", api.MarkItemAsCompleteHandler(DataService))
	http.HandleFunc("DELETE /todoapp/item/", api.DeleteHandler(DataService))
	http.HandleFunc("/todoapp/items/", api.GetAllHandler(DataService))

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	// Clean-up for when 'ctrl+c' is pressed
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	<-sigs
	fmt.Println("Server shutting down...")
	close(stopCh)

	wg.Wait()
	fmt.Println("Server has shut down")
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("cmd/web/pages/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data, err := RequestTodoItems(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		t.Execute(w, data)
	}
}

func RequestTodoItems() (responses.GetAllRes, error) {
	resp, err := http.Get("http://localhost:8080/todoapp/items/")
	if err != nil {
		return responses.GetAllRes{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return responses.GetAllRes{}, fmt.Errorf("failed to retrieve todo items: %s", resp.Status)
	}

	var todoList responses.GetAllRes
	err = json.NewDecoder(resp.Body).Decode(&todoList.Items)
	if err != nil {
		return responses.GetAllRes{}, err
	}

	return todoList, nil
}
