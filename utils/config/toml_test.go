package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()
	r := `
[service]
env = "development"

[log]

[mysql]
url        = "db"
dbName     = "famili-api"
username   = "famili-api"
password   = "password"
`

	// 一時ファイル生成
	tmp, err := ioutil.TempFile("", "toml-")
	assert.NoError(t, err)
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(r); err != nil {
		t.Fatal(err)
	}

	c, err := NewConfig(tmp.Name())
	assert.NoError(t, err)

	cases := []struct {
		value string
		want  string
	}{
		{"env", "development"},
	}

	for _, v := range cases {
		got := c.Service.Env
		if got != v.want {
			t.Errorf("%s: want %s got %s", v.value, v.want, got)
		}
	}
}
