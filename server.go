package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
)

func route(handler PasswordHandler) http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(res,
				http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
			return
		}

		var password Password
		if err := json.NewDecoder(req.Body).Decode(&password); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := handler(&password)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		jsondata, err := json.Marshal(data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(jsondata)
	}
	return http.HandlerFunc(fn)
}

func main() {
	numcpu := runtime.NumCPU()
	numprocs := runtime.GOMAXPROCS(numcpu)

	log.Printf("GOMAXPROCS set to %d, from %d", numcpu, numprocs)

	http.Handle("/hash", route(Hash))
	http.Handle("/compare", route(Compare))

	addr := ":9004"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("Listening at %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
