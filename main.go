package botsbyuberswe

import (
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gorilla/websocket"
	"github.com/syndtr/goleveldb/leveldb"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	cookieName        = "botbyuber"
	clients           = make(map[*websocket.Conn]bool) // connected clients
	broadcast         = make(chan WebsocketMessage)    // broadcast channel
	upgrader          websocket.Upgrader
	db                *leveldb.DB
	clientID          = "***REMOVED***"
	clientSecret      = "***REMOVED***"
	redirectURL       = "https://bots.uberswe.com/callback"
	clientConnections = make(map[string]*twitch.Client)
)

// Define our message object

type Template struct {
	ModifiedHash string
}

type HashRequest struct {
	Hash string
}

func Init() {
	rand.Seed(time.Now().UnixNano())
}

func Run() {

	var err error

	// Log handling
	f, err := os.OpenFile("uberswe.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	// Database handling
	db, err = leveldb.OpenFile("uberswe.db", nil)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Websockets

	initWebsockets()

	// Routes

	routes()

	// Refresh tokens every 10 min
	go refreshHandler()

	// Connect to twitch channels every 5 seconds
	go twitchIRCHandler()

	log.Println("Listening on port 8010")
	err = http.ListenAndServe(":8010", nil)

	if err != nil {
		panic(err)
	}
}
