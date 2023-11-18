package schema

import (
	"testing"
)

func Test_GetBaseType(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want Type
	}{
		{
			name: "Scalar",
			t: Type{
				Kind: ScalarTypeKind,
			},
			want: Type{
				Kind: ScalarTypeKind,
			},
		},
		{
			name: "Enum",
			t: Type{
				Kind: EnumTypeKind,
			},
			want: Type{
				Kind: EnumTypeKind,
			},
		},
		{
			name: "Object",
			t: Type{
				Kind: ObjectTypeKind,
			},
			want: Type{
				Kind: ObjectTypeKind,
			},
		},
		{
			name: "List",
			t: Type{
				Kind: ListTypeKind,
				OfType: &Type{
					Kind: ScalarTypeKind,
				},
			},
			want: Type{
				Kind: ScalarTypeKind,
			},
		},
		{
			name: "NonNull",
			t: Type{
				Kind: NonNullTypeKind,
				OfType: &Type{
					Kind: ScalarTypeKind,
				},
			},

			want: Type{
				Kind: ScalarTypeKind,
			},
		},
		{
			name: "NonNullList",
			t: Type{
				Kind: NonNullTypeKind,
				OfType: &Type{
					Kind: ListTypeKind,
					OfType: &Type{
						Kind: ScalarTypeKind,
					},
				},
			},

			want: Type{
				Kind: ScalarTypeKind,
			},
		},
		{
			name: "NonNullObject",
			t: Type{
				Kind: NonNullTypeKind,
				OfType: &Type{
					Kind: ObjectTypeKind,
				},
			},

			want: Type{
				Kind: ObjectTypeKind,
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
