Монтирование проекта к серверу api-node
```bash
multipass mount . api-node:/home/ubuntu/app
```

Запуск приложения
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