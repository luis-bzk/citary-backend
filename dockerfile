FROM node:22.12.0

WORKDIR /app_backend

COPY package*.json ./
RUN npm install --omit=dev

COPY . .

EXPOSE 3000

CMD ["npm", "start"]
