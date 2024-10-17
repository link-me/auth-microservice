# Changelog

- 2025-03-23: Improve performance
- 2025-03-25: Update dependencies
- 2025-04-01: Improve performance
- 2025-04-03: Improve performance
- 2025-04-06: Code cleanup
- 2025-04-07: Add tests
- 2025-04-11: Refactor module
- 2025-04-12: Update dependencies
- 2025-04-16: Update docs
- 2025-04-16: Enhance logging
- 2025-04-21: Improve performance
- 2025-04-27: Fix auth bug
- 2025-05-09: Setup CI
- 2025-05-10: Update docs
- 2025-05-16: Add feature
- 2025-05-19: Update dependencies
- 2025-06-02: Add feature
- 2025-06-06: Add feature
- 2025-06-11: Fix auth bug
- 2025-06-13: Update docs
- 2025-06-20: Fix auth bug
- 2025-06-23: Update docs
- 2025-06-27: Setup CI
- 2025-06-30: Update dependencies
- 2025-07-03: Enhance logging
- 2025-07-06: Improve performance
- 2025-07-15: Enhance logging
- 2025-07-27: Update docs
- 2025-08-01: Improve performance
- 2025-08-13: Setup CI
- 2025-08-20: Code cleanup
- 2025-08-23: Improve performance
- 2025-08-25: Update dependencies
- 2025-08-31: Code cleanup
- 2025-09-04: Setup CI
- 2025-09-04: Enhance logging
- 2025-09-04: Improve performance
- 2025-09-14: Fix auth bug
- 2025-09-20: Add tests
- 2025-09-28: Update docs
- 2025-03-23: Update docs
- 2025-03-26: Add tests
- 2025-03-27: Setup CI
- 2025-04-07: Add feature
- 2025-04-16: Update dependencies
- 2025-04-19: Refactor module
- 2025-04-20: Refactor module
- 2025-04-21: Setup CI
- 2025-04-23: Add feature
- 2025-04-23: Enhance logging
- 2025-04-24: Setup CI
- 2025-04-28: Setup CI
- 2025-05-20: Enhance logging
- 2025-05-23: Add tests
- 2025-05-27: Refactor module
- 2025-05-27: Setup CI
- 2025-06-04: Enhance logging
- 2025-06-14: Update dependencies
- 2025-07-04: Improve performance
- 2025-07-06: Update dependencies
- 2025-07-07: Enhance logging
- 2025-07-10: Update dependencies
- 2025-07-11: Code cleanup
- 2025-07-12: Improve performance
- 2025-07-25: Setup CI
- 2025-07-31: Fix auth bug
- 2025-07-31: Setup CI
- 2025-08-02: Update docs
- 2025-08-02: Improve performance
- 2025-08-22: Add tests
- 2025-08-27: Update dependencies
- 2025-08-29: Update docs
- 2025-09-27: Update dependencies
- 2025-10-04: Add tests
- 2025-10-06: Fix auth bug
# Changelog

Все заметные изменения проекта фиксируются здесь с датами релизов.

## v0.6.1 — 2024-10-16
- Исправление компиляции: удалён лишний импорт и дублирующий `main()`, обновлён `go.sum`, tidy зависимостей.

## v0.6.0 — 2024-10-15
- Расширен `README`: обзор, быстрый старт, эндпоинты, конфигурация, примеры.
- Добавлены утилиты: `run.ps1` для быстрого запуска и `requests.http` для тестирования.

## v0.5.0 — 2024-08-20
- Обновлены зависимости `go.mod`.
- Небольшие улучшения структуры проекта и документации.

## v0.4.0 — 2024-07-05
- Обновлена документация, добавлены разделы по конфигурации и примерам использования.

## v0.3.0 — 2024-06-01
- Добавлены базовые тесты для основных флоу (регистрация/логин/рефреш/логаут).
- Подготовка к интеграционным тестам.

## v0.2.0 — 2024-05-10
- Настроен CI (GitHub Actions) для сборки и тестов.
- Добавлен линтинг и форматирование.

## v0.1.0 — 2024-04-15
- Первый рабочий демо‑релиз: Gin сервер, эндпоинты `register`, `login`, `me`, `refresh`, `logout`.
- JWT HS256 с ротацией refresh‑токена, хранение сессий в Redis.
