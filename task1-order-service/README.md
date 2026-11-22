# Базовый пример сервиса заказов из задания 1

Отдельный Go-модуль с Go kit, HTTP и in-memory репозиторием. Запуск из этого каталога:

```powershell
go mod tidy
go run .
```

Примеры проверки:

```powershell
$created = Invoke-RestMethod -Method Post -Uri "http://localhost:8080/orders" -ContentType "application/json" -Body '{"customer_id":"cust-1","restaurant_id":"rest-42","order_items":[{"product_code":"P001","name":"Пицца","unit_price":550,"quantity":1}]}'
$id = $created.id
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/orders/$id"
Invoke-RestMethod -Method Patch -Uri "http://localhost:8080/orders/$id/status" -ContentType "application/json" -Body '{"status":"Completed"}'
```
