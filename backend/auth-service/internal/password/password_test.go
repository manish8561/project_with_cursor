package password

import (
	"strings"
	"testing"
)

func TestHashAndCheckSuccess(t *testing.T) {
	hash, err := Hash("secret-password")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}
	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("expected bcrypt prefix, got %q", hash)
	}

	ok, needsRehash := Check(hash, "secret-password")
	if !ok {
		t.Fatal("expected password to match")
	}
	if needsRehash {
		t.Fatal("bcrypt hash should not need rehash")
	}
}

func TestCheckWrongPassword(t *testing.T) {
	hash, err := Hash("secret-password")
	if err != nil {
		t.Fatalf("Hash: %v", err)
	}

	ok, needsRehash := Check(hash, "wrong-password")
	if ok {
		t.Fatal("expected password mismatch")
	}
	if needsRehash {
		t.Fatal("failed bcrypt check should not request rehash")
	}
}

func TestCheckLegacyPlaintextNeedsRehash(t *testing.T) {
	ok, needsRehash := Check("legacy-plain", "legacy-plain")
	if !ok {
		t.Fatal("expected legacy plaintext to match")
	}
	if !needsRehash {
		t.Fatal("legacy plaintext match should need rehash")
	}
}

func TestCheckLegacyPlaintextMismatch(t *testing.T) {
	ok, needsRehash := Check("legacy-plain", "other")
	if ok {
		t.Fatal("expected legacy plaintext mismatch")
	}
	if needsRehash {
		t.Fatal("legacy mismatch should not request rehash")
	}
}
