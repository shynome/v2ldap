package v2ray

import (
	"testing"

	"github.com/shynome/v2ldap/model"
)

func TestConfig(t *testing.T) {
	v2 := &V2ray{}
	config := v2.GenConfig([]model.User{})
	if len(config.Inbound) != 1 {
		t.Errorf(`expect inbound length 1, get %v`, len(config.Inbound))
	}
	return
}
