FROM nginx:latest

# Копируем файлы фронтенда в стандартную папку Nginx
COPY . /usr/share/nginx/html

# Указываем Nginx использовать index.html как главную страницу
RUN echo "server { \
    listen 80; \
    server_name localhost; \
    root /usr/share/nginx/html; \
    index index.html; \
    location / { \
        try_files \$uri /index.html; \
    } \
}" > /etc/nginx/conf.d/default.conf

# Открываем порт 80
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
