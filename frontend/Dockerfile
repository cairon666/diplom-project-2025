# Многоэтапная сборка для React приложения
FROM node:20-alpine AS builder

# Включаем pnpm через corepack
RUN corepack enable

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем package.json и pnpm-lock.yaml для кэширования зависимостей
COPY package.json pnpm-lock.yaml* ./

# Устанавливаем зависимости (включая dev dependencies для сборки)
RUN pnpm install --frozen-lockfile

# Копируем исходный код
COPY . .

# Собираем приложение для продакшена
RUN pnpm run build

# Продакшн стадия с nginx
FROM nginx:alpine

# Копируем собранное приложение
COPY --from=builder /app/dist /usr/share/nginx/html

# Копируем кастомную конфигурацию nginx для SPA
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Открываем порт
EXPOSE 80

# Запускаем nginx
CMD ["nginx", "-g", "daemon off;"] 