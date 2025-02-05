# Сервис для интеграции с judge0

Для сборки приложения:
```bash
docker build -t code-judge-service .
```

Для запуска приложения:
```bash
docker run -p 8080:8080 \
  -e DB_CONN="postgres://user:pass@host/db" \
  -e KAFKA_BROKERS="kafka:9092" \
  code-judge-service
```

Для локального запуска:
Поправить `.env` файл в соответствии с нуждами)))

Запустить миграцию:
```bash
golang-migrate.exe -database "postgres://{DB_USER}:{DB_PASSWORD}@{HOST}:{PORT}/{DB_NAME}?sslmode=disable" -path ./migrations up
```
(ну либо ручками ее бахнуть))) )

Запустить docker-compose
```bash
docker-compose up -d
```

Запустить сервис и радоваться жизни))
(на всякий случай в папке `mock` лежит моковый сервис, который можно раскатать, чтобы проэмулировать деятельность judge0 ;0) )