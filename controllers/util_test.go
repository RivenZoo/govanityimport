package controllers

import (
	"reflect"
	"testing"
)

func Test_assignNonEmptyByFieldName(t *testing.T) {
	type testStruct struct {
		Num   int
		Name  string
		array []int
		ptr   *int
	}
	var testInt int
	type args struct {
		src    interface{}
		dst    interface{}
		fields []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		struct {
			name string
			args args
		}{name: "cover not non empty", args: args{
			src: &testStruct{
				Num:   1,
				Name:  "abc",
				array: []int{2},
				ptr:   &testInt,
			},
			dst: &testStruct{
				Num: 10,
			},
			fields: []string{"Num", "Name", "array", "ptr"},
		}},
		struct {
			name string
			args args
		}{name: "not modify with empty field", args: args{
			src: &testStruct{
				Num:   0,
				Name:  "",
				array: []int{2},
				ptr:   &testInt,
			},
			dst: &testStruct{
				Num:  10,
				Name: "111",
			},
			fields: []string{"Num", "Name", "array", "ptr"},
		}},
		struct {
			name string
			args args
		}{name: "struct", args: args{
			src: testStruct{
				Num:   0,
				Name:  "",
				array: []int{2},
				ptr:   &testInt,
			},
			dst: testStruct{
				Num:  10,
				Name: "111",
			},
			fields: []string{"Num", "Name", "array", "ptr"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var before, after testStruct
			dstVal := reflect.ValueOf(tt.args.dst)
			if dstVal.Kind() == reflect.Ptr {
				dst := dstVal.Elem().Interface()
				before = dst.(testStruct)
			} else {
				before = tt.args.dst.(testStruct)
			}

			assignNonEmptyByFieldName(tt.args.src, tt.args.dst, tt.args.fields...)

			dstVal = reflect.ValueOf(tt.args.dst)
			if dstVal.Kind() == reflect.Ptr {
				dst := dstVal.Elem().Interface()
				after = dst.(testStruct)
			} else {
				after = tt.args.dst.(testStruct)
			}

			var src testStruct
			srcVal := reflect.ValueOf(tt.args.src)
			if srcVal.Kind() == reflect.Ptr {
				srcInterface := srcVal.Elem().Interface()
				src = srcInterface.(testStruct)
			} else {
				src = tt.args.src.(testStruct)
			}

			t.Logf("src: %#v, dst before: %#v, after: %#v", src, before, after)
			if src.Num != 0 {
				if after.Num != src.Num {
					t.Fail()
				}
			}
			if src.Num == 0 {
				if after.Num != before.Num {
					t.Fail()
				}
			}
			if src.Name != "" {
				if after.Name != src.Name {
					t.Fail()
				}
			}
			if src.Name == "" {
				if after.Name != before.Name {
					t.Fail()
				}
			}
		})
	}
}

func Test_isZero(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		struct {
			name string
			args args
			want bool
		}{name: "int", args: args{
			v: reflect.ValueOf(int(0)),
		}, want: true},
		struct {
			name string
			args args
			want bool
		}{name: "string", args: args{
			v: reflect.ValueOf(""),
		}, want: true},
		struct {
			name string
			args args
			want bool
		}{name: "array", args: args{
			v: reflect.ValueOf([]int{0, 0}),
		}, want: true},
		struct {
			name string
			args args
			want bool
		}{name: "struct", args: args{
			v: reflect.ValueOf(struct {
				Str string
				N   int
			}{Str: "",
				N: 0,}),
		}, want: true},
		struct {
			name string
			args args
			want bool
		}{name: "non empty", args: args{
			v: reflect.ValueOf([]int{1}),
		}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isZero(tt.args.v); got != tt.want {
				t.Errorf("isZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
