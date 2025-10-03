package sim

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type submitPayload struct {
	UserID  string `json:"userId"`
	Correct bool   `json:"correct"`
}

func sendOne(baseURL, userId, correctAnswer string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Duration(10+rand.Intn(990)) * time.Millisecond)
	corr := rand.Intn(2) == 0
	body, _ := json.Marshal(submitPayload{UserID: userId, Correct: corr})
	http.Post(baseURL+"/submit", "application/json", bytes.NewBuffer(body))
}

// Run launches N simulated users
func Run(baseURL string, n int, correctAnswer string) {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup
	for i := 1; i <= n; i++ {
		wg.Add(1)
		go sendOne(baseURL, "user_"+strconv.Itoa(i), correctAnswer, &wg)
	}
	wg.Wait()
}
