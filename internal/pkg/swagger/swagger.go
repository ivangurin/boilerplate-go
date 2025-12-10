package swagger

import (
	"bufio"
	"io"
	"net/http"
	"os"
)

func Register(mux *http.ServeMux) {
	// swagger-ui
	fs := http.FileServer(http.Dir("pkg/swagger-ui"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	// swagger.json
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, _ *http.Request) {
		file, err := os.Open("pkg/swagger/swagger.json")
		if err != nil {
			http.Error(w, "swagger not found", http.StatusNotFound)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		reader := bufio.NewReader(file)
		_, err = io.Copy(w, reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
