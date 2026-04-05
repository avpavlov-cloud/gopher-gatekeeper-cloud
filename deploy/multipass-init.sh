#!/bin/bash
# Создаем ноду для Keycloak
multipass launch --name auth-node --cpus 2 --memory 2G
# Создаем ноду для API
multipass launch --name api-node --cpus 1 --memory 1G

# Установка Docker на auth-node (упрощенно)
multipass exec auth-node -- sudo apt-get update
multipass exec auth-node -- sudo apt-get install -y docker-compose
