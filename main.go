package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="es">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Project DevOps - UOC</title>
</head>
<body>
	<h1>Esto es solo una prueba</h1>
	<p>No nos flipemos. What's your flavor?</p>
</body>
</html>`)
	})

	log.Println("Servidor iniciado en el puerto 8090(DEV) 8091(PROD)")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
