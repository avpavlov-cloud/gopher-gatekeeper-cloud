package auth

import (
	"context"
	"fmt"
	"net/http"
	"notes-api/internal/config"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
)

// KeycloakMiddleware проверяет валидность JWT токена
func KeycloakMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := config.LoadConfig()

		// 1. Извлекаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Требуется Bearer токен", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Инициализируем провайдер OIDC (Keycloak)
		ctx := context.Background()
		provider, err := oidc.NewProvider(ctx, fmt.Sprintf("%s/realms/%s", cfg.KeycloakURL, cfg.Realm))
		if err != nil {
			http.Error(w, "Ошибка подключения к Keycloak", http.StatusInternalServerError)
			return
		}

		// 3. Верифицируем подпись токена
		verifier := provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})
		_, err = verifier.Verify(ctx, tokenString)
		if err != nil {
			http.Error(w, "Невалидный токен: "+err.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AdminOnlyMiddleware проверяет наличие роли 'admin' в токене
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		claims, ok := token.Claims.(jwt.MapClaims)

		// Логика Keycloak: роли лежат в realm_access -> roles
		if ok {
			if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
				if roles, ok := realmAccess["roles"].([]interface{}); ok {
					for _, role := range roles {
						if role == "admin" {
							next.ServeHTTP(w, r)
							return
						}
					}
				}
			}
		}

		http.Error(w, "Доступ запрещен: требуется роль admin", http.StatusForbidden)
	})
}
