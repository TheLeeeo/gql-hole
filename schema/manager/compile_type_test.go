package manager

import (
	"testing"

	"github.com/TheLeeeo/gql-test-suite/schema"
	"golang.org/x/exp/slices"
)

func Test_CompileFields(t *testing.T) {
	tests := []struct {
		name string
		t    schema.Type
		want []string
	}{
		{
			name: "Scalar",
			t: schema.Type{
				Kind: schema.ScalarTypeKind,
			},
			want: []string{"__typename"},
		},
		{
			name: "Enum",
			t: schema.Type{
				Kind: schema.EnumTypeKind,
			},
			want: []string{"__typename"},
		},
		{
			name: "Object",
			t: schema.Type{
				Kind: schema.ObjectTypeKind,
				Fields: []schema.Field{
					{
						Name: "field1",
						Type: &schema.Type{
							Kind: schema.ScalarTypeKind,
						},
					},
					{
						Name: "field2",
						Type: &schema.Type{
							Kind: schema.ScalarTypeKind,
						},
					},
				},
			},
			want: []string{"__typename"},
		},
		{
			name: "List",
			t: schema.Type{
				Kind: schema.ListTypeKind,
				OfType: &schema.Type{
					Kind: schema.ObjectTypeKind,
					Fields: []schema.Field{
						{
							Name: "field1",
							Type: &schema.Type{
								Kind: schema.ScalarTypeKind,
							},
						},
						{
							Name: "field2",
							Type: &schema.Type{
								Kind: schema.ScalarTypeKind,
							},
						},
					},
				},
			},
			want: []string{},
		},
		{
			name: "NonNull",
			t: schema.Type{
				Kind: schema.NonNullTypeKind,
				OfType: &schema.Type{
					Kind: schema.ObjectTypeKind,
					Fields: []schema.Field{
						{
							Name: "field1",
							Type: &schema.Type{
								Kind: schema.ScalarTypeKind,
							},
						},
						{
							Name: "field2",
							Type: &schema.Type{
								Kind: schema.ScalarTypeKind,
							},
						},
					},
				},
			},
			want: []string{},
		},
		{
			name: "NonNullList",
			t: schema.Type{
				Kind: schema.NonNullTypeKind,
				OfType: &schema.Type{
					Kind: schema.ListTypeKind,
					OfType: &schema.Type{
						Kind: schema.ObjectTypeKind,
						Fields: []schema.Field{

							{
								Name: "field1",
								Type: &schema.Type{
									Kind: schema.ScalarTypeKind,
								},
							},
							{
								Name: "field2",
								Type: &schema.Type{
									Kind: schema.ScalarTypeKind,
								},
							},
						},
					},
				},
			},
			want: []string{},
		},
	}

	m := New(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.CompileType(tt.t); !slices.Equal(got, tt.want) {
				t.Errorf("CompileFields(), got %v, want %v", got, tt.want)
			}
		})
	}
}
