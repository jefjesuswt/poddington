package fleet

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/jefjesuswt/walroos/config"
	_ "modernc.org/sqlite"
)

func setupTestDb(t *testing.T) *config.Database {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	db, err := config.NewSQLite(dbPath)
	if err != nil {
		t.Fatalf("Could not start DB in RAM")
	}

	return db
}

func TestFleetSQLiteStore(t *testing.T) {
	db := setupTestDb(t)
	defer db.Close()

	store := NewSQLiteStore(db)
	ctx := context.Background()

	t.Run("Should do migrations", func(t *testing.T) {
		err := store.Migrate(ctx)
		if err != nil {
			t.Fatalf("Could not migrate database: %v", err)
		}
	})

	testNode := Node{
		ID:        "test-node",
		Name:      "Test Node",
		Address:   "127.0.0.1:8080",
		Token:     "test-token",
		CreatedAt: time.Now().UTC(),
		LastSeen:  time.Now().UTC().Add(-1 * time.Hour),
	}

	t.Run("Should save", func(t *testing.T) {
		err := store.Save(ctx, testNode)
		if err != nil {
			t.Fatalf("Could not save node: %v", err)
		}
	})

	t.Run("Should return ErrNodeNameAlreadyExists when duplicated", func(t *testing.T) {
		nodeWithSameName := testNode
		nodeWithSameName.ID = "test-node-2"
		nodeWithSameName.Address = "127.0.0.1:8081"

		err := store.Save(ctx, nodeWithSameName)
		if err != ErrNodeNameAlreadyExists {
			t.Fatalf("Expected ErrNodeNameAlreadyExists, got: %v", err)
		}
	})

	t.Run("Should return ErrNodeAddressAlreadyExists when duplicated", func(t *testing.T) {
		nodeWithSameAddress := testNode
		nodeWithSameAddress.ID = "test-node-3"
		nodeWithSameAddress.Name = "Test Node 3"

		err := store.Save(ctx, nodeWithSameAddress)
		if err != ErrNodeAddressAlreadyExists {
			t.Fatalf("Expected ErrNodeNameAlreadyExists, got: %v", err)
		}
	})

	t.Run("Should return all nodes", func(t *testing.T) {
		_, err := store.GetAll(ctx)
		if err != nil {
			t.Fatalf("Could not get nodes: %v", err)
		}
	})

	t.Run("Should return node by ID", func(t *testing.T) {
		node, err := store.GetByID(ctx, testNode.ID)

		if err != nil {
			t.Fatalf("Could not get node by ID: %v", err)
		}

		if node.Name != testNode.Name {
			t.Fatalf("Expected name: %s, got %s", testNode.Name, node.Name)
		}
	})

	t.Run("Should update node", func(t *testing.T) {
		err := store.UpdateLastSeen(ctx, testNode.ID)

		if err != nil {
			t.Fatalf("Could not update last seen: %v", err)
		}

		updatedNode, err := store.GetByID(ctx, testNode.ID)
		if err != nil {
			t.Fatalf("GetByID failed on updated node: %v", err)
		}

		if !updatedNode.LastSeen.After(testNode.LastSeen) {
			t.Fatalf("didn't update last_seen, origina: %s, got: %s", testNode.LastSeen, updatedNode.LastSeen)
		}
	})

	t.Run("Should not update a non existant node", func(t *testing.T) {
		err := store.UpdateLastSeen(ctx, "non-existant")

		if err != ErrNodeNotFound {
			t.Fatalf("Expected ErrNodeNotFound, got: %v", err)
		}
	})

	t.Run("Should detete node", func(t *testing.T) {
		err := store.Delete(ctx, testNode.ID)

		if err != nil {
			t.Fatalf("Could not delete node: %v", err)
		}

		x, err := store.GetByID(ctx, testNode.ID)
		if err != ErrNodeNotFound {
			t.Fatalf("Expected ErrNodeNotFound, got: %v", err)
		}

		if x != (Node{}) {
			t.Fatalf("Expected empty node, got: %v", x)
		}
	})
}
