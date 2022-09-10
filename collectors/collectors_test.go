package collectors

import "testing"

func TestGetComputersInfo(t *testing.T) {
	wsNames := []string{"ws555a", "ws836"}
	_ = GetComputersInfo(wsNames)
}
