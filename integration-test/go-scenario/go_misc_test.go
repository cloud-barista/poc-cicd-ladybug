package goscenario

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGoMisc(t *testing.T) {
	t.Run("go api misc test", func(t *testing.T) {
		SetUpForGrpc()

		// McarApi Set/Get Testing
		McarApi.SetServerAddr("127.0.0.1")
		res, _ := McarApi.GetServerAddr()
		assert.True(t, "127.0.0.1" == res)

		McarApi.SetTLSCA("./tls.ca")
		res, _ = McarApi.GetTLSCA()
		assert.True(t, "./tls.ca" == res)

		sec, _ := time.ParseDuration("10s")
		McarApi.SetTimeout(sec)
		tt, _ := McarApi.GetTimeout()
		assert.True(t, sec == tt)

		McarApi.SetJWTToken("abcdefg")
		res, _ = McarApi.GetJWTToken()
		assert.True(t, "abcdefg" == res)

		McarApi.SetInType("json")
		res, _ = McarApi.GetInType()
		assert.True(t, "json" == res)

		McarApi.SetInType("yaml")
		res, _ = McarApi.GetInType()
		assert.True(t, "yaml" == res)

		err := McarApi.SetInType("text")
		assert.True(t, err != nil)

		McarApi.SetOutType("json")
		res, _ = McarApi.GetOutType()
		assert.True(t, "json" == res)

		McarApi.SetOutType("yaml")
		res, _ = McarApi.GetOutType()
		assert.True(t, "yaml" == res)

		err = McarApi.SetOutType("text")
		assert.True(t, err != nil)

		TearDownForGrpc()
	})
}
