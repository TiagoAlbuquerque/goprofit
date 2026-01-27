package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	shoppinglists "goprofit/shoppingLists"
)

func Start() {
	http.HandleFunc("/api/shopping-lists", getShoppingLists)

	// Serve static files
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	fmt.Println("Server starting on :8080...")
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
