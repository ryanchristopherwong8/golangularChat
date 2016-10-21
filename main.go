// package declaration, main is a special package from go
package main
// import statements
import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

// main entry point to program
func main() {
	// command line flags
	port := flag.Int("port", 8000, "port to serve on")
	dir := flag.String("directory", "web/", "directory of web files")
	flag.Parse()

	// handle all requests by serving a file of the same name
	fs := http.Dir(*dir)
	fileHandler := http.FileServer(fs)
	http.Handle("/", fileHandler)
	//handle func (takes in a function) , rather than a handler object
	http.HandleFunc("/ws", wsHandler)

	log.Printf("Running on port %d\n", *port)

	addr := fmt.Sprintf("127.0.0.1:%d", *port)
	// this call blocks -- the progam runs here forever
	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
