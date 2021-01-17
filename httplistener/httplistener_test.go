package httplistener

import (
	"fmt"
	"math/rand"
	"testing"
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
