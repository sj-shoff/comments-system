# **Comments System** - система управления постами и комментариями с GraphQL API

---

## Описание
Система предоставляет следующие возможности:
- Создание/управление постами (включение/отключение комментариев)
- Добавление комментариев к постам (включая вложенные комментарии)
- Режим реального времени через WebSocket-подписки
- Поддержка двух режимов хранения: in-memory и PostgreSQL

---

## Технологический стек
- **Язык**: Go 1.24.3
- **GraphQL**: gqlgen
- **Базы данных**: PostgreSQL (основная), in-memory (для разработки)
- **Логирование**: slog с красивым выводом для разработки
- **Контейнеризация**: Docker
- **Миграции**: golang-migrate
- **Тестирование**: testify + mocks

---

# Запуск проекта

## Требования
- Go 1.24+
- Docker (для запуска через контейнеры)
- PostgreSQL (для production-режима)

## Обязательный пункт
Создать файл `.env` с паролем:
```env
POSTGRES_PASSWORD=your_password
```

## Запуск через Docker
```bash
# In-memory режим
make docker-inmemory

# PostgreSQL режим
make docker-postgres
```

---

# Структура проекта
```bash
.
├── cmd/
│   ├── comments-system/   # Основная точка входа приложения
│   └── migrator/          # Утилита для миграций
├── configs/
│   ├── inmemory.yaml      # Конфигурация для in-memory хранилища
│   └── postgres.yaml      # Конфигурация для PostgreSQL
├── docker/                # Файлы для Docker
├── internal/              # Внутренние модули приложения
│   ├── config/            # Конфигурация приложения
│   ├── graph/             # Реализация GraphQL
│   ├── models/            # Модели данных
│   ├── pubsub/            # Реализация pub/sub (публикация/подписка)
│   ├── service/           # Бизнес-логика сервиса
│   └── storage/           # Реализация хранилища данных
├── migrations/            # Файлы миграций базы данных (PostgreSQL)
├── pkg/                   # Общие пакеты, которые могут быть использованы в других проектах
│   ├── errors/            # Обработка ошибок
│   ├── logger/            # Логирование
│   └── utils/             # Вспомогательные утилиты
├── .env                   # Переменные окружения
├── .gitignore             # Игнорируемые файлы для Git
├── go.mod                 # Файл зависимостей Go
├── go.sum                 # Контрольные суммы для зависимостей
├── Makefile               # Утилита для сборки и управления проектом
└── README.md              # Описание проекта

```

---

# Примеры GraphQL-запросов

## Запросы (Queries)

### Получить список постов с пагинацией
```graphql
query GetPosts {
  posts(limit: 10, offset: 0) {
    id
    title
    author
    createdAt
  }
}
```

### Получить конкретный пост
```graphql
query GetPost {
  post(id: "1") {
    id
    title
    content
    commentsEnabled
  }
}
```

### Получить комментарии к посту
```graphql
query GetComments {
  comments(postId: "1", limit: 5, offset: 0) {
    total
    comments {
      id
      author
      content
    }
  }
}
```

### Получить ответы на комментарий
```graphql
query GetReplies {
  commentReplies(parentId: "comment_123") {
    id
    author
    content
  }
}
```

## Изменения (Mutations)

### Создать пост
```graphql
mutation CreatePost {
  createPost(input: {
    title: "Новый пост",
    content: "Содержание поста...",
    author: "Иван Иванов",
    commentsEnabled: true
  }) {
    id
    createdAt
  }
}
```

### Создать комментарий
```graphql
mutation CreateComment {
  createComment(input: {
    postId: "1",
    author: "Анна Петрова",
    content: "Отличный пост!"
  }) {
    id
    createdAt
  }
}
```

### Создать ответ на комментарий
```graphql
mutation CreateReply {
  createComment(input: {
    postId: "1",
    parentId: "comment_123",
    author: "Сергей",
    content: "Я согласен с предыдущим комментарием"
  }) {
    id
    content
  }
}
```

### Включить/отключить комментарии
```graphql
mutation ToggleComments {
  toggleComments(postId: "1", enabled: false) {
    id
    commentsEnabled
  }
}
```

## Подписки (Subscriptions)

### Подписаться на новые комментарии
```graphql
subscription OnCommentAdded {
  commentAdded(postId: "1") {
    id
    author
    content
    createdAt
  }
}
```

---

# Конфигурация
Inmemory хранилище:
```yaml
env: dev # local, dev, prod

server:
  port: "8080"

storage: "inmemory"
```

Postgres:
```yaml
env: dev # local, dev, prod
migrations: "./migrations"

server:
  port: "8080"

postgres:
  host: "postgres"
  port: "5432"
  username: "postgres"
  password: "" # из .env
  dbname: "comments"
  sslmode: "disable"

storage: "postgres"
```

---

# Тестирование
```bash
# Запуск всех тестов
make test
```

---

> Приложение доступно по адресу (GraphQL Playground): `http://localhost:8080`  