// package declaration, main is a special package from go
package main
// import statements
import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"os"
)
var connections map[*websocket.Conn]bool

func sendAll(msg []byte) {
	for conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			// delete connection from map connections
			delete(connections, conn)
			return
		}
	}

}

// writer => anything we write to w is returned to the client
// request => client's request
func wsHandler(w http.ResponseWriter, r * http.	Request){
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close();
	connections[conn] = true
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			delete(connections, conn)
			return
		}
		log.Println(string(msg))
		sendAll(msg)
	}

}

// main entry point to program
func main() {
	// command line flags
	//port := flag.Int("port", 8000, "port to serve on")
	port := os.Getenv("PORT")
    /*
    if port == "" {
        log.WithField("PORT", port).Fatal("$PORT must be set")
    }
    */
	dir := flag.String("directory", "./public/web/", "directory of web files")
	flag.Parse()
	connections = make(map[*websocket.Conn]bool)
	//http.Handle("/", http.FileServer(http.Dir("./public")))

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	//handle func (takes in a function) , rather than a handler object
	http.HandleFunc("/ws", wsHandler)

	log.Printf("Running on port %d\n", port)

	addr := fmt.Sprintf(":%d", &port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
