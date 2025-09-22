# Запуск:

Потребуется .env файл

```
DB_HOST=postgres
DB_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_DB=ozon_db
DB_SSLMODE=disable
```

Далее установить зависимости:

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

Внедрена моя собственная библиотека инмемори кэширования, с таймаутами, вытеснением LFU по колличеству запросов

На тестах повторные обращения работали значительно быстрее, примерно в 100 раз

![img.png](.assets/img_cache.png)

![img.png](.assets/img_cache2.png)

![img.png](.assets/img_cache3.png)

# Http-gateway

Я использую протобаф контракт, в котором использую gateway механизм, чтобы перенаправлять запросы с rest на grpc порт,
так же генерируется опенапи спецификация, которая исопльзуется для сваггера

# gRPC

Благодаря Validate в опциях протобафа, можно легко проверять, что входной параметр это ссылка
![img.png](.assets/img14.png)
![img.png](.assets/img13.png)

# Swagger

![img_2.png](.assets/img_2.png)

Создание ссылки
![img.png](.assets/img10.png)

Получение оригинальной
![img.png](.assets/img11.png)

Если сгенерировать повторно, выйдет соответствующая ошибка
![img.png](.assets/img.png)

Если такой ссылки нет, выйдет соответствующая ошибка
![img_1.png](.assets/img_1.png)

# Docker

![img.png](.assets/img12.png)