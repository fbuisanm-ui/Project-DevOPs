package main

import (
	"encoding/json" // Permite convertir estructuras de Go a JSON y viceversa.
	"fmt"           // Permite crear mensajes de texto con formato.
	"log"           // Permite mostrar información y errores en la consola.
	"net/http"      // Proporciona las herramientas para crear el servidor web.
	"strconv"       // Permite convertir cadenas de texto en números enteros.
	"sync"          // Permite proteger los datos cuando llegan varias peticiones a la vez.
	"time"          // Permite configurar los tiempos máximos del servidor.
)

// ---------------------------------------------------------
// ESTRUCTURA DE DATOS
// ---------------------------------------------------------

// Mensaje representa cada mensaje que se puede crear,
// consultar, modificar o eliminar mediante la API.
//
// Los nombres de los campos JSON permiten que los datos
// se intercambien fácilmente con clientes externos.
type Mensaje struct {
	ID    int    `json:"id"`
	Texto string `json:"texto"`
}

// ---------------------------------------------------------
// FUNCIÓN PRINCIPAL
// ---------------------------------------------------------

func main() {

	// Creamos un mapa para almacenar los mensajes en memoria.
	// La clave será el identificador y el valor será el mensaje.
	//
	// Ejemplo:
	// 1 -> {ID: 1, Texto: "Hola"}
	// 2 -> {ID: 2, Texto: "Segundo mensaje"}
	mensajes := make(map[int]Mensaje)

	// Variable que utilizaremos para asignar un identificador
	// diferente a cada mensaje nuevo.
	siguienteID := 1

	// Creamos un Mutex para proteger el mapa de mensajes.
	//
	// Esto evita problemas si dos peticiones intentan modificar
	// los datos al mismo tiempo.
	var mutex sync.Mutex

	// Creamos un multiplexor de rutas.
	//
	// El ServeMux permite definir diferentes rutas y métodos HTTP
	// para que cada petición llegue a su función correspondiente.
	mux := http.NewServeMux()

	// ---------------------------------------------------------
	// ARCHIVOS ESTÁTICOS
	// ---------------------------------------------------------

	// Registramos la carpeta static para servir archivos como
	// imágenes, hojas de estilo CSS o archivos JavaScript.
	//
	// StripPrefix elimina "/static/" de la ruta antes de buscar
	// el archivo dentro de la carpeta "static".
	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	// ---------------------------------------------------------
	// RUTA PRINCIPAL: GET /
	// ---------------------------------------------------------

	// Esta ruta muestra la página principal de la aplicación.
	// El patrón /{$} indica que solamente responde a la raíz.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {

		// Indicamos que la respuesta tendrá contenido HTML.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Mostramos una página HTML sencilla.
		fmt.Fprint(w, `
			<!DOCTYPE html>
			<html lang="es">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Proyecto DevOps</title>
			</head>
			<style>
        		h1 {
            		font-size: 40px;
        			}
    		</style>
			<body>
				<h1>SOY ALUMNO DE LA UOC: JENKINS TEST 5</h1>
				<img src="/static/bladerunner.jpg" alt="Imagen del proyecto" width="600">
			</body>
			</html>
		`)
	})

	// ---------------------------------------------------------
	// RUTA DE SALUD: GET /health
	// ---------------------------------------------------------

	// Esta ruta permite comprobar si el servidor está funcionando.
	// Es útil para monitorización y comprobaciones de disponibilidad.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {

		// Indicamos que la respuesta se enviará en formato JSON.
		w.Header().Set("Content-Type", "application/json")

		// Devolvemos un código HTTP 200 (OK).
		w.WriteHeader(http.StatusOK)

		// Enviamos el estado del servidor.
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	// ---------------------------------------------------------
	// RUTA DE INFORMACIÓN: GET /api/info
	// ---------------------------------------------------------

	// Esta ruta proporciona información general de la aplicación.
	mux.HandleFunc("GET /api/info", func(w http.ResponseWriter, r *http.Request) {

		// Indicamos que la respuesta será JSON.
		w.Header().Set("Content-Type", "application/json")

		// Devolvemos información del proyecto.
		json.NewEncoder(w).Encode(map[string]string{
			"proyecto": "Project DevOps",
			"lenguaje": "Go",
			"estado":   "funcionando",
		})
	})

	// ---------------------------------------------------------
	// CONSULTAR TODOS LOS MENSAJES: GET /api/mensajes
	// ---------------------------------------------------------

	// Esta ruta devuelve todos los mensajes almacenados.
	mux.HandleFunc("GET /api/mensajes", func(w http.ResponseWriter, r *http.Request) {

		// Bloqueamos temporalmente el acceso al mapa para evitar
		// que otra petición lo modifique mientras lo consultamos.
		mutex.Lock()
		defer mutex.Unlock()

		// Indicamos que la respuesta será JSON.
		w.Header().Set("Content-Type", "application/json")

		// Convertimos el mapa de mensajes a formato JSON.
		json.NewEncoder(w).Encode(mensajes)
	})

	// ---------------------------------------------------------
	// CONSULTAR UN MENSAJE: GET /api/mensajes/{id}
	// ---------------------------------------------------------

	// Esta ruta permite consultar un mensaje concreto utilizando
	// su identificador.
	//
	// Por ejemplo:
	// GET /api/mensajes/1
	mux.HandleFunc("GET /api/mensajes/{id}", func(w http.ResponseWriter, r *http.Request) {

		// Recuperamos el identificador de la URL.
		idTexto := r.PathValue("id")

		// Convertimos el identificador de texto a número entero.
		id, err := strconv.Atoi(idTexto)

		// Si el identificador no es un número válido, devolvemos 400.
		if err != nil {
			http.Error(w, "Identificador no válido", http.StatusBadRequest)
			return
		}

		// Bloqueamos el acceso al mapa mientras consultamos el mensaje.
		mutex.Lock()
		defer mutex.Unlock()

		// Buscamos el mensaje mediante su identificador.
		mensaje, existe := mensajes[id]

		// Si no existe, devolvemos un error 404.
		if !existe {
			http.Error(w, "Mensaje no encontrado", http.StatusNotFound)
			return
		}

		// Indicamos que la respuesta será JSON.
		w.Header().Set("Content-Type", "application/json")

		// Devolvemos el mensaje encontrado.
		json.NewEncoder(w).Encode(mensaje)
	})

	// ---------------------------------------------------------
	// CREAR UN MENSAJE: POST /api/mensajes
	// ---------------------------------------------------------

	// Esta ruta permite crear un mensaje nuevo enviando
	// los datos mediante una petición POST.
	mux.HandleFunc("POST /api/mensajes", func(w http.ResponseWriter, r *http.Request) {

		// Creamos una variable para almacenar los datos recibidos.
		var mensaje Mensaje

		// Decodificamos el JSON recibido en el cuerpo de la petición.
		err := json.NewDecoder(r.Body).Decode(&mensaje)

		// Si el JSON no tiene un formato válido, devolvemos 400.
		if err != nil {
			http.Error(w, "JSON incorrecto", http.StatusBadRequest)
			return
		}

		// Comprobamos que el texto del mensaje no esté vacío.
		if mensaje.Texto == "" {
			http.Error(w, "El texto del mensaje no puede estar vacío", http.StatusBadRequest)
			return
		}

		// Bloqueamos el acceso al mapa porque vamos a modificarlo.
		mutex.Lock()
		defer mutex.Unlock()

		// Asignamos un identificador único al mensaje.
		mensaje.ID = siguienteID

		// Guardamos el mensaje en el mapa.
		mensajes[siguienteID] = mensaje

		// Incrementamos el identificador para el siguiente mensaje.
		siguienteID++

		// Indicamos que la respuesta será JSON.
		w.Header().Set("Content-Type", "application/json")

		// Indicamos que el recurso se ha creado correctamente.
		w.WriteHeader(http.StatusCreated)

		// Devolvemos el mensaje creado.
		json.NewEncoder(w).Encode(mensaje)
	})

	// -------------------------------------------------------
	// ACTUALIZAR UN MENSAJE: PUT /api/mensajes/{id}
	// -------------------------------------------------------

	// Esta ruta permite modificar el texto de un mensaje existente.
	mux.HandleFunc("PUT /api/mensajes/{id}", func(w http.ResponseWriter, r *http.Request) {

		// Recuperamos el identificador desde la URL.
		idTexto := r.PathValue("id")

		// Convertimos el identificador a número entero.
		id, err := strconv.Atoi(idTexto)

		// Si el identificador no es válido, devolvemos 400.
		if err != nil {
			http.Error(w, "Identificador no válido", http.StatusBadRequest)
			return
		}

		// Creamos una variable para recibir los nuevos datos.
		var mensajeActualizado Mensaje

		// Decodificamos el JSON enviado por el cliente.
		err = json.NewDecoder(r.Body).Decode(&mensajeActualizado)

		// Si el JSON es incorrecto, devolvemos 400.
		if err != nil {
			http.Error(w, "JSON incorrecto", http.StatusBadRequest)
			return
		}

		// Comprobamos que el nuevo texto no esté vacío.
		if mensajeActualizado.Texto == "" {
			http.Error(w, "El texto del mensaje no puede estar vacío", http.StatusBadRequest)
			return
		}

		// Bloqueamos el mapa porque vamos a consultar y modificar datos.
		mutex.Lock()
		defer mutex.Unlock()

		// Comprobamos si el mensaje existe.
		_, existe := mensajes[id]

		// Si no existe, devolvemos un error 404.
		if !existe {
			http.Error(w, "Mensaje no encontrado", http.StatusNotFound)
			return
		}

		// Conservamos el identificador de la URL y actualizamos el texto.
		mensajeActualizado.ID = id

		// Guardamos el mensaje actualizado en el mapa.
		mensajes[id] = mensajeActualizado

		// Indicamos que la respuesta será JSON.
		w.Header().Set("Content-Type", "application/json")

		// Devolvemos el mensaje actualizado con código 200.
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mensajeActualizado)
	})

	// ---------------------------------------------------------
	// ELIMINAR UN MENSAJE: DELETE /api/mensajes/{id}
	// ---------------------------------------------------------

	// Esta ruta permite eliminar un mensaje mediante su identificador.
	mux.HandleFunc("DELETE /api/mensajes/{id}", func(w http.ResponseWriter, r *http.Request) {

		// Recuperamos el identificador desde la URL.
		idTexto := r.PathValue("id")

		// Convertimos el identificador a número entero.
		id, err := strconv.Atoi(idTexto)

		// Si el identificador no es válido, devolvemos 400.
		if err != nil {
			http.Error(w, "Identificador no válido", http.StatusBadRequest)
			return
		}

		// Bloqueamos el mapa porque vamos a consultar y eliminar datos.
		mutex.Lock()
		defer mutex.Unlock()

		// Comprobamos si el mensaje existe.
		_, existe := mensajes[id]

		// Si no existe, devolvemos un error 404.
		if !existe {
			http.Error(w, "Mensaje no encontrado", http.StatusNotFound)
			return
		}

		// Eliminamos el mensaje del mapa.
		delete(mensajes, id)

		// Devolvemos el código 204, que indica que la operación
		// se ha realizado correctamente y no necesita devolver contenido.
		w.WriteHeader(http.StatusNoContent)
	})

	// ---------------------------------------------------------
	// CONFIGURACIÓN DEL SERVIDOR HTTP
	// ---------------------------------------------------------

	// Creamos una instancia del servidor HTTP con diferentes
	// parámetros para mejorar su funcionamiento y seguridad.
	servidor := &http.Server{

		// Puerto en el que escuchará el servidor.
		Addr: ":8090",

		// Asociamos el multiplexor con todas las rutas registradas.
		Handler: mux,

		// Tiempo máximo para leer la petición completa.
		ReadTimeout: 5 * time.Second,

		// Tiempo máximo para escribir la respuesta.
		WriteTimeout: 10 * time.Second,

		// Tiempo máximo que una conexión puede permanecer inactiva.
		IdleTimeout: 60 * time.Second,
	}

	// Mostramos un mensaje para saber que el servidor está arrancando.
	log.Println("Servidor iniciado en http://localhost:8090")

	// Iniciamos el servidor y esperamos peticiones HTTP.
	//
	// Si se produce un error al arrancar o mantener el servidor,
	// lo mostramos en la consola.
	err := servidor.ListenAndServe()

	// ErrServerClosed es el error esperado cuando el servidor
	// se cierra de forma controlada. No lo tratamos como un fallo.
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
