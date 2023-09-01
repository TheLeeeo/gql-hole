package client

import (
	"testing"

	"github.com/TheLeeeo/gql-test-suite/models"
	"golang.org/x/exp/slices"
)

func Test_GetBaseType(t *testing.T) {
	tests := []struct {
		name string
		t    *models.Type
		want *models.Type
	}{
		{
			name: "Scalar",
			t: &models.Type{
				Kind: models.ScalarTypeKind,
			},
			want: &models.Type{
				Kind: models.ScalarTypeKind,
			},
		},
		{
			name: "Enum",
			t: &models.Type{
				Kind: models.EnumTypeKind,
			},
			want: &models.Type{
				Kind: models.EnumTypeKind,
			},
		},
		{
			name: "Object",
			t: &models.Type{
				Kind: models.ObjectTypeKind,
			},
			want: &models.Type{
				Kind: models.ObjectTypeKind,
			},
		},
		{
			name: "List",
			t: &models.Type{
				Kind: models.ListTypeKind,
				OfType: &models.Type{
					Kind: models.ScalarTypeKind,
				},
			},
			want: &models.Type{
				Kind: models.ScalarTypeKind,
			},
		},
		{
			name: "NonNull",
			t: &models.Type{
				Kind: models.NonNullTypeKind,
				OfType: &models.Type{
					Kind: models.ScalarTypeKind,
				},
			},

			want: &models.Type{
				Kind: models.ScalarTypeKind,
			},
		},
		{
			name: "NonNullList",
			t: &models.Type{
				Kind: models.NonNullTypeKind,
				OfType: &models.Type{
					Kind: models.ListTypeKind,
					OfType: &models.Type{
						Kind: models.ScalarTypeKind,
					},
				},
			},

			want: &models.Type{
				Kind: models.ScalarTypeKind,
			},
		},
		{
			name: "NonNullObject",
			t: &models.Type{
				Kind: models.NonNullTypeKind,
				OfType: &models.Type{
					Kind: models.ObjectTypeKind,
				},
			},

			want: &models.Type{
				Kind: models.ObjectTypeKind,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.GetBaseType(); got.Kind != tt.want.Kind {
				t.Errorf("GetBaseType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CompileFields(t *testing.T) {
	tests := []struct {
		name string
		t    *models.Type
		want []string
	}{
		{
			name: "Scalar",
			t: &models.Type{
				Kind: models.ScalarTypeKind,
			},
			want: []string{"__typename"},
		},
		{
			name: "Enum",
			t: &models.Type{
				Kind: models.EnumTypeKind,
			},
			want: []string{"__typename"},
		},
		{
			name: "Object",
			t: &models.Type{
				Kind: models.ObjectTypeKind,
				Fields: []models.Field{
					{
						Name: "field1",
						Type: &models.Type{
							Kind: models.ScalarTypeKind,
						},
					},
					{
						Name: "field2",
						Type: &models.Type{
							Kind: models.ScalarTypeKind,
						},
					},
				},
			},
			want: []string{"__typename"},
		},
		{
			name: "List",
			t: &models.Type{
				Kind: models.ListTypeKind,
				OfType: &models.Type{
					Kind: models.ObjectTypeKind,
					Fields: []models.Field{
						{
							Name: "field1",
							Type: &models.Type{
								Kind: models.ScalarTypeKind,
							},
						},
						{
							Name: "field2",
							Type: &models.Type{
								Kind: models.ScalarTypeKind,
							},
						},
					},
				},
			},
			want: []string{},
		},
		{
			name: "NonNull",
			t: &models.Type{
				Kind: models.NonNullTypeKind,
				OfType: &models.Type{
					Kind: models.ObjectTypeKind,
					Fields: []models.Field{
						{
							Name: "field1",
							Type: &models.Type{
								Kind: models.ScalarTypeKind,
							},
						},
						{
							Name: "field2",
							Type: &models.Type{
								Kind: models.ScalarTypeKind,
							},
						},
					},
				},
			},
			want: []string{},
		},
		{
			name: "NonNullList",
			t: &models.Type{
				Kind: models.NonNullTypeKind,
				OfType: &models.Type{
					Kind: models.ListTypeKind,
					OfType: &models.Type{
						Kind: models.ObjectTypeKind,
						Fields: []models.Field{

							{
								Name: "field1",
								Type: &models.Type{
									Kind: models.ScalarTypeKind,
								},
							},
							{
								Name: "field2",
								Type: &models.Type{
									Kind: models.ScalarTypeKind,
								},
							},
						},
					},
				},
			},
			want: []string{},
		},
	}

	cl := New(&Config{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cl.CompileType(tt.t); !slices.Equal(got, tt.want) {
				t.Errorf("CompileFields(), got %v, want %v", got, tt.want)
			}
		})
	}
}
