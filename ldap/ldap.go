package ldap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	ldap "github.com/go-ldap/ldap"
	"github.com/parnurzeal/gorequest"
)

// LDAP 实例
type LDAP struct {
	Host     string
	BaseDN   string
	Filter   string
	Attr     string
	BindDN   string
	Password string
	USERS    string
}

// NewLDAP NewLDAP
func NewLDAP(ld *LDAP) (err error) {

	ld.USERS = os.Getenv("LDAP_USERS")
	if ld.USERS != "" {
		return
	}

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

var fileProtocol = "file://"
var fileProtocolLength = len(fileProtocol)

func getUsersFromURL(url string) (users []string, err error) {
	var body []byte
	if strings.HasPrefix(url, fileProtocol) {
		body, err = ioutil.ReadFile(url[fileProtocolLength:])
	} else {
		_, bodyStr, errs := gorequest.New().Get(url).End()
		if len(errs) > 0 {
			err = errs[0]
			return
		}
		body = []byte(bodyStr)
	}
	err = json.Unmarshal([]byte(body), &users)
	return
}

// GetUsers by filter
func (ld LDAP) GetUsers() (users []string, err error) {
	if ld.USERS != "" {
		return getUsersFromURL(ld.USERS)
	}
	l, err := ld.bind()
	if err != nil {
		return
	}
	defer l.Close()

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
