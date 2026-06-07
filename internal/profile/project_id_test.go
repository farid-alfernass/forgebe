package profile

import "testing"

func TestProjectIDStableAndSafe(t *testing.T) {
	id1 := ProjectID("/Users/me/My Service", "git@example.com:org/repo.git")
	id2 := ProjectID("/Users/me/My Service", "git@example.com:org/repo.git")
	if id1 != id2 {
		t.Fatalf("expected stable id, got %q and %q", id1, id2)
	}
	if id1 == "" {
		t.Fatal("expected non-empty id")
	}
	for _, r := range id1 {
		if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != '-' && r != '_' && r != '.' {
			t.Fatalf("unsafe character %q in id %q", r, id1)
		}
	}
}

func TestProjectIDDifferentRemoteChangesHash(t *testing.T) {
	id1 := ProjectID("/repo", "remote-a")
	id2 := ProjectID("/repo", "remote-b")
	if id1 == id2 {
		t.Fatalf("expected different ids for different remotes")
	}
}
