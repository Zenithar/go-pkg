package runtime

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestStatsCollector(t *testing.T) {
	s := &Stats{}
	s.Collect()

	out, _ := json.Marshal(s)
	fmt.Println(string(out))
}
