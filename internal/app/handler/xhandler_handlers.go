package handler

import "net/http"

func (x *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	x.Error(w, "Not Found", http.StatusNotFound)
}
