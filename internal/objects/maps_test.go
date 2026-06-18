package objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapGetMultiDimensional(t *testing.T) {
	type testCase struct {
		name      string
		m         *Map
		indices   []Object
		expected  Object
		expectErr error
	}

	tests := []testCase{
		{
			name: "simple map get",
			m: &Map{
				Pairs: map[HashKey]HashPair{
					(&String{Value: "key"}).HashKey(): {
						Key:   &String{Value: "key"},
						Value: &Integer{Value: 1},
					},
				},
			},
			indices:   []Object{&String{Value: "key"}},
			expected:  &Integer{Value: 1},
			expectErr: nil,
		},
		{
			name: "map with list get",
			m: &Map{
				Pairs: map[HashKey]HashPair{
					(&String{Value: "key"}).HashKey(): {
						Key: &String{Value: "key"},
						Value: &List{
							Elements: []Object{
								&Integer{Value: 1},
								&Integer{Value: 2},
							},
						},
					},
				},
			},
			indices:   []Object{&String{Value: "key"}, &Integer{Value: 1}},
			expected:  &Integer{Value: 2},
			expectErr: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.m.GetMultiDimensional(tc.indices)
			assert.Equal(t, tc.expected, result)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}
