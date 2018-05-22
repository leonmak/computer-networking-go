// (i) create a connection socket when contacted by a client (browser);
// (ii) receive the HTTP request from this connection;
// (iii) parse the request to determine the specific file being requested;
// (iv) get the requested file from the server’s file system;
// (v) create an HTTP response message consisting of the requested file preceded by header lines; and
// (vi) send the response over the TCP connection to the requesting browser. If a browser requests a file
// that is not present in your server, your server should return a “404 Not Found” error message.

package main

import (
	"net"
	"log"
	"net/http"
	"bufio"
	"strings"
	"io/ioutil"
	"os"
	"path"
	"io"
	"fmt"
)

func main() {

	port := ":8080"
	addr, _ := net.ResolveTCPAddr("tcp", port)
	listener, _ := net.ListenTCP("tcp4", addr)

	// Using net.http libraries:
	//dir := flag.String("d", ".", "dir")
	//flag.Parse()
	//http.Handle("/", http.FileServer(http.Dir(*dir)))
	//http.ListenAndServe(port, nil)

	// or (see below):
	// handleWithNetHttp()
	for {
		if conn, err := listener.Accept(); err == nil {  // (i)
			go handleClient(conn) // concurrent
		} else {
			log.Fatal(err)
		}
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// (ii) net.Conn satisfies io.Reader interface
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Fatal(err)
	}

	// (iii) can also actually parse the []bytes in conn to get Path
	filePath := strings.TrimPrefix(req.URL.Path, "/")
	log.Println(filePath)

	// (iv)
	var responseStrLines []string
	var responseDat []byte
	statusLine := "HTTP-Version = HTTP/1.1"	// status line
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		// (vi)
		//log.Fatal(err)
		responseStrLines = []string {
			statusLine,
			"HTTP/1.1 404 Not Found",
		}
		responseDat = []byte("<html><body>404 Not Found :(</body></html>")
	} else {
		// (v)
		responseStrLines = []string {
			statusLine,
			"HTTP/1.1 200 OK",			// headers
		}
		if dat, err := ioutil.ReadFile(filePath); err == nil {
			responseDat = dat
		}
	}

	// optional message body (with 1 line above)
	responseStr := strings.Join(responseStrLines, "\n") + "\n\n"

	responseBytes := append([]byte(responseStr), responseDat...)
	fmt.Println(string(responseBytes))
	_, err = conn.Write(responseBytes)
	if err != nil {
		log.Fatal(err)
	}
}

//////////////////////////////////////////////////////////////////
// Using net.http												//
//////////////////////////////////////////////////////////////////

type SimpleHandler interface { // http.Handler interface
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type FileMount struct {
	root string
}

func (fm *FileMount) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	f, err := os.Open(path.Join(fm.root, urlPath))
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(w, f)
}

func FileServer(root string) SimpleHandler {
	return &FileMount{root} // implements
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func handleWithNetHttp() {
	http.Handle("/", FileServer("."))    // Use struct which implements handler ServeHTTP
	http.HandleFunc("/static/", FileHandler)  // Use handler with same signature as ServeHTTP
	http.ListenAndServe(":8080", nil)
}

//////////////////////////////////////////////////////////////////
// END Using net.http											//
//////////////////////////////////////////////////////////////////
