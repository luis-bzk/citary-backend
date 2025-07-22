# Etapa base
FROM node:22 AS base
WORKDIR /app_backend
COPY package*.json ./

# Etapa de desarrollo
FROM base AS development
RUN npm install
COPY . .
EXPOSE 3001
CMD ["npm", "run", "dev"]

# Etapa de build
FROM base AS builder
RUN npm install
COPY . .
RUN npm run build

# Etapa de producción
FROM node:22 AS production
WORKDIR /app_backend
COPY --from=builder /app_backend ./
# Copiar entrypoint.sh desde el contexto del build (nubrik-backend)
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
# Instalar netcat y dos2unix
RUN apt-get update && apt-get install -y netcat-traditional dos2unix && dos2unix /entrypoint.sh && apt-get clean && rm -rf /var/lib/apt/lists/*
EXPOSE 3000
ENTRYPOINT ["/entrypoint.sh"]