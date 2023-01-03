package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/mohnish/dashboard/component"
	"github.com/mohnish/dashboard/hub"
	"github.com/mohnish/dashboard/subs"
)

const (
	StorageDir     string = "storage"
	PluginsDir     string = "plugins"
	StorageDirPerm        = 0777
	PluginsDirPerm        = 0777
)

var (
	count   int = 0
	dashHub     = hub.New()
	ticker      = time.NewTicker(time.Second * 10)
	port        = os.Getenv("PORT")
)

func main() {
	// Handle landing page
	// This function will
	// 1. render a template and load react
	// 2. respond with real time data to the front-end
	go func() {
		defer ticker.Stop()

		for {
			select {
			case client := <-dashHub.Register:
				dashHub.Add(client)
			case client := <-dashHub.Unregister:
				dashHub.Remove(client)
			case msg := <-dashHub.Broadcast:
				dashHub.Publish(msg)
			case <-dashHub.Shutdown:
				dashHub.Close()
			case <-ticker.C:
				log.Println("pruned connections", dashHub.Prune())
				log.Println("connected clients", len(dashHub.Clients))
				log.Println("active go routines", runtime.NumGoroutine())
				PrintMemUsage()
			}
		}
	}()

	go func() {
		// Spin up the plugins
		// 1. read all the json files from /plugins dir
		// 2. create subscriptions
		// 3. start the subscriptions
		// 4. listen to the subscription streams
		// 5. push data down to clients every time we get the response
		files, err := ioutil.ReadDir(PluginsDir)

		if err != nil {
			if os.IsNotExist(err) {
				log.Println("Plugins dir not found. creating one...")

				err := os.Mkdir(PluginsDir, PluginsDirPerm)

				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}

		var pluginConfigs []subs.PluginConfig

		// FIXME: (MT) this needs to be done concurrently in multiple goroutines
		for _, file := range files {
			var pluginConfig subs.PluginConfig

			// TODO: (MT) Store some sort of invalid files list
			// if file.Name() == ".DS_Store" {
			//		continue
			// }

			log.Println("reading file", file.Name())
			body, _ := ioutil.ReadFile(PluginsDir + "/" + file.Name())
			err := json.Unmarshal(body, &pluginConfig)

			if err != nil {
				log.Fatal(err)
			} else {
				pluginConfigs = append(pluginConfigs, pluginConfig)
			}
		}

		log.Println("found configs", pluginConfigs)
		var subscriptions []subs.Subscription

		for _, pluginConfig := range pluginConfigs {
			interval, err := strconv.Atoi(pluginConfig.Interval)

			if err != nil {
				log.Fatal(err)
			}

			log.Println("subscribing to", pluginConfig.Url, interval)
			f := subs.Fetch(pluginConfig.Url, time.Duration(interval)*time.Second)
			s := subs.Subscribe(f) // starts the sub
			subscriptions = append(subscriptions, s)
		}

		megaSub := subs.Merge(subscriptions)

		for update := range megaSub.Updates() {
			log.Println("fetched", update)

			if err != nil {
				log.Fatal(err)
			} else {
				go func(up interface{}) {
					dashHub.Broadcast <- up
				}(update)
			}
		}
	}()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Requested URL:", r.Method+" "+r.URL.Path)

		w.Write([]byte("ok"))
	})

	http.HandleFunc("/static/css/app.css", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Requested URL:", r.Method+" "+r.URL.Path)

		contents, _ := ioutil.ReadFile("static/css/app.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write(contents)
	})

	http.HandleFunc("/static/js/app.js", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Requested URL:", r.Method+" "+r.URL.Path)

		contents, _ := ioutil.ReadFile("static/js/app.js")
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(contents)
	})

	// casting socket func to websocket handler
	http.Handle("/ws", websocket.Handler(socket))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Requested URL:", r.Method+" "+r.URL.Path)

		files, err := ioutil.ReadDir(StorageDir)

		if err != nil {
			if os.IsNotExist(err) {
				log.Println("Storage dir not found. creating one...")

				err := os.Mkdir(StorageDir, StorageDirPerm)

				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
		}

		var comps []component.Component

		// FIXME: (MT) this needs to be done concurrently in multiple goroutines
		for _, file := range files {
			var comp component.Component

			body, err := ioutil.ReadFile(StorageDir + "/" + file.Name())

			if err != nil {
				log.Fatal("failed to read file", file.Name(), err)
			}

			err = json.Unmarshal(body, &comp)

			log.Println("~~~", comp)

			if err != nil {
				log.Fatal(err)
			} else {
				comps = append(comps, comp)
			}
		}

		t, _ := template.ParseFiles("static/index.html")
		t.Execute(w, comps)
	})

	// Handle push notifications
	http.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Requested URL:", r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}

		var comp component.Component

		// FIXME: (MT) This is a stream. This needs to be switched to use
		// dec.More() with a for loop since the incoming stream `r.Body`
		// might be bigger and so it'd take more than once decode call
		json.NewDecoder(r.Body).Decode(&comp)

		log.Println("Request body", comp.App)
		comp.Save()

		go func() {
			dashHub.Broadcast <- comp
		}()

		json.NewEncoder(w).Encode(comp)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Helper functions
// TOOD: (MT) create a struct and send these to the front end
// Use these values in the UI
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	log.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	log.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	log.Printf("\tSys = %v MiB", bToMb(m.Sys))
	log.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func socket(client *websocket.Conn) {
	count++

	var goRoutineId = count
	heartBeat := time.NewTicker(10 * time.Second)
	defer heartBeat.Stop()
	log.Println("socket connection established")

	go func() {
		dashHub.Register <- client
	}()

	// Heartbeat implementation to keep the conn alive
	for {
		select {
		case <-heartBeat.C:
			log.Println("tick", goRoutineId)
			err := websocket.JSON.Send(client, ".")

			if err != nil {
				dashHub.Unregister <- client
				return
			}
		}
	}
}
