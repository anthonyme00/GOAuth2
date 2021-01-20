package util

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestOpenPort(t *testing.T) {
	port := GetRandomOpenPort()

	//try to connect through that port
	for i := 0; i < 10; i++ {
		server := http.Server{
			Addr: net.JoinHostPort("127.0.0.1", fmt.Sprintf("%v", port)),
		}

		server_err := make(chan error)
		go func() {
			server_err <- server.ListenAndServe()
			close(server_err)
		}()

		dur, _ := time.ParseDuration("0.1s")

		time.Sleep(dur)
		server.Close()

		err, _ := <-server_err
		if err != http.ErrServerClosed {
			t.Errorf("Failed ! Can't open listener on that port. Error: %s", err)
			t.FailNow()
		}
	}

	t.Log("Success !")
}
