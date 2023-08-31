package models

// func Test_GetBaseType(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		t    *Type
// 		want *Type
// 	}{
// 		{
// 			name: "Scalar",
// 			t: &Type{
// 				Kind: ScalarTypeKind,
// 			},
// 			want: &Type{
// 				Kind: ScalarTypeKind,
// 			},
// 		},
// 		{
// 			name: "Enum",
// 			t: &Type{
// 				Kind: EnumTypeKind,
// 			},
// 			want: &Type{
// 				Kind: EnumTypeKind,
// 			},
// 		},
// 		{
// 			name: "Object",
// 			t: &Type{
// 				Kind: ObjectTypeKind,
// 			},
// 			want: &Type{
// 				Kind: ObjectTypeKind,
// 			},
// 		},
// 		{
// 			name: "List",
// 			t: &Type{
// 				Kind: ListTypeKind,
// 				OfType: &Type{
// 					Kind: ScalarTypeKind,
// 				},
// 			},
// 			want: &Type{
// 				Kind: ScalarTypeKind,
// 			},
// 		},
// 		{
// 			name: "NonNull",
// 			t: &Type{
// 				Kind: NonNullTypeKind,
// 				OfType: &Type{
// 					Kind: ScalarTypeKind,
// 				},
// 			},

// 			want: &Type{
// 				Kind: ScalarTypeKind,
// 			},
// 		},
// 		{
// 			name: "NonNullList",
// 			t: &Type{
// 				Kind: NonNullTypeKind,
// 				OfType: &Type{
// 					Kind: ListTypeKind,
// 					OfType: &Type{
// 						Kind: ScalarTypeKind,
// 					},
// 				},
// 			},

// 			want: &Type{
// 				Kind: ScalarTypeKind,
// 			},
// 		},
// 		{
// 			name: "NonNullObject",
// 			t: &Type{
// 				Kind: NonNullTypeKind,
// 				OfType: &Type{
// 					Kind: ObjectTypeKind,
// 				},
// 			},

// 			want: &Type{
// 				Kind: ObjectTypeKind,
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.t.GetBaseType(); got.Kind != tt.want.Kind {
// 				t.Errorf("GetBaseType() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_CompileFields(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		t    *Type
// 		want []string
// 	}{
// 		{
// 			name: "Scalar",
// 			t: &Type{
// 				Kind: ScalarTypeKind,
// 			},
// 			want: []string{},
// 		},
// 		{
// 			name: "Enum",
// 			t: &Type{
// 				Kind: EnumTypeKind,
// 			},
// 			want: []string{},
// 		},
// 		{
// 			name: "Object",
// 			t: &Type{
// 				Kind: ObjectTypeKind,
// 				Fields: []Field{
// 					{
// 						Name: "field1",
// 						Type: &Type{
// 							Kind: ScalarTypeKind,
// 						},
// 					},
// 					{
// 						Name: "field2",
// 						Type: &Type{
// 							Kind: ScalarTypeKind,
// 						},
// 					},
// 				},
// 			},
// 			want: []string{"field1", "field2"},
// 		},
// 		{
// 			name: "List",
// 			t: &Type{
// 				Kind: ListTypeKind,
// 				OfType: &Type{
// 					Kind: ObjectTypeKind,
// 					Fields: []Field{
// 						{
// 							Name: "field1",
// 							Type: &Type{
// 								Kind: ScalarTypeKind,
// 							},
// 						},
// 						{
// 							Name: "field2",
// 							Type: &Type{
// 								Kind: ScalarTypeKind,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want: []string{},
// 		},
// 		{
// 			name: "NonNull",
// 			t: &Type{
// 				Kind: NonNullTypeKind,
// 				OfType: &Type{
// 					Kind: ObjectTypeKind,
// 					Fields: []Field{
// 						{
// 							Name: "field1",
// 							Type: &Type{
// 								Kind: ScalarTypeKind,
// 							},
// 						},
// 						{
// 							Name: "field2",
// 							Type: &Type{
// 								Kind: ScalarTypeKind,
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want: []string{},
// 		},
// 		{
// 			name: "NonNullList",
// 			t: &Type{
// 				Kind: NonNullTypeKind,
// 				OfType: &Type{
// 					Kind: ListTypeKind,
// 					OfType: &Type{
// 						Kind: ObjectTypeKind,
// 						Fields: []Field{

// 							{
// 								Name: "field1",
// 								Type: &Type{
// 									Kind: ScalarTypeKind,
// 								},
// 							},
// 							{
// 								Name: "field2",
// 								Type: &Type{
// 									Kind: ScalarTypeKind,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 			want: []string{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := tt.t.CompileFields(); !slices.Equal(got, tt.want) {
// 				t.Errorf("CompileFields() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
