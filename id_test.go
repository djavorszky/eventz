package eventz

import "testing"

func TestID(t *testing.T) {
	var id int
	curID = 1
	for i := 1; i < 100; i++ {
		id = ID()

		if id != i {
			t.Errorf("ID() = '%d', should be '%d'", id, i)
		}
	}
}
