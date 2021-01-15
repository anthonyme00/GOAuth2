package util

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
)

// Get a random open port.
//
// Source:
// https://stackoverflow.com/questions/43424787/how-to-use-next-available-port-in-http-listenandserve
func GetRandomOpenPort() int {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()
	return port
}

// Open a url in the default browser.
//
// Source:
// https://gist.github.com/nanmu42/4fbaf26c771da58095fa7a9f14f23d27
func OpenInBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}
