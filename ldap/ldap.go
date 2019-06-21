package ldap

import (
	"fmt"
	"os"
	"reflect"

	ldap "github.com/go-ldap/ldap"
)

// LDAP 实例
type LDAP struct {
	Host     string
	BaseDN   string
	Filter   string
	Attr     string
	BindDN   string
	Password string
}

// NewLDAP NewLDAP
func NewLDAP(ld *LDAP) (err error) {

	v := reflect.ValueOf(ld)
	v = v.Elem()
	c := v.NumField()
	for i := 0; i < c; i++ {
		val := v.Field(i).Interface()
		if val != "" {
			continue
		}
		key := v.Type().Field(i).Name
		envName := fmt.Sprintf("LDAP_%s", key)
		envVal := os.Getenv(envName)
		if envVal == "" {
			return fmt.Errorf("Env: %s is required", envName)
		}
		v.Field(i).SetString(envVal)
	}
	return

}

func (ld LDAP) bind() (l *ldap.Conn, err error) {
	l, err = ldap.DialURL(ld.Host)
	if err != nil {
		return
	}
	if err = l.Bind(ld.BindDN, ld.Password); err != nil {
		return
	}
	return
}

// GetUsers by filter
func (ld LDAP) GetUsers() (users []string, err error) {
	l, err := ld.bind()
	defer l.Close()
	if err != nil {
		return
	}

	result, err := l.Search(&ldap.SearchRequest{
		BaseDN:     ld.BaseDN,
		Scope:      ldap.ScopeWholeSubtree,
		Filter:     ld.Filter,
		Attributes: []string{ld.Attr},
	})
	if err != nil {
		return
	}

	for _, entry := range result.Entries {
		username := entry.GetAttributeValue(ld.Attr)
		if username == "" {
			continue
		}
		users = append(users, username)
	}

	return
}
