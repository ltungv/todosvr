package middleware

import "net/http"

// Middleware tạo một handler mới dựa trên handler được truyền vào.
type Middleware func(http.Handler) http.Handler
