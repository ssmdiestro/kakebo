# ---- build ----
    FROM node:20-alpine AS build
    WORKDIR /app
    
    # 1) Copiamos solo manifests para cachear deps
    COPY package.json package-lock.json ./
    
    # 2) Fuerza instalar devDependencies aunque exista .npmrc con production=true
    #    (npm v9+ acepta --include=dev; equivalente a no omitir dev)
    RUN npm ci --include=dev
    
    # 3) Copiamos el resto del código
    COPY . .
    
    # 4) (Guardarraíles) Arregla posibles CRLF y permisos en shims binarios
    RUN apk add --no-cache dos2unix \
     && find node_modules/.bin -maxdepth 1 -type f -exec dos2unix {} \; \
     && chmod +x node_modules/.bin/*
    
    # 5) Comprobación opcional (útil para diagnosticar)
    RUN npx vite --version
    
    # 6) Build
    RUN npm run build   # genera /app/dist
    
    
    # ---- run ----
    FROM nginx:alpine
    # (opcional) elimina WORKDIR, Nginx ya sirve /usr/share/nginx/html por defecto
    COPY --from=build /app/dist /usr/share/nginx/html
    COPY ./nginx.conf /etc/nginx/nginx.conf
    EXPOSE 80
    CMD ["nginx","-g","daemon off;"]
    