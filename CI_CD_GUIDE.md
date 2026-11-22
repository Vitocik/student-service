# Лабораторная работа № 7 — CI/CD для student-service

## 1. Что находится в проекте

Это сервис «Студенты» из лабораторной работы № 6 (вариант 3). Он предоставляет HTTP API на порту `8080` и gRPC API на порту `9090`. В проекте уже есть модульные тесты в каталогах `implementation` и `repository`.

Для основной части ЛР 7 добавлены:

- `.github/workflows/ci.yml` — автоматическая проверка кода на Go 1.26 и Go 1.27, `go vet`, golangci-lint, сборка, тесты с race detector и сохранение `coverage.out` как artifact;
- `Dockerfile` и `.dockerignore` — multi-stage сборка контейнера;
- `.github/workflows/docker.yml` — сборка и публикация Docker-образа в Docker Hub;
- `.github/workflows/release.yml` — индивидуальный вариант 3 ЛР 7: при push тега `v*` собираются бинарные файлы Linux, Windows и macOS и создаётся GitHub Release.

**Важно о версии Go.** В `go.mod` корневого проекта указано `go 1.26.0`. Поэтому приведённый в методичке пример с Go 1.21/1.22 и Docker-образом `golang:1.22-alpine` к этому архиву не подходит. В workflow используются Go 1.26.x и 1.27.x. Официальная история релизов Go подтверждает, что обе ветки доступны на дату выполнения работы.

## 2. Подготовка локального компьютера

Откройте корневую папку `lab6_student_service` в VS Code и откройте встроенный терминал: **Terminal → New Terminal**.

Проверьте установленные программы:

```powershell
go version
git --version
docker --version
```

Для корневого сервиса требуется Go не ниже 1.26.0. Если `go version` показывает версию ниже 1.26, установите Go 1.26.x или 1.27.x и перезапустите терминал. Docker Desktop должен быть запущен, если планируете локально запускать контейнер.

Проверьте проект до загрузки в GitHub:

```powershell
go mod download
go mod verify
go test ./...
go test -race ./...
go vet ./...
go build ./...
```

Каждая команда должна завершиться без ошибки. Команда `go test ./...` проверяет пакеты корневого сервиса. Каталог `task1-order-service` — отдельный вложенный Go-модуль из примера задания 1, он не входит в `./...` корневого модуля.

## 3. Создание репозитория GitHub и первая загрузка ЛР 6

1. Откройте GitHub и создайте новый репозиторий, например `student-service`.
2. На этапе создания не добавляйте README, `.gitignore` или LICENSE: файлы уже есть в локальном проекте.
3. Скопируйте HTTPS-адрес репозитория. Ниже вместо `YOUR_GITHUB_USERNAME` укажите свой логин GitHub.

В терминале, находясь именно в корне `lab6_student_service`, выполните:

```powershell
git status
```

Если команда сообщает, что это ещё не Git-репозиторий, выполните:

```powershell
git init
git add .
git commit -m "Add student service from lab 6"
git branch -M main
```

Если Git-репозиторий уже существует, `git init` повторять не нужно. Добавьте новые файлы и создайте коммит обычным способом:

```powershell
git add .
git commit -m "Prepare project for lab 7"
git branch -M main
```

Проверьте подключённые адреса:

```powershell
git remote -v
```

Если `origin` ещё не настроен:

```powershell
git remote add origin https://github.com/YOUR_GITHUB_USERNAME/student-service.git
```

Если `origin` уже есть, но указывает на другой репозиторий:

```powershell
git remote set-url origin https://github.com/YOUR_GITHUB_USERNAME/student-service.git
```

Загрузите проект:

```powershell
git push -u origin main
```

Если GitHub просит пароль при HTTPS-аутентификации, обычный пароль от аккаунта не подходит. Используйте вход через Git Credential Manager/браузер либо Personal Access Token в качестве пароля. Никому не отправляйте токен и не вставляйте его в исходники.

После push откройте страницу репозитория и убедитесь, что видны `go.mod`, `main.go`, `implementation`, `repository`, `student`, `Dockerfile` и `.github/workflows`.

## 4. Настройка Docker Hub и секретов GitHub

Для публикации Docker-образа:

1. Создайте/откройте аккаунт Docker Hub.
2. Создайте репозиторий с именем `student-service` (лучше Public, чтобы можно было скачать образ без авторизации).
3. В настройках Docker Hub создайте Personal Access Token с разрешением Read & Write. Скопируйте токен сразу и храните его как пароль.
4. В GitHub откройте репозиторий → **Settings → Secrets and variables → Actions → New repository secret**.
5. Создайте два секрета:
   - `DOCKERHUB_USERNAME` — именно логин Docker Hub, не адрес электронной почты;
   - `DOCKERHUB_TOKEN` — созданный Personal Access Token.

Не записывайте токен в `docker.yml`, README или команду Git. Workflow получает его через `${{ secrets.DOCKERHUB_TOKEN }}`.

## 5. CI: автоматическая сборка, тесты и линтинг

Файл `.github/workflows/ci.yml` запускается при push в `main`, при pull request в `main`, а также вручную из вкладки Actions.

Внутри job выполняются последовательно:

1. `actions/checkout` получает исходный код;
2. `actions/setup-go` устанавливает Go из матрицы: 1.26.x и 1.27.x;
3. `go mod download` загружает модули, `go mod verify` проверяет их целостность;
4. `go vet ./...` выполняет статический анализ стандартными средствами Go;
5. `golangci-lint` запускает набор линтеров;
6. `go build -v ./...` собирает пакеты;
7. `go test -v -race -coverprofile=coverage.out ./...` запускает тесты, детектор гонок и создаёт профиль покрытия;
8. `actions/upload-artifact` сохраняет `coverage.out` отдельно для каждой версии Go.

Кэширование включено в `actions/setup-go`: оно использует `go.sum` как зависимость ключа и повторно использует кэш модулей и сборки при следующих запусках.

После загрузки файлов откройте в GitHub вкладку **Actions**. У workflow **CI** должны появиться два результата матрицы: `Build and Test (Go 1.26.x)` и `Build and Test (Go 1.27.x)`. Дождитесь зелёных галочек у обоих.

Чтобы посмотреть логи, откройте запуск workflow → соответствующий job → шаги `Build`, `Run tests with race detector and coverage` или `Run golangci-lint`.

Чтобы показать, что линтер действительно блокирует ошибочный код (задание 2), после успешного первого запуска создайте отдельную ветку, а не меняйте `main`:

```powershell
git switch -c test/linter-check
@'
package main

func lintCheck() {
    unused := 123
}
'@ | Set-Content -Encoding utf8 lint_check.go

git add lint_check.go
git commit -m "Test lint pipeline"
git push -u origin test/linter-check
```

Затем откройте Pull Request из `test/linter-check` в `main`. CI должен обнаружить неиспользуемую переменную/ошибку компиляции и завершить проверку неуспешно. Не выполняйте merge этой версии. Удалите тестовый файл, отправьте исправление в ту же ветку и дождитесь зелёного CI:

```powershell
Remove-Item lint_check.go
git add -A
git commit -m "Remove intentional lint error"
git push
```

После этого Pull Request должен стать доступным для слияния, если остальные проверки пройдены.

Чтобы скачать покрытие: откройте успешный запуск CI, найдите блок **Artifacts** внизу страницы и скачайте `coverage-go-1.26.x` или `coverage-go-1.27.x`. После распаковки в терминале из каталога с файлом выполните:

```powershell
go tool cover -html=coverage.out -o coverage.html
```

Откройте `coverage.html` в браузере.

## 6. Docker: сборка и публикация образа

`Dockerfile` использует multi-stage build: сначала исходный код собирается в образе Go, затем готовый бинарник переносится в небольшой runtime-образ Alpine. В контейнере открыты порты `8080` (HTTP) и `9090` (gRPC). Процесс запускается не от root.

Workflow `.github/workflows/docker.yml` запускается при push в `main`, при push тега `v*` и при pull request. Для pull request образ только собирается, но не отправляется в Docker Hub. При push в `main` или тега `v*` workflow авторизуется в Docker Hub и публикует образ.

После успешного запуска откройте Docker Hub → `student-service` → **Tags**. Для обычного push в `main` ожидаются теги `main`, `latest` и `sha-...`. Для версии с Git-тегом `v1.0.0` появится Docker-тег `1.0.0`.

### Локальная проверка опубликованного образа

Замените `YOUR_DOCKERHUB_USERNAME` на логин Docker Hub:

```powershell
docker pull YOUR_DOCKERHUB_USERNAME/student-service:main
docker run --rm --name student-service -p 8080:8080 -p 9090:9090 YOUR_DOCKERHUB_USERNAME/student-service:main
```

Не закрывая окно с контейнером, откройте второй терминал PowerShell и проверьте HTTP API:

```powershell
$body = @{
  full_name = "Иванов Иван Иванович"
  group = "ПИЖ-б-о-23-2"
  gpa = 4.6
} | ConvertTo-Json

$created = Invoke-RestMethod -Method Post `
  -Uri "http://localhost:8080/students" `
  -ContentType "application/json" `
  -Body $body

$created
$id = $created.id
Invoke-RestMethod -Method Get -Uri "http://localhost:8080/students/$id"
```

Если ответ содержит ID и данные студента, HTTP API в контейнере работает. Остановите контейнер сочетанием `Ctrl+C` в первом терминале.

Если порт занят, остановите локально запущенный Go-сервис или измените порт слева в параметре `-p`, например `-p 8081:8080`, а запрос выполняйте на `localhost:8081`.

## 7. Защита ветки main

Настройте защиту после первого успешного запуска CI:

1. GitHub → репозиторий → **Settings → Branches**.
2. Выберите **Add branch protection rule**.
3. В `Branch name pattern` укажите `main`.
4. Включите **Require a pull request before merging**.
5. Включите **Require status checks to pass before merging**.
6. В списке проверок выберите оба job:
   - `Build and Test (Go 1.26.x)`;
   - `Build and Test (Go 1.27.x)`.
7. Сохраните правило.

Проверка: создайте тестовую ветку, внесите небольшое изменение README, отправьте ветку и откройте Pull Request в `main`. В PR должны запуститься проверки CI. Слияние станет доступно после успешного прохождения обязательных проверок.

На GitHub Free защита веток доступна для публичных репозиториев; для частных репозиториев правила защиты требуют платный план Pro или выше. Не делайте репозиторий публичным, если в нём есть секретные/личные данные.

## 8. Индивидуальный вариант 3 ЛР 7 — релизы по тегам

Workflow `.github/workflows/release.yml` запускается только при отправке тега, начинающегося с `v`.

После того как основной CI прошёл, создайте версионный тег:

```powershell
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

Откройте вкладку Actions и дождитесь workflow **Release**. Он собирает три бинарных файла под Linux, Windows и macOS (amd64), загружает их как артефакты, затем создаёт GitHub Release с этими файлами и контрольными суммами `SHA256SUMS.txt`.

Одновременно workflow Docker собирает образ для тега `v1.0.0` и публикует его с тегом `1.0.0`.

Для следующего релиза используйте новую версию, например `v1.0.1`. Не создавайте повторно тот же тег.

## 9. Проверка заданий и что сохранить для отчёта

Перед сдачей проверьте:

- проект ЛР 6 загружен в GitHub в ветку `main`;
- CI запускается на push и pull request;
- две версии Go в матрице прошли сборку, тесты и линтер;
- в Actions доступен artifact с `coverage.out`;
- в Docker Hub опубликован `student-service` с тегами;
- образ скачивается и запускается локально, HTTP-запрос успешно создаёт и читает студента;
- защита `main` блокирует merge до прохождения CI;
- при push `v1.0.0` появляется GitHub Release с тремя бинарными файлами.

Для отчёта достаточно сделать скриншоты: структуру файлов проекта; workflow CI с зелёными проверками двух версий Go; шаги Build/Test/Lint; скачанный artifact; workflow Docker и страницу Tags в Docker Hub; работающий контейнер и ответ HTTP; настройки защиты `main` и Pull Request; GitHub Release с файлами для трёх ОС.

## 10. Частые проблемы

**`go: go.mod requires go >= 1.26.0`** — локальная версия Go слишком старая. Установите Go 1.26.x/1.27.x, закройте и заново откройте терминал, проверьте `go version`.

**`remote: Invalid username or token`** — вместо пароля GitHub используйте авторизацию через браузер/Git Credential Manager или Personal Access Token. Не используйте пароль аккаунта.

**Docker workflow не проходит вход** — проверьте точные имена секретов `DOCKERHUB_USERNAME` и `DOCKERHUB_TOKEN`, а также разрешение Read & Write у токена Docker Hub.

**`denied: requested access to the resource is denied`** — убедитесь, что репозиторий `student-service` создан в Docker Hub под тем же аккаунтом, чей логин лежит в `DOCKERHUB_USERNAME`.

**`port is already allocated`** — порт 8080 или 9090 занят. Остановите другой сервис или выберите свободный внешний порт при `docker run`.

**В Actions не видно обязательной проверки для Branch protection** — сначала выполните хотя бы один push/PR, чтобы GitHub зарегистрировал job, затем вернитесь к настройке правила и выберите нужные status checks.
