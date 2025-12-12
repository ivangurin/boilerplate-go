package swagger

import (
	_ "embed"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:embed swagger.json
var swaggerJSON []byte

func Register(mux *http.ServeMux) {
	// Swagger UI
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	// swagger.json
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(swaggerJSON)
	})
}
