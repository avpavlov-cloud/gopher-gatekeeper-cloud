## Для инициализации multipass
Нужно запустить содержимое ``multipass-init.sh`` и просмотреть IP
```bash
 multipass list
 ```

Монтирование проекта к серверу api-node
```bash
multipass mount . api-node:/home/ubuntu/app
```

## Запуск приложения
```bash
 multipass shell api-node
 sudo snap install go --classic
 cd app
 go run cmd/api/main.go
```

Проверка сервиса
```bash
curl -i http://10.157.62.174:8000/health
```

## Запуск аудентификации
```bash
multipass mount deploy/keycloak auth-node:/home/ubuntu/keycloak-deploy
multipass shell auth-node
cd ~/keycloak-deploy
sudo docker-compose up -d
```

Админка Keycloack находится по адресу ``http://10.157.62.188:8080``. В неё заходим с пользователем `admin` и паролем `admin`

 
## Настройка Keycloak

   1. Создание Realm: В меню Master -> Create Realm -> Name: my-project.
   2. Создание Client: Clients -> Create client -> ID: notes-api.
   * В Capability Config включите Client authentication (ON) -> Save.
      * Во вкладке Credentials скопируйте Client Secret.
   3. Создание Роли: Realm roles -> Create role -> Name: admin.
   4. Создание Пользователя: Users -> Add user -> Username: testuser.
   * Вкладка Credentials: Set password (напр. 12345), выключите Temporary.
      * Важно: Во вкладке Details очистите поле Required user actions, если там что-то есть.
      * Вкладка Role mapping: Assign role -> выберите admin.

### Настройка Keycloak для Go-сервиса (audience "notes-api")

Настройка Client Scopes

1. В админке Keycloak перейди **Client Scopes** → **Create client scope**  
   - Name: `notes-api-scope`  
   - Protocol: `OpenID Connect`  
2. Сохрани.  
3. В созданном scope перейди **Mappers** → **Configure a new mapper** → выбери **Audience**  
   - Name: `notes-api-audience`  
   - Included Client Audience: выбери `notes-api`  
4. Сохрани.

Привязка scope к клиенту

1. Clients → `notes-api` → вкладка **Client Scopes** → **Add client scope**  
2. Выбери `notes-api-scope` → тип **Default** → **Add**

Проверка клиента

- Client authentication: ON  
- Authorization: OFF  
- Standard flow: ON  
- Direct access grants: ON  

После настройки

1. Получи новый токен (старый не работает).  
2. Проверить токен на [jwt.io](https://jwt.io) → в Payload должно быть:  

```json
"aud": "notes-api"
```

## Получение токена
```bash
~$ curl -X POST "http://10.157.62.188:8080/realms/my-project/protocol/openid-connect/token" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "client_id=notes-api" \
     -d "client_secret=hz5HnBIakHeAkRqLNOCtVvNKRn2R7Ps2" \
     -d "username=testuser" \
     -d "password=12345" \
     -d "grant_type=password"
```

Отлично! Теперь токен у тебя корректный — у него есть нужный audience (`"aud":["notes-api","account"]`) и scope для твоего API. Ниже — рабочие запросы к твоему `notes-api` с этим токеном.

---

### 🔹 Сохраняем токен в переменную

```bash
TOKEN="eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI2bkxqQXhTS21lR0pNbklxTzMxTlBsWFNpQVQ2YzdNQWlaMVpXa1JxNk5vIn0.eyJleHAiOjE3NzUzODg4NTAsImlhdCI6MTc3NTM4ODU1MCwianRpIjoiMmU5ZjBmOTktOWEyNy00ODgzLTk4YmUtYzk3NTE5ODlkMzdlIiwiaXNzIjoiaHR0cDovLzEwLjE1Ny42Mi4xODg6ODA4MC9yZWFsbXMvbXktcHJvamVjdCIsImF1ZCI6WyJub3Rlcy1hcGkiLCJhY2NvdW50Il0sInN1YiI6ImQyNTBkNzdkLTZhNDMtNGExNS1iNTFmLTIyNjNjMjk5Y2E2YyIsInR5cCI6IkJlYXJlciIsImF6cCI6Im5vdGVzLWFwaSIsInNlc3Npb25fc3RhdGUiOiIyN2ZhN2YyOC03NmNiLTRhMzMtYWFhZC04YTYzZmQ3NTVjYjgiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbIi8qIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJkZWZhdWx0LXJvbGVzLW15LXByb2plY3QiLCJvZmZsaW5lX2FjY2VzcyIsImFkbWluIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6ImVtYWlsIHByb2ZpbGUgbm90ZXMtYXBpLXNjb3BlIiwic2lkIjoiMjdmYTdmMjgtNzZjYi00YTMzLWFhYWQtOGE2M2ZkNzU1Y2I4IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJ0ZXN0dXNlciIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiJ9.BJjKgEwMQL2XXnU79XGlLAckB5k-IHtEOYAa4e4RIrWFPZPdQMpn0xwoWBtV_jk12Hz3NiLKeHQoRdQEnbNMCvKONc49DzqFmWHCyp0LKobNGgEWja0I4LYxREfNa8mxQMTO4IHdySFhZVnYS_JwODlWqg8BA4OrcB7ywE30NH696y64CBu5J908ZGXtwtDNa4Wv6sAk_rS91W4eNY4Ku3wgjqjoDsMlAeF_ueIkDfX__I8FRj51LVVHfo7vA4aheHgN-f6dr6hLxCCAjGmSQEn71ET6m3G7IMSBqR2kTz6uF9G5B3UUTT1Xhob55n-aPJMamQiMHbbmxvs93hTWqQ"
```

---

### 1️⃣ Проверка `/health`

```bash
curl -H "Authorization: Bearer $TOKEN" http://10.157.62.174:8000/health
```

Ожидаемый ответ:

```
Service is up
```

---

### 2️⃣ GET всех заметок `/api/v1/notes`

```bash
curl -H "Authorization: Bearer $TOKEN" http://10.157.62.174:8000/api/v1/notes/
```

### 3️⃣ POST новой заметки

```bash
curl -X POST http://10.157.62.174:8000/api/v1/notes/ \
 -H "Content-Type: application/json" \
 -H "Authorization: Bearer $TOKEN" \
 -d '{"title":"New Note","content":"Content of the note"}'
```

### 4️⃣ DELETE заметки по `id`

Допустим, `id` = `1`:

```bash
curl -X DELETE http://10.157.62.174:8000/api/v1/notes/1 \
 -H "Authorization: Bearer $TOKEN"
```

