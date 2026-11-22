# Лабораторная работа № 6 — основы микросервисной архитектуры

**Индивидуальный вариант 3: микросервис «Студенты».** Реализация на Go и Go kit. Один сервис предоставляет одинаковые бизнес-операции через HTTP/JSON (`:8080`) и gRPC (`:9090`). Для хранения используется потокобезопасный in-memory репозиторий, поэтому после перезапуска сервера данные сбрасываются.

## Структура

```text
lab6_student_service/
├── student/                  # Student, Service, Repository и доменные ошибки
├── implementation/           # бизнес-логика и модульные тесты
├── repository/                # потокобезопасное хранилище map + RWMutex
├── transport/                 # Go kit endpoints и DTO
├── httptransport/             # HTTP-маршруты, JSON-кодирование и ошибки
├── grpctransport/             # Go kit gRPC transport и преобразование DTO
├── proto/student.proto        # контракт StudentService
├── gen/student/v1/            # Go-описания сообщений и gRPC-интерфейса
├── task1-order-service/       # пример из задания 1 (HTTP, сервис заказов)
├── main.go                    # запуск HTTP и gRPC параллельно
└── go.mod
```

Вариантный сервис разделяет бизнес-логику и сетевые протоколы. Путь запроса: `HTTP/gRPC transport → endpoint → Service → Repository`. Обе транспортные оболочки используют один экземпляр сервиса и одно in-memory хранилище.

## Требования

- Go 1.26.0 или новее (версия указана в корневом `go.mod`).
- Доступ к сети при первом запуске `go mod tidy`, чтобы Go мог загрузить зависимости.
- Для повторной генерации protobuf-файлов: `protoc`, установленный в `PATH`, и доступ к сети для установки Go-плагинов.
- Для проверки gRPC через CLI: `grpcurl`.

## Запуск варианта 3

Откройте PowerShell или терминал в корневом каталоге `lab6_student_service`:

```powershell
go mod tidy
go test ./...
go run .
```

По умолчанию запускаются HTTP на `localhost:8080` и gRPC на `localhost:9090`. Порты можно изменить флагами:

```powershell
go run . -http.addr=:8081 -grpc.addr=:9091
```

Остановка через `Ctrl+C` запускает корректное завершение HTTP- и gRPC-серверов.

## Проверка HTTP

В PowerShell выполните команды по порядку. Первый запрос создаёт студента; возвращённый идентификатор затем используется для чтения и изменения группы.

```powershell
$created = Invoke-RestMethod -Method Post `
  -Uri "http://localhost:8080/students" `
  -ContentType "application/json" `
  -Body '{"full_name":"Иванов Иван Иванович","group":"ПИЖ-б-о-23-2(2)","gpa":4.6}'

$id = $created.id
$created

Invoke-RestMethod -Method Get -Uri "http://localhost:8080/students/$id"

Invoke-RestMethod -Method Patch `
  -Uri "http://localhost:8080/students/$id/group" `
  -ContentType "application/json" `
  -Body '{"group":"ПИЖ-б-о-23-3"}'

Invoke-RestMethod -Method Get -Uri "http://localhost:8080/students/$id"
```

HTTP-контракт:

| Метод и путь | Назначение | Успешный ответ |
|---|---|---|
| `POST /students` | Создать студента | `{"id":"..."}` |
| `GET /students/{id}` | Получить студента | `{"student":{"id":"...","full_name":"...","group":"...","gpa":4.6,"enrolled_on":...}}` |
| `PATCH /students/{id}/group` | Изменить группу | `{"changed":true}` |

Поле `enrolled_on` задаётся сервисом в формате Unix timestamp (секунды), если его не передали. `id` генерируется автоматически. Пустое ФИО или группа считаются некорректными данными. Для неизвестного ID возвращается HTTP 404, для некорректного запроса — HTTP 400.

## Проверка gRPC

Сервис использует тот же доменный слой. Укажите файл `.proto` вручную, поэтому серверная reflection для проверки не требуется. Команды ниже предназначены для `grpcurl`:

```powershell
grpcurl -plaintext -import-path . -proto proto/student.proto `
  -d '{"student":{"fullName":"Иванов Иван Иванович","group":"ПИЖ-б-о-23-2(2)","gpa":4.6}}' `
  localhost:9090 student.v1.StudentService/Create
```

Скопируйте `id` из ответа, затем выполните:

```powershell
grpcurl -plaintext -import-path . -proto proto/student.proto `
  -d '{"id":"ВСТАВЬТЕ_ID_ИЗ_ОТВЕТА"}' `
  localhost:9090 student.v1.StudentService/GetByID

grpcurl -plaintext -import-path . -proto proto/student.proto `
  -d '{"id":"ВСТАВЬТЕ_ID_ИЗ_ОТВЕТА","group":"ПИЖ-б-о-23-3"}' `
  localhost:9090 student.v1.StudentService/ChangeGroup
```

В gRPC используется стандартное JSON-имя protobuf `fullName` для поля `full_name`. Ошибка отсутствующего студента передаётся с кодом gRPC `NotFound`, некорректный ввод — `InvalidArgument`.

## Генерация protobuf-кода

Контракт расположен в `proto/student.proto`. Скрипт для PowerShell устанавливает генераторы и вызывает `protoc`:

```powershell
.\scripts\generate-proto.ps1
```

Перед запуском убедитесь, что `protoc.exe` доступен через `PATH`. Генерация заменяет файлы в `gen/student/v1` официальным кодом плагинов `protoc-gen-go` и `protoc-gen-go-grpc`.

## Проверка качества

Из корневого каталога проекта:

```powershell
go test ./...
go test -race ./...
go vet ./...
```

Модульные тесты проверяют создание студента, автоматическое формирование ID и времени зачисления, получение по ID, смену группы и отказ при пустом ФИО. В среде подготовки успешно выполнены `go test ./student ./repository` и `go test -race ./student ./repository`. Полную сборку и тестирование HTTP/gRPC здесь подтвердить не удалось: Go Module Proxy был недоступен, а внешние зависимости не были загружены.

## Задание 1 — пример сервиса заказов

В каталоге `task1-order-service` сохранён самостоятельный модуль, повторяющий основную архитектурную идею примера из методички. Для его запуска откройте отдельный терминал:

```powershell
cd task1-order-service
go mod tidy
go run .
```

Он поднимает HTTP на `localhost:8080` и предоставляет `POST /orders`, `GET /orders/{id}`, `PATCH /orders/{id}/status`. Остановите вариантный сервис, если порт `8080` уже занят.


## Лабораторная работа № 7 — CI/CD

В проект добавлены GitHub Actions workflows для автоматической сборки и тестирования (`.github/workflows/ci.yml`), публикации Docker-образа (`.github/workflows/docker.yml`) и создания релизов по тегам `v*` (`.github/workflows/release.yml`). Основной сервис собирается на Go 1.26.x и 1.27.x; `coverage.out` сохраняется как artifact. Для Docker Hub в настройках GitHub Actions необходимо задать секреты `DOCKERHUB_USERNAME` и `DOCKERHUB_TOKEN`.

Подробная пошаговая инструкция по настройке GitHub, Docker Hub, проверке контейнера, защите ветки `main` и созданию релиза находится в [CI_CD_GUIDE.md](CI_CD_GUIDE.md).

После создания GitHub-репозитория можно добавить в начало README бейджи, заменив `YOUR_GITHUB_USERNAME` и `student-service` на фактические значения:

```markdown
[![CI](https://github.com/YOUR_GITHUB_USERNAME/student-service/actions/workflows/ci.yml/badge.svg)](https://github.com/YOUR_GITHUB_USERNAME/student-service/actions/workflows/ci.yml)
[![Docker](https://github.com/YOUR_GITHUB_USERNAME/student-service/actions/workflows/docker.yml/badge.svg)](https://github.com/YOUR_GITHUB_USERNAME/student-service/actions/workflows/docker.yml)
```
<!-- test branch protection -->