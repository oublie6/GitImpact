// jwt_test.go 验证 JWT 中间件在合法与缺失令牌下的行为。
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	pkgAuth "gitimpact/backend/pkg/auth"

	"github.com/gin-gonic/gin"
)

func TestJWTMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", JWT("secret"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	tok, err := pkgAuth.GenerateToken("secret", time.Hour, 1, "u", "user")
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
