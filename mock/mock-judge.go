package mock

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Submission struct {
	Token string `json:"token"`
}

type Status struct {
	Status struct {
		Id          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

var statuses = map[string]Status{}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/submissions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		token := fmt.Sprintf("mock-%d", rand.Intn(1000000))
		statuses[token] = Status{
			Status: struct {
				Id          int    `json:"id"`
				Description string `json:"description"`
			}{
				Id:          2,
				Description: "Processing",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Submission{Token: token})
	})

	http.HandleFunc("/submissions/", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Path[len("/submissions/"):]
		status, exists := statuses[token]
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if rand.Intn(6)%5 == 0 {
			status.Status.Id = 3
			status.Status.Description = "Accepted"
			status.Stdout = "Hello, World!\n"
			statuses[token] = status
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	})

	fmt.Println("Mock Judge0 server started at :2358")
	http.ListenAndServe(":2358", nil)
}
