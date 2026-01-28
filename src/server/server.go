package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	shoppinglists "goprofit/shoppingLists"
	"goprofit/utils"
)

func Start() {
	http.HandleFunc("/api/shopping-lists", getShoppingLists)
	http.HandleFunc("/api/performance", getPerformance)

	// Serve static files
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	fmt.Println("Server starting on :8080...")
	fmt.Println("Pprof available at http://localhost:8080/debug/pprof/")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getShoppingLists(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get top 9 lists
	lists := shoppinglists.GetTopDTO(9)

	json.NewEncoder(w).Encode(lists)
}

func getPerformance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	metrics := utils.GetMetrics()
	json.NewEncoder(w).Encode(metrics)
}
