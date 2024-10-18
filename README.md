# Auth Microservice

![status](https://img.shields.io/badge/status-stable-green)
![since](https://img.shields.io/badge/since-2024-blue)
![license](https://img.shields.io/badge/license-MIT-lightgrey)

Stack: Go (Gin) + Redis

## Overview
- Демонстрационный сервис аутентификации с базовыми функциями: регистрация, логин, выдача `access`/`refresh` токенов (JWT HS256), ротация `refresh`, `logout`, `GET /me`.
- Хранение `refresh`-токенов/сессий в Redis по ключу `session:{email}:{jti}` с TTL.
- Конфигурация через переменные окружения, простое локальное развёртывание.

## Features
- Регистрация пользователя (in-memory хранилище для демо).
- Логин с проверкой пароля (`bcrypt`).
- JWT `access` (короткий TTL) и `refresh` (длинный TTL) с ротацией.
- `GET /me` по `Authorization: Bearer <access>`.
- `logout` с отзывом `refresh`-токена.

## Quick Start (Windows PowerShell)
1) Установите переменные окружения:
   - `$env:JWT_SECRET = "dev-secret"`
   - `$env:ACCESS_TTL = "900s"`  (15 минут)
   - `$env:REFRESH_TTL = "720h"` (30 дней)
   - `$env:REDIS_URL = "redis://localhost:6379"`
2) Установите зависимости и запустите:
   - `go mod tidy`
   - `go run ./src`
3) Быстрый запуск:
   - `./run.ps1` (скрипт выставит env и запустит сервер)

## Endpoints
- `POST /register` `{email,password}`
- `POST /login` `{email,password}` → `{access_token, refresh_token}`
- `GET  /me` `Authorization: Bearer <access_token>` → `{email,exp}`
- `POST /refresh` `{refresh_token}` → новая пара `{access_token, refresh_token}` (ротация)
- `POST /logout` `{refresh_token}` → отзыв refresh

## Configuration
- `JWT_SECRET` — секрет для подписи HS256 (default: `dev-secret`)
- `ACCESS_TTL` — TTL для access (default: `900s`)
- `REFRESH_TTL` — TTL для refresh (default: `720h`)
- `REDIS_URL` — подключение к Redis (default: `redis://localhost:6379`)
- `ADDR` — адрес сервера (default: `:8080`)

## How it works
- При логине генерируется `access` и `refresh` (claims: `sub=email`, `jti` для refresh).
- `refresh` фиксируется в Redis по ключу `session:{email}:{jti}` с TTL.
- `refresh` при обновлении ротируется: старый отзывается, новый записывается.
- `me` читает `access` и возвращает сведения из claims.

## Examples (PowerShell)
- Регистрация:
  `Invoke-RestMethod -Uri http://localhost:8080/register -Method Post -Body (@{email='a@b.c';password='Str0ng!'} | ConvertTo-Json) -ContentType 'application/json'`
- Логин:
  `$login = Invoke-RestMethod -Uri http://localhost:8080/login -Method Post -Body (@{email='a@b.c';password='Str0ng!'} | ConvertTo-Json) -ContentType 'application/json'`
  `$access = $login.access_token; $refresh = $login.refresh_token`
- Me:
  `Invoke-RestMethod -Uri http://localhost:8080/me -Headers @{Authorization="Bearer $access"}`
- Refresh:
  `Invoke-RestMethod -Uri http://localhost:8080/refresh -Method Post -Body (@{refresh_token=$refresh} | ConvertTo-Json) -ContentType 'application/json'`
- Logout:
  `Invoke-RestMethod -Uri http://localhost:8080/logout -Method Post -Body (@{refresh_token=$refresh} | ConvertTo-Json) -ContentType 'application/json'`

## Production Notes
- Секреты только из env/secret storage, не хардкодить.
- Короткий TTL для `access`, длинный для `refresh`.
- Логирование и аудит входов/выходов; корреляция по `request_id`.
- Ограничение попыток `/login` (интеграция с rate limiter).
- Ротация секретов; отзыв refresh при компрометации.

## Metrics (набросок)
- `auth_logins_total{status}`
- `auth_refresh_total{status}`
- `auth_logout_total{status}`
- `auth_token_latency_ms{type}`

## Roadmap
- RU: см. `projects/auth-microservice/ROADMAP.md`
- EN: see `projects/auth-microservice/ROADMAP.en.md`

## Language
- Private usage notes exist in two languages and are local-only:
  - RU: `projects/auth-microservice/PRIVATE_USAGE.txt`
  - EN: `projects/auth-microservice/PRIVATE_USAGE.en.txt`
  - Note: both files are ignored by Git.

## Releases

Смотрите раздел Tags в репозитории и `CHANGELOG.md`.
Ключевые версии:
- `v0.1.0` — первый демо‑релиз (апрель 2024)
- `v0.2.0` — CI и базовый пайплайн (май 2024)
- `v0.3.0` — базовые тесты (июнь 2024)
- `v0.4.0` — обновление документации (июль 2024)
- `v0.5.0` — обновления зависимостей (август 2024)
- `v0.6.0` — расширенный README, утилиты запуска (октябрь 2024)
- `v0.6.1` — исправление сборки, tidy (октябрь 2024)
