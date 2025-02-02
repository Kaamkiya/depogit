package main

import (
	"flag"
	"net"
	"strconv"

	"codeberg.org/Kaamkiya/depogit/internal/server"
)

func main() {
	host := flag.String("host", "0.0.0.0", "The host to run the server on.")
	port := flag.Int("port", 34567, "The port on which to host the server.")
	scanPath := flag.String("path", "/home/zm/projects/", "Where the git repos are stored.")

	server.Start(net.JoinHostPort(*host, strconv.Itoa(*port)), *scanPath)
}
