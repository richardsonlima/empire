package slugs

import (
	"reflect"
	"testing"

	"github.com/remind101/empire/images"
	"github.com/remind101/empire/processes"
)

func TestProcfileExtractor(t *testing.T) {
	image := &images.Image{
		Repo: "ejholmes/docker-statsd",
		ID:   "1234",
	}

	tests := []struct {
		procfile string
		pm       processes.CommandMap
	}{
		{
			`web: ./bin/web`,
			processes.CommandMap{"web": "./bin/web"},
		},
	}

	for _, tt := range tests {
		c := &dockerClient{procfile: tt.procfile}
		e := &ProcfileExtractor{Client: c}

		pm, err := e.Extract(image)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := pm, tt.pm; !reflect.DeepEqual(got, want) {
			t.Fatalf("Extract => %q; %q", got, want)
		}
	}
}

func TestParseProcfile(t *testing.T) {
	tests := []struct {
		in string
		pm processes.CommandMap
	}{
		{
			`web: ./bin/web`,
			processes.CommandMap{
				"web": "./bin/web",
			},
		},

		{
			`
web: ./bin/web
worker: ./bin/worker -port=$PORT -sock=unix:///var/run/docker.sock
`,
			processes.CommandMap{
				"web":    "./bin/web",
				"worker": "./bin/worker -port=$PORT -sock=unix:///var/run/docker.sock",
			},
		},
	}

	for _, tt := range tests {
		pm, err := ParseProcfile([]byte(tt.in))
		if err != nil {
			t.Fatal(err)
		}

		if got, want := pm, tt.pm; !reflect.DeepEqual(got, want) {
			t.Errorf("parseProcfile(%s) => %q; want %q", tt.in, got, want)
		}
	}
}
