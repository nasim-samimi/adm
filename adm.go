package adm

import (
	"context"
	"net/http"
	// "gitlab.com/rwxrob/uniq"
)

const defaultSequence = "0"

// configurations passed to the plugin
type Config struct {
	Sequence     string `json:"sequence,omitempty"`
	UpdatePrefix string `json:"updatePrefix,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Sequence:     defaultSequence,
		UpdatePrefix: "/web-app1",
	}
}

// necessary components of a traefik plugin
type MKAdm struct {
	next         http.Handler
	sequence     string
	name         string
	updatePrefix string // add this
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Sequence) == 0 {
		config.Sequence = defaultSequence
	}

	return &MKAdm{
		next:         next,
		sequence:     config.Sequence,
		name:         name,
		updatePrefix: config.UpdatePrefix,
	}, nil

}

func (m *MKAdm) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// if r.URL.Path == m.updatePrefix+"/admission-update" {
	// 	// You can rotate, reset, or log here
	// 	m.sequence, _ = rotateSequence(m.sequence)

	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("Admission sequence updated\n"))
	// 	return
	// }

	// log.Println("Plugin hit! Path:", r.URL.Path)

	// decision := "0"
	// m.sequence, decision = rotateSequence(m.sequence)

	//header injection to the backend service
	w.Header().Set("X-Plugin-Test", "intercepted")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Plugin intercepted request\n"))

	//header injection to the client response
	// w.Header().Add("Sequence", decision)
	// if decision == "0" {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	m.next.ServeHTTP(w, r)
}

// func rotateSequence(s string) (string, string) {
// 	runes := []rune(s)
// 	rotated := append(runes[1:], runes[0])
// 	return string(rotated), string(runes[0])
// }
