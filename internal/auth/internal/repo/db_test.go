package repo

import (
	"auth/internal/models"
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func BenchmarkCreateToken(b *testing.B) {
	multipliers := []int{10, 16, 32, 64}

	for _, m := range multipliers {
		b.Run(fmt.Sprintf("Count of simbols %d", m), func(b *testing.B) {
			_ = createToken(m)
		})

	}
}

func TestDataBase_CreateSession(t *testing.T) {

	userID := uuid.New()

	repo := &DataBase{}
	testSession := models.Session{
		UserID: userID,
	}

	s, errCode := repo.CreateSession(testSession)
	if errCode != models.OK {
		t.Errorf("Error occured")
	}

	if s.Token != testSession.Token {
		t.Errorf("Error occured")
	}

	if s.Token != "" {
		t.Errorf("Error occured")
	}

}
