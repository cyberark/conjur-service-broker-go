// Package main provides sample app for E2E tests it exposes a HTML page with env variables
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

var tmpl = must(template.New("page").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>sample-app</title>
	</head>
	<body>
		{{range .}}<div>{{ . }}</div>{{else}}<div><strong>no data</strong></div>{{end}}
	</body>
</html>`))

func must(tmpl *template.Template, err error) *template.Template {
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("starting server on port %s", port)
	http.HandleFunc("/", handler)
	srv := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
	var data []string
	envs := []string{"ORG_SECRET", "SPACE_SECRET", "APP_SECRET"}
	for _, e := range envs {
		env := os.Getenv(e)
		if len(env) > 0 {
			data = append(data, env)
		}
	}
	if err := tmpl.Execute(w, data); err != nil {
		_, _ = fmt.Fprint(w, err)
	}
}
