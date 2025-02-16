#!/bin/bash

set -exuo pipefail

PROJECT_DIR="/home/jwl/avito-shop-service"
CONFIG_FILE="$PROJECT_DIR/.golangci.yml"

# Проверяем, существует ли конфиг-файл
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Ошибка: Конфигурационный файл $CONFIG_FILE не найден!"
    exit 1
fi

# Переходим в директорию проекта
cd "$PROJECT_DIR"

# Запускаем линтер для всех Go-файлов в проекте
golangci-lint run -c "$CONFIG_FILE" ./...
