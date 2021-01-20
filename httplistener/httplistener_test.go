package httplistener

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestURLGeneration(t *testing.T) {
	randPort := rand.Intn(65534-8000) + 8000
	listener := HTTPListener{
		port: randPort,
	}
	expectedOutput := fmt.Sprintf("http://127.0.0.1:%v/", randPort)

	if listener.GetUrl() != expectedOutput {
		t.Errorf("Failed ! Listener doesn't create correct url, expected %s got %s", expectedOutput, listener.GetUrl())
		t.FailNow()
	}
}

func TestListener(t *testing.T) {
	listener := HTTPListener{}

	handle := listener.OpenListener(func(w http.ResponseWriter, req *http.Request) {})
	url := listener.GetUrl()

	var req *http.Request
	go func() {
		req = listener.GetResponse(handle)
	}()

	http.Post(url, "text", strings.NewReader("abc"))

	dur, _ := time.ParseDuration("0.5s")
	time.Sleep(dur)

	if req == nil {
		t.Errorf("Failed! Listener failed to get a response")
		t.FailNow()
	}

	t.Log("Success ! ")
}
