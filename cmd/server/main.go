package main

import (
	"log"
	"net/http"

	"github.com/justine-george/nexus-decentralized-messaging/internal/directory"
	"github.com/justine-george/nexus-decentralized-messaging/internal/websocket"
)

func main() {
	// Initialize directory service
	dirService := directory.NewService()

	// Start directory service gRPC server
	go dirService.Serve()

	// Set up HTTP server
	http.HandleFunc("/ws", websocket.HandleConnection(dirService))
	http.HandleFunc("/", serveHome)

	// Start HTTP server
	log.Println("Starting server on :8443")
	err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "web/templates/index.html")
}
