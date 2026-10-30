# Project-DevOPs
Proyecto colaborativo UOC con Go y Docker, con entornos de desarrollo y producción.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Mensaje struct {
	ID    int    `json:"id"`
	Texto string `json:"texto"`
}

func main() {

	mensajes := make(map[int]Mensaje)
	siguienteID := 1
	var mutex sync.Mutex
	mux := http.NewServeMux()

	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("./static")),
		),
	)

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
				<h1>SOY ALUMNO DE LA UOC</h1>
				<img src="/static/bladerunner.jpg" alt="Imagen del proyecto" width="600">
			</body>
			</html>
		`)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("GET /api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"proyecto": "Project DevOps",
			"lenguaje": "Go",
			"estado":   "funcionando",
		})
	})

	mux.HandleFunc("GET /api/mensajes", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mensajes)
	})

	mux.HandleFunc("GET /api/mensajes/{id}", func(w http.ResponseWriter, r *http.Request) {
		idTexto := r.PathValue("id")
		id, err := strconv.Atoi(idTexto)
		if err != nil {
			http.Error(w, "Identificador no válido", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		mensaje, existe := mensajes[id]
		if !existe {
			http.Error(w, "Mensaje no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mensaje)
	})

	mux.HandleFunc("POST /api/mensajes", func(w http.ResponseWriter, r *http.Request) {
		var mensaje Mensaje
		err := json.NewDecoder(r.Body).Decode(&mensaje)
		if err != nil {
			http.Error(w, "JSON incorrecto", http.StatusBadRequest)
			return
		}

		if mensaje.Texto == "" {
			http.Error(w, "El texto del mensaje no puede estar vacío", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		mensaje.ID = siguienteID
		mensajes[siguienteID] = mensaje
		siguienteID++

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(mensaje)
	})

	mux.HandleFunc("PUT /api/mensajes/{id}", func(w http.ResponseWriter, r *http.Request) {
		idTexto := r.PathValue("id")
		id, err := strconv.Atoi(idTexto)
		if err != nil {
			http.Error(w, "Identificador no válido", http.StatusBadRequest)
			return
		}

		var mensajeActualizado Mensaje
		err = json.NewDecoder(r.Body).Decode(&mensajeActualizado)
		if err != nil {
			http.Error(w, "JSON incorrecto", http.StatusBadRequest)
			return
		}

		if mensajeActualizado.Texto == "" {
			http.Error(w, "El texto del mensaje no puede estar vacío", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		_, existe := mensajes[id]
		if !existe {
			http.Error(w, "Mensaje no encontrado", http.StatusNotFound)
			return
		}

		mensajeActualizado.ID = id
		mensajes[id] = mensajeActualizado

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mensajeActualizado)
	})

	mux.HandleFunc("DELETE /api/mensajes/{id}", func(w http.ResponseWriter, r *http.Request) {
		idTexto := r.PathValue("id")
		id, err := strconv.Atoi(idTexto)
		if err != nil {
			http.Error(w, "Identificador no válido", http.StatusBadRequest)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		_, existe := mensajes[id]
		if !existe {
			http.Error(w, "Mensaje no encontrado", http.StatusNotFound)
			return
		}

		delete(mensajes, id)
		w.WriteHeader(http.StatusNoContent)
	})

	servidor := &http.Server{
		Addr:         ":8090",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Servidor iniciado en http://localhost:8090")

	err := servidor.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}