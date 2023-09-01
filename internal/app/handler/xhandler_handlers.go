// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"net/http"
)

func (x *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	x.Error(w, "Not Found", http.StatusNotFound)
}
