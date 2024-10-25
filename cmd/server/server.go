package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var (
	listen   = flag.String("listen", ":8080", "listen address")
	dir      = flag.String("dir", ".", "directory to serve")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan struct{})
)

func main() {
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := initializeWatcher(absDir)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go startHTTPServer(absDir)
	go handleWebSocketConnections()

	watchForChanges(watcher)
}

func initializeWatcher(absDir string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			log.Printf("Watching %q for changes...", path)
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return watcher, nil
}

func startHTTPServer(absDir string) {
	http.HandleFunc("/ws", handleWebSocket)
	http.Handle("/", http.FileServer(http.Dir(absDir)))
	log.Printf("listening on %q...", *listen)
	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			delete(clients, conn)
			break
		}
	}
}

func handleWebSocketConnections() {
	for {
		<-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
			if err != nil {
				log.Println("write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func watchForChanges(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
				recompileAndReload()
				broadcast <- struct{}{}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func recompileAndReload() {
	cmd := exec.Command("go", "build", "-o", "./cmd/markdown/main.wasm", "./cmd/markdown")
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("recompile error: %v\nOutput:\n%s", err, output)
		return
	}
	log.Println("recompiled successfully")
}
