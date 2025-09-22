# Запуск:

потребуется .env файл
```
DB_HOST=postgres
DB_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_DB=ozon_db
DB_SSLMODE=disable
```

далее установить зависимости:
```bash
  make grpc-deps
```

И запустить контейнеры
```bash
  make up
```
или
```bash
  docker-compose up
```
# cache

внедрена моя собственная библиотека инмемори кэширования, с таймаутами, вытеснением LFU по колличеству запросов

на тестах повторные обращения работали значительно быстрее, примерно в 100 раз

![img.png](.assets/img_cache.png)

![img.png](.assets/img_cache2.png)

![img.png](.assets/img_cache3.png)


# Http-gateway

Я исользую прото баф контракт, в котором использую gateway механизм, чтобы перенаправлять запросы с rest на grpc порт, так же генерируется опенапи спецификация, которая исопльзуется для сваггера


# gRPC 

Благодаря Validate в опциях протобафа, можно легко проверять, что входной параметр это ссылка
![img.png](.assets/img14.png)
![img.png](.assets/img13.png)

# Swagger

![img_2.png](.assets/img_2.png)

создание ссылки
![img.png](.assets/img10.png)

получение оригинальной
![img.png](.assets/img11.png)

если сгенерировать повторно, выйдет соответствующая ошибка
![img.png](.assets/img.png)

если такой ссылки нет, выйдет соответствующая ошибка
![img_1.png](.assets/img_1.png)

# Docker
![img.png](.assets/img12.png)