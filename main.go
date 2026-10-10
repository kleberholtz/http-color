package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	defaultColor = "#dc2626"
	defaultName  = "Vermelho"
)

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="pt-BR">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{.Name}} | Go Color</title>
	<style>
		* { box-sizing: border-box; }
		body {
			min-height: 100vh;
			margin: 0;
			display: grid;
			place-items: center;
			background: {{.Color}};
			color: white;
			font: 16px/1.5 system-ui, sans-serif;
		}
		main { padding: 2rem; text-align: center; }
		h1 { margin: 0 0 1rem; font-size: clamp(3rem, 12vw, 7rem); }
		p { margin: .5rem 0; }
		.hostname { font-family: ui-monospace, monospace; }
	</style>
</head>
<body>
	<main>
		<h1>{{.Name}}</h1>
		<p>Cor servida por uma aplicação Go</p>
		<p>Hostname: <span class="hostname">{{.Hostname}}</span></p>
	</main>
</body>
</html>`))

type pageData struct {
	Name     string
	Color    string
	Hostname string
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "desconhecido"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := page.Execute(w, pageData{Name: defaultName, Color: defaultColor, Hostname: hostname}); err != nil {
			log.Printf("erro ao renderizar a página: %v", err)
		}
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "8080"
	}
	log.Printf("servidor iniciado na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
