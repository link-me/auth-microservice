package main

import (
    "context"
    "errors"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    jwt "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    redis "github.com/redis/go-redis/v9"
    "golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}

var (
    users       = map[string]string{} // email -> bcrypt hash (демо-память)
    jwtSecret   = []byte(getEnv("JWT_SECRET", "dev-secret"))
    accessTTL   = parseDuration(getEnv("ACCESS_TTL", "900s"))
    refreshTTL  = parseDuration(getEnv("REFRESH_TTL", "720h"))
    redisClient = mustRedis(getEnv("REDIS_URL", "redis://localhost:6379"))
)

func main() {
    r := gin.Default()

    r.POST("/register", handleRegister)
    r.POST("/login", handleLogin)
    r.POST("/refresh", handleRefresh)
    r.POST("/logout", handleLogout)
    r.GET("/me", handleMe)

    addr := getEnv("ADDR", ":8080")
    log.Printf("Auth microservice demo listening on %s", addr)
    if err := r.Run(addr); err != nil {
        log.Fatal(err)
    }
}

func handleRegister(c *gin.Context) {
    var req RegisterRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": "bad_request", "message": "invalid JSON"})
        return
    }
    if req.Email == "" || req.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"code": "bad_request", "message": "email and password required"})
        return
    }
    if _, exists := users[req.Email]; exists {
        c.JSON(http.StatusConflict, gin.H{"code": "conflict", "message": "user already exists"})
        return
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": "internal", "message": "hash error"})
        return
    }
    users[req.Email] = string(hash)
    c.JSON(http.StatusCreated, gin.H{"status": "registered"})
}

func handleLogin(c *gin.Context) {
    var req LoginRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"code": "bad_request", "message": "invalid JSON"})
        return
    }
    hash, ok := users[req.Email]
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"code": "unauthorized", "message": "invalid credentials"})
        return
    }
    if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"code": "unauthorized", "message": "invalid credentials"})
        return
    }

    access, err := generateAccessToken(req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": "internal", "message": "token error"})
        return
    }
    refresh, jti, err := generateRefreshToken(req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": "internal", "message": "token error"})
        return
    }

    if err := storeRefresh(c, req.Email, jti, refreshTTL); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": "internal", "message": "session store error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func handleMe(c *gin.Context) {
    auth := c.GetHeader("Authorization")
    if len(auth) < 8 || auth[:7] != "Bearer " {
        c.JSON(http.StatusUnauthorized, gin.H{"code": "unauthorized", "message": "bearer required"})
        return
    }
    tokenStr := auth[7:]
    claims, err := parseToken(tokenStr)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"code": "unauthorized", "message": "invalid token"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"email": claims.Subject, "exp": claims.ExpiresAt.Time.Unix()})
}

func handleRefresh(c *gin.Context) {
    var req RefreshRequest
    if err := c.BindJSON(&req); err != nil || req.RefreshToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"code": "bad_request", "message": "refresh_token required"})
        return
    }
    claims, err := parseToken(req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"code": "unauthorized", "message": "invalid refresh"})
        return
    }
    jti := claims.ID
    if jti == "" || claims.Subject == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"code": "unauthorized", "message": "invalid refresh"})
        return
    }
    // Проверяем, что refresh активен в Redis
    if ok, _ := isRefreshActive(c, claims.Subject, jti); !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"code": "unauthorized", "message": "refresh revoked"})
        return
    }
    // Ротация: отзываем старый, выдаём новый
    _ = revokeRefresh(c, claims.Subject, jti)

    access, err := generateAccessToken(claims.Subject)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": "internal", "message": "token error"})
        return
    }
    refresh, newJti, err := generateRefreshToken(claims.Subject)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": "internal", "message": "token error"})
        return
    }
    if err := storeRefresh(c, claims.Subject, newJti, refreshTTL); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": "internal", "message": "session store error"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func handleLogout(c *gin.Context) {
    var req RefreshRequest
    if err := c.BindJSON(&req); err != nil || req.RefreshToken == "" {
        c.JSON(http.StatusBadRequest, gin.H{"code": "bad_request", "message": "refresh_token required"})
        return
    }
    claims, err := parseToken(req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"code": "unauthorized", "message": "invalid refresh"})
        return
    }
    _ = revokeRefresh(c, claims.Subject, claims.ID)
    c.JSON(http.StatusOK, gin.H{"status": "logged_out"})
}

// ===== Helpers =====

func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func parseDuration(s string) time.Duration {
    d, err := time.ParseDuration(s)
    if err != nil {
        return 15 * time.Minute
    }
    return d
}

func mustRedis(url string) *redis.Client {
    opt, err := redis.ParseURL(url)
    if err != nil {
        log.Printf("redis url parse error: %v", err)
        // Падаем в демо-режим без Redis (nil), некоторые операции будут недоступны
        return nil
    }
    cli := redis.NewClient(opt)
    if err := cli.Ping(context.Background()).Err(); err != nil {
        log.Printf("redis ping error: %v", err)
        return nil
    }
    return cli
}

func generateAccessToken(sub string) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   sub,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTTL)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func generateRefreshToken(sub string) (string, string, error) {
    jti := uuid.NewString()
    claims := jwt.RegisteredClaims{
        Subject:   sub,
        ID:        jti,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTTL)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    str, err := token.SignedString(jwtSecret)
    return str, jti, err
}

func parseToken(tokenStr string) (*jwt.RegisteredClaims, error) {
    t, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return jwtSecret, nil
    })
    if err != nil || !t.Valid {
        return nil, errors.New("invalid token")
    }
    claims, ok := t.Claims.(*jwt.RegisteredClaims)
    if !ok {
        return nil, errors.New("invalid claims")
    }
    return claims, nil
}

func storeRefresh(c *gin.Context, sub, jti string, ttl time.Duration) error {
    if redisClient == nil {
        // демо-режим, без сохранения
        return nil
    }
    key := refreshKey(sub, jti)
    return redisClient.Set(c.Request.Context(), key, "1", ttl).Err()
}

func isRefreshActive(c *gin.Context, sub, jti string) (bool, error) {
    if redisClient == nil {
        // демо-режим: считаем активным
        return true, nil
    }
    key := refreshKey(sub, jti)
    res, err := redisClient.Exists(c.Request.Context(), key).Result()
    return res == 1, err
}

func revokeRefresh(c *gin.Context, sub, jti string) error {
    if redisClient == nil {
        return nil
    }
    key := refreshKey(sub, jti)
    return redisClient.Del(c.Request.Context(), key).Err()
}

func refreshKey(sub, jti string) string {
    return "session:" + sub + ":" + jti
}

import "fmt"
func main(){ fmt.Println("Demo start") }
// Improve performance
// Update dependencies
// Improve performance
// Add tests
// Update docs
// Enhance logging
// Improve performance
// Fix auth bug
// Setup CI
// Update docs
// Add feature
// Update dependencies
// Add feature
// Add feature
// Fix auth bug
// Fix auth bug
// Setup CI
// Enhance logging
// Improve performance
// Enhance logging
// Update docs
// Improve performance
// Setup CI
// Code cleanup
// Improve performance
// Update dependencies
// Code cleanup
// Setup CI
// Improve performance
// Add tests
// Update docs
// Setup CI
// Update dependencies
// Refactor module
// Refactor module
// Setup CI
// Add feature
// Enhance logging
// Setup CI
// Setup CI
// Add tests
// Refactor module
// Improve performance
// Update dependencies
// Enhance logging
// Update dependencies
// Improve performance
// Fix auth bug
// Update docs
// Improve performance
// Add tests
// Update dependencies
// Add tests
