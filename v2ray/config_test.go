package v2ray

import (
	"testing"
)

func TestConfig(t *testing.T) {
	v2 := &V2ray{
		DB: getDB(),
	}
	config := v2.GetConfig()
	if len(config.Inbound) != 2 {
		t.Errorf(`expect inbound length 2, get %v`, len(config.Inbound))
	}
	return
}
