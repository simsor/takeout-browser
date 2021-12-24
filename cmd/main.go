package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/simsor/takeout-browser/takeout"
	"github.com/simsor/takeout-browser/web"
)

var (
	folder = flag.String("folder", "", "Google Takeout root folder")
	listen = flag.String("listen", "0.0.0.0:8080", "Host:port to listen on")
)

func main() {
	flag.Parse()

	if *folder == "" {
		fmt.Fprintln(os.Stderr, "folder is required")
		flag.Usage()
		os.Exit(3)
	}

	folders, err := takeout.LoadPhotoTakeoutFolders(*folder)
	if err != nil {
		log.Fatalf("LoadPhotoTakeoutFolders: %v", err)
	}
	log.Printf("Loaded %d folders", len(folders))

	server := web.NewServer(folders)

	log.Printf("Serving on %s...", *listen)
	log.Fatalf("ListenAndServe: %v", server.ListenAndServe(*listen))
}
