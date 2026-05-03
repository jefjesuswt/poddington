package fleet

import (
	"context"
	"testing"
	"time"
)

type mockNodeStore struct {
	MigrateFn        func(ctx context.Context) error
	SaveFn           func(ctx context.Context, n Node) error
	GetAllFn         func(ctx context.Context) ([]Node, error)
	GetByIDFn        func(ctx context.Context, id string) (Node, error)
	DeleteFn         func(ctx context.Context, id string) error
	UpdateLastSeenFn func(ctx context.Context, id string) error
}

func (m *mockNodeStore) Migrate(ctx context.Context) error {
	if m.MigrateFn != nil {
		return m.MigrateFn(ctx)
	}
	panic("not implemented")
}

func (m *mockNodeStore) Save(ctx context.Context, n Node) error {
	if m.SaveFn != nil {
		return m.SaveFn(ctx, n)
	}
	panic("not implemented")
}

func (m *mockNodeStore) GetAll(ctx context.Context) ([]Node, error) {
	if m.GetAllFn != nil {
		return m.GetAllFn(ctx)
	}
	panic("not implemented")
}

func (m *mockNodeStore) GetByID(ctx context.Context, id string) (Node, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	panic("not implemented")
}

func (m *mockNodeStore) Delete(ctx context.Context, id string) error {
	if m.DeleteFn != nil {
		return m.DeleteFn(ctx, id)
	}
	panic("not implemented")
}

func (m *mockNodeStore) UpdateLastSeen(ctx context.Context, id string) error {
	if m.UpdateLastSeenFn != nil {
		return m.UpdateLastSeenFn(ctx, id)
	}
	panic("not implemented")
}

func TestService_Init(t *testing.T) {

	var calls int

	store := &mockNodeStore{
		MigrateFn: func(ctx context.Context) error {
			calls++
			return nil
		},
	}

	service := NewService(store)
	err := service.Init(context.Background())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestService_RegisterNode(t *testing.T) {
	var calls int
	var capturedNode Node

	store := &mockNodeStore{
		SaveFn: func(ctx context.Context, n Node) error {
			calls++
			capturedNode = n
			return nil
		},
	}

	service := NewService(store)

	nodeResult, err := service.RegisterNode(context.Background(), "test-node", "10.10.10.10")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}

	if capturedNode.Name != "test-node" || capturedNode.Address != "10.10.10.10" {
		t.Errorf("Expected node name and address to be set properly but got %+v", capturedNode)
	}

	if capturedNode.ID == "" || capturedNode.Token == "" {
		t.Errorf("Expected node ID and token to be set but got empty string")
	}

	if nodeResult.ID != capturedNode.ID {
		t.Errorf("Returned node does not equal to saved node")
	}
}

func TestService_ListNodes(t *testing.T) {
	var calls int

	store := &mockNodeStore{
		GetAllFn: func(ctx context.Context) ([]Node, error) {
			calls++
			return []Node{}, nil
		},
	}

	service := NewService(store)

	_, err := service.ListNodes(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}
}

func TestService_RemoveNode(t *testing.T) {
	var calls int

	store := &mockNodeStore{
		DeleteFn: func(ctx context.Context, id string) error {
			calls++
			return nil
		},
	}

	service := NewService(store)

	err := service.RemoveNode(context.Background(), "test-id")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if calls != 1 {
		t.Errorf("Expected 1 call, got %d", calls)
	}
}

func TestService_PingNode(t *testing.T) {

	t.Run("Should update last seen", func(t *testing.T) {
		var updateCalls int
		var getCalls int
		var capturedId string

		store := &mockNodeStore{
			GetByIDFn: func(ctx context.Context, id string) (Node, error) {
				getCalls++
				capturedId = id
				return Node{
					ID:        capturedId,
					Name:      "test",
					Address:   "10.10.10.10",
					Token:     "test-token",
					CreatedAt: time.Now().UTC(),
					LastSeen:  time.Now().UTC(),
				}, nil
			},
			UpdateLastSeenFn: func(ctx context.Context, id string) error {
				updateCalls++
				capturedId = id
				return nil
			},
		}

		service := NewService(store)

		err := service.PingNode(context.Background(), "test-id", "test-token")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updateCalls != 1 {
			t.Errorf("Expected 1 update call, got %d", updateCalls)
		}
		if getCalls != 1 {
			t.Errorf("Expected 1 get call, got %d", getCalls)
		}
		if capturedId != "test-id" {
			t.Errorf("Expected id to be set to test-id, got %s", capturedId)
		}
	})

	t.Run("Should not update last seen with invalid token", func(t *testing.T) {

		var updateCalls int
		var getCalls int
		var capturedId string

		store := &mockNodeStore{
			GetByIDFn: func(ctx context.Context, id string) (Node, error) {
				getCalls++
				capturedId = id
				return Node{
					ID:        capturedId,
					Name:      "test",
					Address:   "10.10.10.10",
					Token:     "test-token",
					CreatedAt: time.Now().UTC(),
					LastSeen:  time.Now().UTC(),
				}, nil
			},
			UpdateLastSeenFn: func(ctx context.Context, id string) error {
				updateCalls++
				capturedId = id
				return nil
			},
		}

		service := NewService(store)

		err := service.PingNode(context.Background(), "test-id", "not-test-token")

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if updateCalls != 0 {
			t.Errorf("Expected 0 update call, got %d", updateCalls)
		}
		if getCalls != 1 {
			t.Errorf("Expected 1 get call, got %d", getCalls)
		}
	})
}
