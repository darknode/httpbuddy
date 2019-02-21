package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var instanceName = ""
var showEnvironment = false
var htmlTemplate = `<!DOCTYPE html><html><head>
<style>body { margin: 0; padding: 4px; } table { border: 1px solid #EEE; border-collapse: collapse; padding: 0; margin: 0; } table tr,td { padding: 4px; margin: 2px; } table tr td.head { font-weight: bold; background-color: #EEE; }</style>
</head>
<body>
	<h3>HttpBuddy: {{ .Name }}</h3>
	<table>
		<tr>
			<td class="head">Request URL:</td>
			<td>{{ .URI }}</td>
		</tr>
		<tr>
			<td class="head">Remote address:</td>
			<td>{{ .RemoteAddress }}</td>
		</tr>
		<tr><td class="head" colspan="2">Headers</td></tr>
		{{ range $key, $value := .Headers }}<tr><td>{{ $key }}</td><td>{{ $value }}</td></tr>
		{{end}}
		{{ if .Environment }}
		<tr><td class="head" colspan="2">Environment</td></tr>
		{{ range $key, $value := .Environment }}<tr><td>{{ $key }}</td><td>{{ $value }}</td></tr>
		{{end}}
		{{end}}
	</table>
</body>
</html>`

func handler(w http.ResponseWriter, r *http.Request) {
	looksLikeBrowser := false
	acceptHeader, ok := r.Header[http.CanonicalHeaderKey("Accept")]
	if ok {

	}
	if acceptHeader != nil {
		for _, accept := range acceptHeader {
			if strings.Contains(accept, "text/html") || strings.Contains(accept, "*/*") {
				looksLikeBrowser = true
				break
			}
		}
	}

	data := struct {
		Name          string
		RemoteAddress string
		URI           string
		Headers       map[string][]string
		Environment   map[string]string
	}{
		Name:          instanceName,
		RemoteAddress: r.RemoteAddr,
		URI:           r.RequestURI,
		Headers:       r.Header,
		Environment:   make(map[string]string),
	}

	if showEnvironment {
		for _, e := range os.Environ() {
			pair := strings.Split(e, "=")
			data.Environment[pair[0]] = pair[1]
		}
	}

	if looksLikeBrowser {
		html, err := template.New("foo").Parse(htmlTemplate)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		html.Execute(w, data)
	} else {
		jsonString, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonString)
	}
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if len(port) <= 0 {
		port = "8080"
	}

	instanceName = os.Getenv("NAME")
	if len(instanceName) <= 0 {
		instanceName = "HttpBuddy"
	}

	flag.BoolVar(&showEnvironment, "env", false, "Adds environment variables into the output")
	flag.Parse()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
