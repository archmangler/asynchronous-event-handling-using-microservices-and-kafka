# Milestone 2
=============

*Create a new Order microservice in Go. Create an HTTP GET endpoint to check the microservice’s health. Test that the service works as expected by running the service and executing a call to the health endpoint. Ensure it responds accordingly.*

```
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

/*
curl -X POST  http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"Name": "Something","Manufacturer":"Toshiba", "ID": 78912,"Inpark":"324552424322","Height":"A new Order"}'
*/

type orderHandlers struct {
	sync.Mutex
	store map[string]Order
}

func (h *orderHandlers) orders(w http.ResponseWriter, r *http.Request) {
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


```

*Create a function that publishes an event to a topic in Kafka. Test that this function works correctly by creating a main program in Go that will use this function to publish an event to the OrderReceived topic. Verify that it was received by using the appropriate Kafka command-line operation.*


```


```

*Testing*

- Create topic:


```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-topics.sh --create --topic order-received-events --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
Created topic order-received-events.
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % 
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % 
```

- check topic status:

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-topics.sh \                                                                                                             
  --describe \
  --topic order-received-events \
  --bootstrap-server localhost:9092
Topic: order-received-events    TopicId: 2OmoeJMgTzyEJphSvh6FmA PartitionCount: 1       ReplicationFactor: 1    Configs: segment.bytes=1073741824
        Topic: order-received-events    Partition: 0    Leader: 0       Replicas: 0     Isr: 0
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 % 
```


- Publish to topic via the microservice code:

*CLI Test*

```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %                                                                                                                                    
bin/kafka-console-producer.sh \
  --topic order-received-events \
   --bootstrap-server localhost:9092
>Be Bop Be Bop
>%                                                                                                                                                                                                                                                                                                                        (base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \
  --topic order-received-events \
  --from-beginning \             
  --bootstrap-server localhost:9092 
Be Bop Be Bop
Processed a total of 1 messages
```

*Testing the microservice*

Send a message:

```
(base) welcome@Traianos-MacBook-Pro order-api % 
(base) welcome@Traianos-MacBook-Pro order-api % curl -X POST  http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"Name":"newOrder", "ID":"78912","Time":"223232113111","Data":"new order", "Eventname":"newOrder"}'

```

*Running the microservice*

```
(base) welcome@Traianos-MacBook-Pro order-api % go build orders.go
(base) welcome@Traianos-MacBook-Pro order-api % ./orders          
debug>Handling request ...
debug 5> Publishing event to queue ... {newOrder 1638694818672239000 223232113111 new order newOrder}
wrote: [123 34 78 97 109 101 34 58 34 110 101 119 79 114 100 101 114 34 44 32 34 73 68 34 58 34 55 56 57 49 50 34 44 34 84 105 109 101 34 58 34 50 50 51 50 51 50 49 49 51 49 49 49 34 44 34 68 97 116 97 34 58 34 110 101 119 32 111 114 100 101 114 34 44 32 34 69 118 101 110 116 110 97 109 101 34 58 34 110 101 119 79 114 100 101 114 34 125]
```


*Read the topic (with cli consumer)*


```
(base) welcome@Traianos-MacBook-Pro kafka_2.13-3.0.0 %  bin/kafka-console-consumer.sh \
  --topic order-received-events \
  --from-beginning \
  --bootstrap-server localhost:9092
Be Bop Be Bop
.
.
\
^[[D{"Name":"newOrder", "ID":"78912","Time":"223232113111","Data":"new order", "Eventname":"newOrder"}


```







