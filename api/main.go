package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// ตั้งค่า CORS เพื่อให้ Angular (พอร์ต 4200) คุยกับ Go (พอร์ต 8080) ได้
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Content-Type", "application/json")

	response := Message{Text: "สวัสดีจาก Go Backend! 🐹"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	fmt.Println("Backend server started at :8080")
	http.ListenAndServe(":8080", nil)
}