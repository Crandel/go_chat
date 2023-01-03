package memory_test

import (
	"testing"

	m "github.com/Crandel/go_chat/internal/storage/memory"
)

func TestAddRoom(t *testing.T) {
	s := NewTestStorage()
	rName := "test room"
	gotName, err := s.AddRoom(rName)
	if err != nil {
		t.Errorf("AddRoom failed with error '%s'", err.Error())
	}

	if rName != gotName {
		t.Errorf("Got room name '%s' is different from expected '%s'", gotName, rName)
	}

	_, err = s.AddRoom("")
	if err == nil {
		t.Errorf("AddRoom with empty parameter failed without error")
	}
	if err.Error() != m.EmptyRoomNameError {
		t.Error("here")
	}
}
