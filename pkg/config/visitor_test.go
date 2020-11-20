package config

import (
	"testing"

	"gopkg.in/ini.v1"

	"github.com/fatedier/frp/pkg/consts"

	"github.com/stretchr/testify/assert"
)

const testVisitorPrefix = "test."

func Test_Visitor_UnmarshalFromIni(t *testing.T) {
	assert := assert.New(t)

	testcases := []struct {
		sname    string
		source   []byte
		expected VisitorConf
	}{
		{
			sname: "secret_tcp_visitor",
			source: []byte(`
				[secret_tcp_visitor]
				role = visitor
				type = stcp
				server_name = secret_tcp
				sk = abcdefg
				bind_addr = 127.0.0.1
				bind_port = 9000
				use_encryption = false
				use_compression = false
			`),
			expected: &STCPVisitorConf{
				BaseVisitorConf: BaseVisitorConf{
					ProxyName:  testVisitorPrefix + "secret_tcp_visitor",
					ProxyType:  consts.STCPProxy,
					Role:       "visitor",
					Sk:         "abcdefg",
					ServerName: testVisitorPrefix + "secret_tcp",
					BindAddr:   "127.0.0.1",
					BindPort:   9000,
				},
			},
		},
		{
			sname: "p2p_tcp_visitor",
			source: []byte(`
				[p2p_tcp_visitor]
				role = visitor
				type = xtcp
				server_name = p2p_tcp
				sk = abcdefg
				bind_addr = 127.0.0.1
				bind_port = 9001
				use_encryption = false
				use_compression = false
			`),
			expected: &XTCPVisitorConf{
				BaseVisitorConf: BaseVisitorConf{
					ProxyName:  testVisitorPrefix + "p2p_tcp_visitor",
					ProxyType:  consts.XTCPProxy,
					Role:       "visitor",
					Sk:         "abcdefg",
					ServerName: testProxyPrefix + "p2p_tcp",
					BindAddr:   "127.0.0.1",
					BindPort:   9001,
				},
			},
		},
	}

	for _, c := range testcases {
		f, err := ini.LoadSources(testLoadOptions, c.source)
		assert.NoError(err)

		visitorType := f.Section(c.sname).Key("type").String()
		assert.NotEmpty(visitorType)

		actual := DefaultVisitorConf(visitorType)
		assert.NotNil(actual)

		err = actual.UnmarshalFromIni(testVisitorPrefix, c.sname, f.Section(c.sname))
		assert.NoError(err)
		assert.Equal(c.expected, actual)
	}
}