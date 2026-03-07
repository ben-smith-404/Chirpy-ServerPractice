package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNoErrorInMakeJWT(t *testing.T) {
	_, resultErr := MakeJWT(uuid.Nil, "test", time.Minute)
	if resultErr != nil {
		t.Errorf("error has been generated in makeJWT: %v", resultErr)
	}
}

func TestJWTLifecycle(t *testing.T) {
	wantUUID := uuid.New()
	createdToken, _ := MakeJWT(wantUUID, "test", time.Minute)

	gotUUID, gotErr := ValidateJWT(createdToken, "test")

	if wantUUID != gotUUID {
		t.Errorf("token is not valid: %v", gotErr)
	}
}

func TestTimeoutToken(t *testing.T) {
	createdToken, _ := MakeJWT(uuid.New(), "test", time.Second)
	time.Sleep(time.Second)
	_, responseErr := ValidateJWT(createdToken, "test")

	if !strings.Contains(responseErr.Error(), "token is expired by") {
		t.Errorf("incorrect error: %v", responseErr)
	}
}
