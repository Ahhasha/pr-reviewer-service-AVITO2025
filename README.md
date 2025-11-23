PR Reviewer Service

Сервис для автоматического назначения ревьюверов на пул-реквесты.

Что делает:
- Создает команды разработчиков
- При создании PR автоматически назначает до 2 ревьюверов из команды автора
- Позволяет переназначать ревьюверов
- Управляет активностью пользователей

API методы:

1. Создать команду
POST /team/add
{
  "team_name": "backend",
  "members": [
    {"user_id": "user1", "username": "Alice", "is_active": true},
    {"user_id": "user2", "username": "Bob", "is_active": true}
  ]
}

2. Получить команду
GET /team/get?team_name=backend

3. Изменить активность пользователя
POST /users/setIsActive
{
  "user_id": "user2", 
  "is_active": false
}

4. Создать PR (автоназначение ревьюверов)
POST /pullRequest/create
{
  "pull_request_id": "pr-1",
  "pull_request_name": "New feature",
  "author_id": "user1"
}

5. Переназначить ревьювера
POST /pullRequest/reassign
{
  "pull_request_id": "pr-1",
  "old_user_id": "user2"
}

6. Замержить PR
POST /pullRequest/merge
{
  "pull_request_id": "pr-1"
}

7. Получить PR пользователя
GET /users/getReview?user_id=user2

Как запустить:

1. Установите Docker
2. Выполните: docker-compose up --build
3. Сервис будет на http://localhost:8080

Проверка здоровья:
curl http://localhost:8080/health

Пример работы:

1. Создаем команду
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "dev",
    "members": [
      {"user_id": "alice", "username": "Alice", "is_active": true},
      {"user_id": "bob", "username": "Bob", "is_active": true},
      {"user_id": "charlie", "username": "Charlie", "is_active": true}
    ]
  }'

2. Создаем PR - автоматически назначатся 2 ревьювера
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "pr-1", 
    "pull_request_name": "Feature",
    "author_id": "alice"
  }'

Особенности:

- Назначаются только активные пользователи (is_active = true)
- Автор PR исключается из списка ревьюверов
- Можно назначить 0, 1 или 2 ревьювера в зависимости от доступности
- После мерджа PR нельзя менять ревьюверов
- Переназначение ищет замену из команды автора

Ошибки:

409 Conflict - когда нельзя выполнить операцию (PR уже мерджен, нет кандидатов)
404 Not Found - когда ресурс не найден
400 Bad Request - невалидные данные

Для тестирования используйте postman_collection.json

Технические детали:

- Go 1.24
- PostgreSQL 
- Docker
- Миграции применяются автоматически при docker-compose up