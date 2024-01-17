package adm

import (
	"context"
	"net/http"
	// "gitlab.com/rwxrob/uniq"
)

const defaultSequence = "0"

// configurations passed to the plugin
type Config struct {
	Sequence string `json:"sequence,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Sequence: defaultSequence,
	}
}

// necessary components of a traefik plugin
type MKAdm struct {
	next     http.Handler
	sequence string
	name     string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Sequence) == 0 {
		config.Sequence = defaultSequence
	}

	return &MKAdm{
		next:     next,
		sequence: config.Sequence,
		name:     name,
	}, nil

}

func (m *MKAdm) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decision := "0"
	m.sequence, decision = rotateSequence(m.sequence)
	//header injection to the backend service
	// r.Header.Set("X-Request-ID", uid)

	//header injection to the client response
	w.Header().Add("Sequence", decision)
	if decision == "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// w.Write([]byte("Hello, World!"))
	m.next.ServeHTTP(w, r)
}

func rotateSequence(s string) (string, string) {
	runes := []rune(s)
	rotated := append(runes[1:], runes[0])
	return string(rotated), string(runes[0])
}
