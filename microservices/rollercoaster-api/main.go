package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

/*
curl -v http://127.0.0.1:8080/coasters \
  -H 'Content-Type: application/json' \
  -X POST \
  -d '{ "eventName":"newOrder", "productCode":"00000001", "productQuantity":"2", "customerCode":"101010022"}'
*/
type Order struct {
	eventName       string `json:eventName`
	eventTimeStamp  string `json:eventTimeStamp`
	eventID         string `json:eventID`
	productCode     string `json:productCode`
	productQuantity string `json:productQuantity`
	customerCode    string `json:productCode`
}

type ordersHandlers struct {
	sync.Mutex
	store map[string]Order
}

func (h *ordersHandlers) orders(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *ordersHandlers) post(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := ioutil.ReadAll(r.Body)
	fmt.Println("debug> accepting posted data!!")

	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content type application/json but got '%s'", ct)))
	}

	var order Order

	err = json.Unmarshal(bodyBytes, &order)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	order.eventID = fmt.Sprintf("%d", time.Now().UnixNano())

	h.Lock()
	h.store[order.eventID] = order
	defer h.Unlock()
	fmt.Println("debug> accepting posted data!!")
}

func (h *ordersHandlers) getOrder(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.String(), "/")

	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		h.getRandomOrder(w, r)
		return
	}

	h.Lock()

	order, ok := h.store[parts[2]]

	h.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(order)

	if err != nil {
		//TD
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonBytes)
	fmt.Println("debug> accepting GET data!!")

}

func (h *ordersHandlers) get(w http.ResponseWriter, r *http.Request) {

	orders := make([]Order, len(h.store))

	h.Lock()
	i := 0
	for _, order := range h.store {
		orders[i] = order
		i++
	}

	h.Unlock()

	jsonBytes, err := json.Marshal(orders)

	if err != nil {
		//TD
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")

	w.WriteHeader(http.StatusInternalServerError)

	w.Write(jsonBytes)

	fmt.Println("debug> accepting GET data!!")

}

//Admin Panel
type adminPortal struct {
	password string
}

func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}

	return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - unauthorized"))
	}

	w.Write([]byte("<html><h1>Super Secret Admin Portal</h1></html>"))
}

//Constructor function
func newOrderHandlers() *ordersHandlers {
	return &ordersHandlers{

		store: map[string]Order{},
	}
}

func (h *ordersHandlers) getRandomOrder(w http.ResponseWriter, r *http.Request) {

	ids := make([]string, len(h.store))

	h.Lock()

	i := 0

	for id := range h.store {
		ids[i] = id
		i++
	}

	defer h.Unlock()

	var target string

	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return

	} else if len(ids) == 1 {
		target = ids[0]

	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}

	fmt.Println(target)

	w.Header().Add("location", fmt.Sprintf("/orders/%s", target))
	w.WriteHeader(http.StatusFound)
}

func main() {

	orderHandlers := newOrderHandlers()

	admin := newAdminPortal()

	http.HandleFunc("/orders", orderHandlers.orders)

	http.HandleFunc("/orders/", orderHandlers.getOrder)
	http.HandleFunc("/admin", admin.handler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}

}
