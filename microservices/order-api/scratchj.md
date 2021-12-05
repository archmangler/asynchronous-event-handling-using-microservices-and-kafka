



func (h *orderHandlers) publish(w http.ResponseWriter, r *http.Request) {
	
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
		w.Write([]byte(fmt.Sprintf("We don't speak '%s' around these parts ...", ct)))
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

    fmt.Println("debug> publishing: \n '%s'\n",order)
    
}