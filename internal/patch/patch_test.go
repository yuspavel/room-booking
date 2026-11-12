package patch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		payload     string
		expectedSet bool
		expectedNil bool
		expectedVal string
	}{
		{
			name:        "with value set",
			payload:     `{"name":"Pavel"}`,
			expectedSet: true,
			expectedNil: false,
			expectedVal: "Pavel",
		},
		{
			name:        "without value set",
			payload:     `{}`,
			expectedSet: false,
			expectedNil: true,
		},
		{
			name:        "with null value",
			payload:     `{"name":null}`,
			expectedSet: true,
			expectedNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var Req struct {
				Name Field[string] `json:"name"`
			}

			err := json.Unmarshal([]byte(tc.payload), &Req)

			require.NoError(t, err)
			require.Equal(t, tc.expectedSet, Req.Name.Set)
			if tc.expectedNil {
				require.Nil(t, Req.Name.Val)
			} else {
				require.NotNil(t, Req.Name.Val)
				require.Equal(t, tc.expectedVal, *Req.Name.Val)
			}
		})
	}
}
