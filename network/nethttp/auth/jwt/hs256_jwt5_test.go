package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHS256JWTV5Manager(t *testing.T) {
	manager := NewHS256JWTV5Manager("your-256-bit-secret", time.Hour)

	token, err := manager.Generate("Issuer", "user1", "admin")
	if err != nil {
		panic(err)
	}
	fmt.Println("Generated Token:", token)

	claims, err := manager.Verify(token)
	if err != nil {
		panic(err)
	}
	// &{UserID:user1 Role:admin StandardClaims:{Audience: ExpiresAt:1739414141 Id: IssuedAt:1739410541 Issuer: NotBefore:0 Subject:}}
	fmt.Printf("Verified Claims: %+v\n", claims)
	assert.Equal(t, "user1", claims.UserID)
}
