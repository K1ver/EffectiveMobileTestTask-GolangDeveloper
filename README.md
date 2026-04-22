# EffectiveMobileTestTask-GolangDeveloper
Test task for effective mobile

# Генерация swagger-документации
swag init -g cmd/main.go -o docs

# Запуск через docker compose
docker compose up --build

# Тесты
go test ./... -v

# Swagger UI доступен по адресу
# http://localhost:8080/swagger/index.html