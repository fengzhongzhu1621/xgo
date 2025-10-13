package unmarshaler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYamlCodec_Unmarshal(t *testing.T) {
	t.Run("interface", func(t *testing.T) {
		var tt interface{}
		tt = map[string]interface{}{}
		require.Nil(t, GetCodec("yaml").Unmarshal([]byte("[1, 2]"), &tt))
	})
	t.Run("map[string]interface{}", func(t *testing.T) {
		tt := map[string]interface{}{}
		require.NotNil(t, GetCodec("yaml").Unmarshal([]byte("[1, 2]"), &tt))
	})
}
