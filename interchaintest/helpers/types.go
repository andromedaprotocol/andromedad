package helpers

import "testing"

// EntryPoint
type QueryMsg struct {
	// IBCHooks
	GetCount      *GetCountQuery      `json:"get_count,omitempty"`
	GetTotalFunds *GetTotalFundsQuery `json:"get_total_funds,omitempty"`
}

func debugOutput(t *testing.T, stdout string) {
	if true {
		t.Log(stdout)
	}
}
