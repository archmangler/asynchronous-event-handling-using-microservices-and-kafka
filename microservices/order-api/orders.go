package main

//go get -u github.com/segmentio/kafka-go

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

//produce
var ctx = context.Background()

// the topic and broker address are initialized as constants
const (
	topic          = "order-received-events"
	broker1Address = "localhost:9092"
	broker2Address = "localhost:9092"
	broker3Address = "localhost:9092"
)

/*TESTING:

curl -X POST  http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"Name":"newOrder", "ID":"78912","Time":"223232113111","Data":"new order", "Eventname":"newOrder"}'

curl -s http://localhost:8080/orders| jq -r

*/

type Order struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Time      string `json:"time"`
	Data      string `json:"data"`
	Eventname string `json:"eventname"`
}

type orderHandlers struct {
	sync.Mutex
	store map[string]Order
}

func (h *orderHandlers) orders(w http.ResponseWriter, r *http.Request) {

	fmt.Println("debug>Handling request ...")

	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		//h.post(w, r)
		h.publish(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func (h *orderHandlers) get(w http.ResponseWriter, r *http.Request) {
	orders := make([]Order, len(h.store))

	fmt.Print("debug> getting order")

	h.Lock()
	i := 0
	for _, order := range h.store {
		fmt.Println("debug> ", order)
		orders[i] = order
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *orderHandlers) getRandomOrder(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Add("location", fmt.Sprintf("/orders/%s", target))
	w.WriteHeader(http.StatusFound)
}

func (h *orderHandlers) getOrder(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *orderHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
		return
	}

	var order Order
	err = json.Unmarshal(bodyBytes, &order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	order.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	h.Lock()
	h.store[order.ID] = order
	defer h.Unlock()
}

func (h *orderHandlers) publish(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("debug 1> Publishing event to queue ...")

	bodyBytes, err := ioutil.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//fmt.Println("debug 2> Publishing event to queue ...")

	ct := r.Header.Get("content-type")

	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("We don't speak '%s' around these parts ...", ct)))
		return
	}

	//fmt.Println("debug 3> Publishing event to queue ...")

	var order Order

	err = json.Unmarshal(bodyBytes, &order)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//fmt.Println("debug 4> Publishing event to queue ...")

	order.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	h.Lock()
	h.store[order.ID] = order
	defer h.Unlock()

	fmt.Println("debug 5> Publishing event to queue ...", order)
	produce(bodyBytes, ctx)

}

//Publish the message to kafka
func produce(message []byte, ctx context.Context) {
	// initialize a counter
	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address, broker2Address, broker3Address},
		Topic:   topic,
	})

	//for {
	// each kafka message has a key and value. The key is used
	// to decide which partition (and consequently, which broker)
	// the message gets published on
	err := w.WriteMessages(ctx, kafka.Message{
		Key: []byte(strconv.Itoa(i)),

		// create an arbitrary message payload for the value
		Value: []byte(message),
	})

	if err != nil {
		panic("could not write message " + err.Error())
	}

	// log a confirmation once the message is written
	fmt.Println("wrote:", message)
	i++

	// sleep for a second
	time.Sleep(time.Second)
	//}
}

func newOrderHandlers() *orderHandlers {
	return &orderHandlers{
		store: map[string]Order{},
	}

}

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
		return
	}

	w.Write([]byte("<html><h1>Super secret admin portal</h1></html>"))
}

func main() {
	//admin := newAdminPortal()

	orderHandlers := newOrderHandlers()

	http.HandleFunc("/orders", orderHandlers.orders)

	http.HandleFunc("/orders/", orderHandlers.getOrder)

	//        http.HandleFunc("/admin", admin.handler)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err)
	}
}
