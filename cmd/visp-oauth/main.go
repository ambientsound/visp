package main

import (
	"net/http"
	"os"

	"github.com/ambientsound/visp/spotify/proxyserver"
	log "github.com/sirupsen/logrus"
)

// Simple HTTP server that lets users authenticate with Spotify.
// Access tokens are sent back to the client.

func main() {
	clientID := os.Getenv("VISP_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("VISP_OAUTH_CLIENT_SECRET")
	redirectURL := os.Getenv("VISP_OAUTH_REDIRECT_URL")
	listenAddr := os.Getenv("VISP_OAUTH_LISTEN_ADDR")

	renderer := &spotify_proxyserver.JSONRenderer{}
	server := spotify_proxyserver.New(clientID, clientSecret, redirectURL, renderer)

	log.Infof("Visp oauth proxy starting")
	log.Infof("Listening for connections on %s...\n", listenAddr)
	handler := spotify_proxyserver.Router(server)
	err := http.ListenAndServe(listenAddr, handler)

	if err != nil {
		log.Errorf("Fatal error: %s", err)
		os.Exit(1)
	}

	log.Errorf("Visp oauth proxy terminated")
}
