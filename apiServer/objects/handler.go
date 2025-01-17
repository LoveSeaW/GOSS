package objects

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodPost {
		post(w, r)
		return
	}
	if method == http.MethodPut {
		put(w, r)
		return
	}
	if method == http.MethodGet {
		get(w, r)
		return
	}
	if method == http.MethodDelete {
		del(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
