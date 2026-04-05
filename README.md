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
   
Получение токена
```bash
~$ curl -X POST "http://10.157.62.188:8080/realms/my-project/protocol/openid-connect/token" \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "client_id=notes-api" \
     -d "client_secret=hz5HnBIakHeAkRqLNOCtVvNKRn2R7Ps2" \
     -d "username=testuser" \
     -d "password=12345" \
     -d "grant_type=password"
```

Проверка полученного токена
```bash
 curl -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI2bkxqQXhTS21lR0pNbklxTzMxTlBsWFNpQVQ2YzdNQWlaMVpXa1JxNk5vIn0.eyJleHAiOjE3NzUzODQ3NzksImlhdCI6MTc3NTM4NDQ3OSwianRpIjoiYjNiN2E3YjYtODc5ZC00Y2YzLTlkMjQtODg5Zjg4ZGI5ZjUwIiwiaXNzIjoiaHR0cDovLzEwLjE1Ny42Mi4xODg6ODA4MC9yZWFsbXMvbXktcHJvamVjdCIsImF1ZCI6ImFjY291bnQiLCJzdWIiOiJkMjUwZDc3ZC02YTQzLTRhMTUtYjUxZi0yMjYzYzI5OWNhNmMiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJub3Rlcy1hcGkiLCJzZXNzaW9uX3N0YXRlIjoiNzZhMjZkYjktMzQ1Zi00Y2FlLTk1MDMtYjcxYWRiZTUxMzIwIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1teS1wcm9qZWN0Iiwib2ZmbGluZV9hY2Nlc3MiLCJhZG1pbiIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJlbWFpbCBwcm9maWxlIiwic2lkIjoiNzZhMjZkYjktMzQ1Zi00Y2FlLTk1MDMtYjcxYWRiZTUxMzIwIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJ0ZXN0dXNlciIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiJ9.oYvvDE2BkJqyh3lLfJItNbB-TsRqmkXSehlbfMOUpx9PLCIstbQZWGb6VQayo-MhreVqrA5Z_h4S6_AdyqeUWO_ZOSSX0VwN8jnIr9GexjNdkbYEtRyJjiQQzTe06cimUGMpPq1I--0xg2wGI6QuA7yNeeaN0o_SAuBBjZRMG5ruusl-yAIKZrpqFkOSnQSiHpEcmnB_CZ8Xsqa7JM7T2IuFBH1pEhWQtzKSUY4wSfoF043K1nJzm2GeOW7h5yQVPfDSaZSFHuSq93JMKIANrHc3rPc72Vio_hFbnQyKk6uTElwI1kgdR6ZbZmwOAk20mWBeHVq24u0v4VtWFQdVgg" http://10.157.62.188:8080
 ```