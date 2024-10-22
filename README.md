# Система постов и комментариев к ним

## Описание

Этот проект представляет собой систему постов и комментариев, аналогичную системам, используемым на таких платформах, как Хабр или Reddit. Пользователи могут создавать посты, просматривать список постов, а также добавлять комментарии. Комментарии могут быть вложенными, без ограничения на глубину вложенности. Автор поста может запретить комментарии к своему посту.

## Разворот сервисов (локально)

Для локальной разработки используется docker compose.

Поднятие всей инфраструктуры осуществляется с помощью команды `make up`.

По завершении работы возможно остановить все сервисы командой `make down`.

База данных Postgres доступна внутри docker сети по адресу `postgres:5432`,
а также в локальной сети по адресу `0.0.0.0:5432`.

## Форматирование и стиль кода.

Для запуска форматирования используется команда `make fmt`.

Для запуска проверки линтера используется команда `make lint`.

## Тестирование.

Для запуска тестов используется команда `make test`.


## GraphQL API

После запуска всех сервисов GraphQL API будет доступен по адресу `http://localhost:8080`.

### Примеры запросов:

- Получение списка постов:
```graphql
query Post($id: ID!, $parentCommentID: ID, $cursor: String, $limit: Int) {
  post(id: $id, parentCommentId: $parentCommentID, cursor: $cursor, limit: $limit) {
    id
    title
    content
    commentsDisabled
    comments(cursor: $cursor, limit: $limit) {
      edges {
        cursor
        node {
          id
          postId
          text
          parentCommentId
          createdAt
          childCommentCount
        }
      }
      pageInfo {
        endCursor
        hasNextPage
      }
    }
  }
}


Variables (пример).
```json
{
  "cursor": "2",
  "limit": 5,
  "commentsLimit": null
}
