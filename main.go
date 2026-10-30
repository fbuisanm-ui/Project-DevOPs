package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	// ============================================================
	// CONFIGURACIÓN DE ARCHIVOS ESTÁTICOS
	// ============================================================

	// HANDLE registra una ruta para servir archivos estáticos.
	//
	// Los archivos ubicados en la carpeta "static" podrán
	// consultarse desde el navegador mediante la ruta /static/.
	//
	// STRIPPrefix elimina "/static/" de la petición antes de
	// buscar el archivo dentro de la carpeta local "static".
	//
	// Por ejemplo:
	// http://localhost:8090/static/imagen.jpg
	//
	// buscará el archivo:
	// ./static/bladerunner.jpg
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ============================================================
	// CONFIGURACIÓN DE LA RUTA PRINCIPAL
	// ============================================================

	// HANDLEFunc registra una función que atenderá las peticiones
	// HTTP dirigidas a la ruta raíz ("/").
	//
	// Esta función se ejecutará cuando alguien acceda a:
	// http://localhost:8090/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// HEADER establece la cabecera HTTP de la respuesta.
		//
		// Content-Type indica al navegador qué tipo de contenido
		// está recibiendo.
		//
		// text/html indica que estamos enviando HTML.
		// charset=utf-8 permite mostrar correctamente caracteres
		// especiales, como tildes y eñes.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// FPRINTF escribe el contenido HTML en la respuesta HTTP.
		//
		// El navegador recibe este código y lo interpreta
		// para mostrar la página web.
		fmt.Fprint(w, `<!DOCTYPE html>
<html lang="es">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Project DevOps - UOC</title>
		<style>
		h1 {
			font-size: 120px;
		}
		</style>
	</head>
	<body>
		<h1>SOY ALUMNO DE LA UOC</h1>
		
		<!-- Imagen almacenada en la carpeta static del proyecto. -->
		<img src="/static/bladerunner.jpg" alt="Imagen del proyecto">
	</body>
</html>`)
	})

	// ============================================================
	// INICIO DEL SERVIDOR HTTP
	// ============================================================

	// PRINTLN muestra un mensaje en los registros (logs)
	// para indicar que el servidor está iniciado.
	log.Println("Servidor iniciado en el puerto 8090(DEV) 8091(PROD)")

	// LISTENANDSERVE inicia el servidor HTTP en el puerto 8090.
	//
	// El primer parámetro indica el puerto donde escuchará.
	// El segundo parámetro (nil) utiliza el gestor de rutas
	// predeterminado de Go, donde hemos registrado las rutas
	// anteriores mediante HANDLE y HANDLEFunc.
	//
	// LOG.Fatal muestra cualquier error y termina el programa
	// si el servidor no puede iniciarse o deja de funcionar.
	log.Fatal(http.ListenAndServe(":8090", nil))
}
