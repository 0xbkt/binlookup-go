package binlookup

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
)

const (
	CorrectBIN          = "5288230"
	CorrectButOrphanBIN = "9999999"
	IncorrectBIN        = "0812436"
)

func TestSearchWithCorrectBIN(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(10*time.Second))
	defer cancel()

expedition:
	_, err := Search(CorrectBIN)
	if err != nil {
		sce, ok := errors.Cause(err).(StatusCodeError)
		if ok && sce == http.StatusTooManyRequests {
			select {
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			default:
				goto expedition
			}
		}

		t.Fatalf("%+v", err)
	}
}

func TestSearchWithIncorrectBIN(t *testing.T) {
	_, err := Search(IncorrectBIN)
	if err == nil {
		t.Fatalf("%v is an incorrect BIN but Search returned nil error.", IncorrectBIN)
	}
}

func TestSearchWithCorrectButOrphanBIN(t *testing.T) {
	_, err := Search(CorrectButOrphanBIN)
	if err == nil {
		t.Fatalf("%v is an orphan BIN but Search returned nil for error.", CorrectButOrphanBIN)
	}
}

func TestStatusCodeError(t *testing.T) {
	ise := StatusCodeError(http.StatusInternalServerError)

	if ise.Error() != "500 Internal Server Error" {
		t.FailNow()
	}
}
