# DEV ONLY
FROM golang:1.23-alpine AS dev
WORKDIR /app

# Para que respete la versión del go.mod (1.23.x)
ENV GOTOOLCHAIN=auto

# Herramientas útiles (opcional)
RUN apk add --no-cache git bash

# Cachea dependencias
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct && go mod download

# El código real vendrá por volumen; esto evita copiar todo aquí.
# Puedes copiarlo si quieres: COPY . .

# Ejecuta la app en caliente
CMD ["go", "run", "./cmd/api/main.go"]
