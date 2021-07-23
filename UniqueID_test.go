package acccore

import "testing"

func TestRandomGenUniqueIDGenerator_NewUniqueID(t *testing.T) {
	testMap := make(map[string]bool)
	gen := &RandomGenUniqueIDGenerator{
		Length:     16,
		LowerAlpha: false,
		UpperAlpha: true,
		Numeric:    true,
	}
	iterations := 10000000
	count := 0
	for count <= iterations {
		count++
		id := gen.NewUniqueID()
		if _, exist := testMap[id]; exist {
			t.Errorf("ID conflict on %d iteration", count)
			t.FailNow()
		}
		testMap[id] = true
	}
}
