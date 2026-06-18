package objects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListGetMultiDimensional(t *testing.T) {
	type testCase struct {
		name      string
		list      *List
		indices   []Object
		expected  Object
		expectErr error
	}

	tests := []testCase{
		{
			name: "simple list get",
			list: &List{
				Elements: []Object{
					&Integer{Value: 1},
					&Integer{Value: 2},
					&Integer{Value: 3},
				},
			},
			indices:   []Object{&Integer{Value: 1}},
			expected:  &Integer{Value: 2},
			expectErr: nil,
		},
		{
			name: "nested list get",
			list: &List{
				Elements: []Object{
					&Integer{Value: 1},
					&List{
						Elements: []Object{
							&Integer{Value: 2},
							&Integer{Value: 3},
						},
					},
					&Integer{Value: 4},
				},
			},
			indices:   []Object{&Integer{Value: 1}, &Integer{Value: 0}},
			expected:  &Integer{Value: 2},
			expectErr: nil,
		},
		{
			name: "list with map get",
			list: &List{
				Elements: []Object{
					&Integer{Value: 1},
					&Map{
						Pairs: map[HashKey]HashPair{
							(&String{Value: "key"}).HashKey(): {
								Key:   &String{Value: "key"},
								Value: &Integer{Value: 2},
							},
						},
					},
					&Integer{Value: 3},
				},
			},
			indices:   []Object{&Integer{Value: 1}, &String{Value: "key"}},
			expected:  &Integer{Value: 2},
			expectErr: nil,
		},
		{
			name: "list with map and list get",
			list: &List{
				Elements: []Object{
					&Integer{Value: 1},
					&Map{
						Pairs: map[HashKey]HashPair{
							(&String{Value: "key"}).HashKey(): {
								Key: &String{Value: "key"},
								Value: &List{
									Elements: []Object{
										&Integer{Value: 2},
										&Integer{Value: 3},
									},
								},
							},
						},
					},
					&Integer{Value: 4},
				},
			},
			indices:   []Object{&Integer{Value: 1}, &String{Value: "key"}, &Integer{Value: 0}},
			expected:  &Integer{Value: 2},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.list.GetMultiDimensional(tt.indices)
			assert.Equal(t, tt.expectErr, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
