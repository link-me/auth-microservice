# Roadmap

Goals for Auth Microservice (Go Gin + Redis)

Short-term (1–2 weeks)
- Core endpoints: register, login, refresh, logout, `GET /me`.
- Tokens: JWT (access/refresh), secure signing (HS256/ES256), refresh rotation.
- Password hashing: `argon2id` (or `bcrypt`), salt, OWASP‑aligned parameters.
- Sessions: store refresh tokens/sessions in Redis (blacklist/allowlist).
- Configuration: env (`JWT_SECRET`, `ACCESS_TTL`, `REFRESH_TTL`, `REDIS_URL`), validation.
- Input validation: email/password, normalization, login rate‑limit.
- Tests: unit (hashing/JWT), integration (endpoints, Redis), Postman collection.
- Code quality: `golangci-lint`, `gofmt`, `go vet`, pre‑commit hooks.
- CI: GitHub Actions (Go 1.x), linters, tests, module cache.

Mid-term (3–6 weeks)
- RBAC: roles and permissions, access‑check middleware.
- OAuth2/OpenID Connect: Google/GitHub providers, account linking.
- Notifications: email verification, password reset, rate limiting.
- Audit: login/logout/events logging, `request_id` correlation.
- Admin API: user block, force logout, role management.
- Documentation: OpenAPI/Swagger, request examples, token schemas.
- Containerization: Dockerfile, compose (Gin + Redis + [DB]).

Long-term (6+ weeks)
- Multi‑tenant: tenant boundaries and data isolation.
- SSO/SAML: enterprise provider integration.
- Security policies: 2FA/TOTP/WebAuthn, password requirements, deactivation.
- Scaling: stateless access tokens, horizontal refresh rotation, Redis Cluster.
- Observability: metrics (success/error, latency), tracing (OTEL), alerts.

Technical notes
- JWT: short TTL for access, longer for refresh; rotate on every refresh.
- Redis keys: `session:{user_id}:{refresh_id}` with TTL; store fingerprint/ua/ip.
- Errors: unified JSON with codes and reasons; avoid leaking details.
- Security: CSRF not relevant for pure API, but enforce CORS; limit body size.
- Config: separate dev/prod; secrets only from env/secret storage.