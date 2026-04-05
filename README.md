# 🛡️ Gopher Gatekeeper Cloud

Учебный проект по созданию защищённого API на Go с использованием Keycloak в качестве Identity Provider, развернутый в изолированной облачной среде на базе Multipass.

Проект демонстрирует:

* аутентификацию через OAuth2 / OIDC
* проверку JWT токенов
* ролевую модель доступа (RBAC)
* разделение инфраструктуры на auth и API ноды

---

## 🏗 Архитектура

Система состоит из двух виртуальных машин (Ubuntu 24.04):

| Нода          | IP              | Назначение                     |
| ------------- | --------------- | ------------------------------ |
| **auth-node** | `10.157.62.188` | Keycloak + PostgreSQL (Docker) |
| **api-node**  | `10.157.62.174` | Go API сервис                  |

### Поток работы:

1. Пользователь проходит аутентификацию в Keycloak
2. Получает JWT токен
3. Отправляет запрос в API
4. API валидирует токен и проверяет роли

---

## 🚀 Быстрый старт

### 1. Инициализация Multipass

```bash
bash multipass-init.sh
multipass list
```

---

### 2. Создание и подготовка нод

```bash
multipass launch --name auth-node --cpus 2 --memory 2G
multipass launch --name api-node --cpus 1 --memory 1G
```

Монтирование проекта:

```bash
multipass mount . api-node:/home/ubuntu/app
```

---

## 🔐 Запуск Keycloak (auth-node)

```bash
multipass mount deploy/keycloak auth-node:/home/ubuntu/keycloak-deploy
multipass shell auth-node
cd ~/keycloak-deploy
sudo docker-compose up -d
```

Админка доступна по адресу:
👉 [http://10.157.62.188:8080](http://10.157.62.188:8080)
Логин / пароль: `admin / admin`

---

## ⚙️ Настройка Keycloak

### 1. Realm

```
my-project
```

---

### 2. Client

* ID: `notes-api`
* Client authentication: ON
* Standard flow: ON
* Direct access grants: ON

Скопируйте **Client Secret**.

---

### 3. Роли

```
admin
```

---

### 4. Пользователь

* Username: `testuser`
* Password: `12345`
* Назначить роль: `admin`

---

### 5. Настройка Audience (ВАЖНО)

#### Создание Client Scope

* Name: `notes-api-scope`
* Protocol: OpenID Connect

#### Mapper

* Type: Audience
* Included Client Audience: `notes-api`

#### Привязка к клиенту

Clients → notes-api → Client Scopes → Add → Default

---

### ✔ Проверка токена

Откройте: JWT.io

В payload должно быть:

```json
"aud": "notes-api"
```

---

## 🔑 Получение токена

```bash
curl -X POST "http://10.157.62.188:8080/realms/my-project/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "client_id=notes-api" \
  -d "client_secret=YOUR_CLIENT_SECRET" \
  -d "username=testuser" \
  -d "password=12345" \
  -d "grant_type=password"
```

Сохраните токен:

```bash
TOKEN="your_jwt_token"
```

---

## 🚀 Запуск API (api-node)

```bash
multipass shell api-node
sudo snap install go --classic

cd app
go mod tidy
go run cmd/api/main.go
```

### Переменные окружения (.env)

```env
KEYCLOAK_URL=http://10.157.62.188:8080
KEYCLOAK_REALM=my-project
KEYCLOAK_CLIENT_ID=notes-api
PORT=8000
```

---

## 🛠 API Endpoints

| Метод  | Endpoint             | Доступ | Описание         |
| ------ | -------------------- | ------ | ---------------- |
| GET    | `/health`            | Public | Проверка сервиса |
| GET    | `/api/v1/notes`      | JWT    | Список заметок   |
| POST   | `/api/v1/notes`      | JWT    | Создание заметки |
| DELETE | `/api/v1/notes/{id}` | Admin  | Удаление заметки |

---

## 🧪 Тестирование API

### Проверка health

```bash
curl http://10.157.62.174:8000/health
```

---

### Получение заметок

```bash
curl -H "Authorization: Bearer $TOKEN" \
http://10.157.62.174:8000/api/v1/notes/
```

---

### Создание заметки

```bash
curl -X POST http://10.157.62.174:8000/api/v1/notes/ \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer $TOKEN" \
 -d '{"title":"New Note","content":"Content of the note"}'
```

---

### Удаление заметки

```bash
curl -X DELETE http://10.157.62.174:8000/api/v1/notes/1 \
 -H "Authorization: Bearer $TOKEN"
```

---

### Автотест

```bash
cd tests
chmod +x check_notes_api.sh
./check_notes_api.sh
```

Скрипт проверяет:

* получение токена
* публичный endpoint
* защищённые маршруты
* RBAC

---

## 🔒 Безопасность

* **OIDC / OAuth2** — проверка JWT через публичные ключи Keycloak
* **Audience Validation** — защита от чужих токенов
* **RBAC** — роли из `realm_access.roles`
* **Bearer Token** — защита всех приватных endpoint'ов

---

## 📦 Технологии

* Go
* Keycloak
* PostgreSQL
* Docker / Docker Compose
* Multipass

