package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	_ "time/tzdata"
)

const (
	defaultColor = "#2563eb"
	defaultName  = "Azul"
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
		<p>IP privado: <span class="hostname">{{.PrivateIP}}</span></p>
		<p>Data e hora: <span class="hostname">{{.CurrentTime}}</span></p>
	</main>
</body>
</html>`))

type pageData struct {
	Name        string
	Color       string
	Hostname    string
	PrivateIP   string
	CurrentTime string
}

func privateIPv4() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "não encontrado"
	}

	for _, address := range addresses {
		ipNet, ok := address.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP.To4()
		if ip != nil && ip.IsPrivate() {
			return ip.String()
		}
	}
	return "não encontrado"
}

func disableCache(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "desconhecido"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		disableCache(w)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := pageData{
			Name:        defaultName,
			Color:       defaultColor,
			Hostname:    hostname,
			PrivateIP:   privateIPv4(),
			CurrentTime: time.Now().Format("02/01/2006 15:04:05.000 MST"),
		}
		if err := page.Execute(w, data); err != nil {
			log.Printf("erro ao renderizar a página: %v", err)
		}
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		disableCache(w)
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
