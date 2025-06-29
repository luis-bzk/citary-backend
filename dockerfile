FROM node:22 AS builder
WORKDIR /app_backend
COPY package*.json ./
RUN npm install
COPY . .

FROM node:22
WORKDIR /app_backend
COPY --from=builder /app_backend ./
# Copiar entrypoint.sh desde el contexto del build (citary-backend)
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
# Instalar netcat y dos2unix
RUN apt-get update && apt-get install -y netcat-traditional dos2unix && dos2unix /entrypoint.sh && apt-get clean && rm -rf /var/lib/apt/lists/*
EXPOSE 3000
ENTRYPOINT ["/entrypoint.sh"]