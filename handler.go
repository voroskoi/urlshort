package urlshort

import (
	"log"
	"net/http"

	"github.com/go-yaml/yaml"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if url, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

// fields must be exported, otherwise yaml package do no unmarshall them
type ptoURL struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYAML(input []byte) ([]ptoURL, error) {
	var parsed []ptoURL
	// XXX: it is possible to Unmarshal to map, try it!
	err := yaml.Unmarshal(input, &parsed)
	if err != nil {
		return nil, err
	}
	log.Println(parsed)
	return parsed, nil
}

func buildMap(parsedYaml []ptoURL) map[string]string {
	pathMap := make(map[string]string, len(parsedYaml))
	for _, val := range parsedYaml {
		pathMap[val.Path] = val.Url
	}
	return pathMap
}
