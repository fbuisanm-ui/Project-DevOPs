
# ============================================================
# ETAPA 1: COMPILACIÓN DE LA APLICACIÓN (BUILDER)
# ============================================================

# FROM indica qué imagen base vamos a utilizar.
# Una imagen es una plantilla que contiene un sistema y herramientas
# necesarias para realizar determinadas tareas.
#
# En este caso utilizamos golang:1.24, que incluye:
# - Sistema operativo Linux.
# - Compilador de Go.
# - Herramientas necesarias para construir nuestra aplicación.
#
# AS builder asigna el nombre "builder" a esta etapa.
# Esto nos permitirá recuperar posteriormente el ejecutable generado.
FROM golang:1.24 AS builder


# WORKDIR establece el directorio de trabajo dentro del contenedor.
# A partir de aquí, las instrucciones relativas se ejecutarán
# dentro de /app.
#
# Si el directorio no existe, Docker lo crea automáticamente.
WORKDIR /app


# COPY permite copiar archivos desde nuestro ordenador (contexto Docker)
# al interior de la imagen que estamos construyendo.
#
# Copiamos go.mod porque contiene la información del módulo de Go,
# como su nombre y versión del lenguaje.
COPY go.mod ./


# Copiamos el código fuente de nuestra aplicación.
# En este proyecto, el código principal está en main.go.
COPY main.go ./


# RUN ejecuta un comando durante la construcción de la imagen.
#
# CGO_ENABLED=0:
# Desactiva CGO para generar un ejecutable que no dependa de bibliotecas
# externas de C.
#
# GOOS=linux:
# Indica que queremos generar un ejecutable compatible con Linux.
# Esto es importante porque nuestro contenedor utilizará Alpine Linux.
#
# go build:
# Compila el código fuente de Go.
#
# -o server:
# Indica que el archivo ejecutable resultante se llamará "server".
#
# El punto final (.) indica que debe compilar el paquete del directorio
# actual, es decir, /app.
RUN CGO_ENABLED=0 GOOS=linux go build -o server .


# ============================================================
# ETAPA 2: EJECUCIÓN DE LA APLICACIÓN
# ============================================================

# Empezamos una segunda etapa utilizando Alpine Linux.
#
# Alpine es una distribución Linux muy ligera.
# No necesitamos instalar aquí el compilador de Go porque
# la aplicación ya está compilada.
#
# Esta técnica se llama compilación multi-stage (multi-stage build).
# Permite separar las herramientas de compilación de los elementos
# necesarios para ejecutar la aplicación.
FROM alpine:3.22


# Establecemos /app como directorio de trabajo de esta segunda etapa.
# Es importante recordar que cada etapa tiene su propio sistema
# de archivos.
WORKDIR /app


# COPY --from=builder permite copiar archivos desde otra etapa
# de construcción, en lugar de copiarlos desde nuestro ordenador.
#
# En este caso:
# /app/server = ejecutable generado en la primera etapa.
# ./server   = ubicación donde lo guardamos en esta segunda etapa.
#
# Así conseguimos ejecutar nuestra aplicación sin incluir
# el compilador completo de Go en la imagen final.
COPY --from=builder /app/server ./server


# Copiamos la carpeta static desde nuestro ordenador al contenedor.
#
# Esta carpeta contiene los recursos estáticos de nuestra página web,
# como imágenes, hojas de estilo CSS o archivos JavaScript.
#
# En nuestro proyecto necesitamos copiarla porque contiene imagen.jpg.
#
# Sin esta instrucción, Go podría arrancar correctamente, pero el
# navegador no encontraría la fotografía al solicitarla.
COPY static ./static


# EXPOSE documenta el puerto que utiliza nuestra aplicación
# dentro del contenedor.
#
# Nuestro servidor Go escucha en el puerto 8090.
#
# IMPORTANTE:
# EXPOSE no publica el puerto automáticamente en nuestro ordenador.
# Para eso utilizamos la configuración "ports" de Docker Compose.
EXPOSE 8090


# CMD define el comando que se ejecutará cuando arranque el contenedor.
#
# En este caso ejecutamos el archivo "server", que contiene
# nuestra aplicación Go ya compilada.
#
# Es la instrucción que pone en marcha nuestro servidor web.
CMD ["./server"]
