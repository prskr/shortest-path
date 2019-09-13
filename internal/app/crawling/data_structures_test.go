package crawling

import "testing"

func Test_stringSet_Add(t *testing.T) {
	type fields struct {
		items map[string]bool
	}
	type args struct {
		value string
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantAlreadyPresent bool
		resultingSetSize   int
	}{
		{
			name: "Test add into empty set",
			fields: fields{
				items: make(map[string]bool),
			},
			args: args{
				value: "hello, world",
			},
			wantAlreadyPresent: false,
			resultingSetSize:   1,
		},
		{
			name: "Test add already existing string",
			fields: fields{
				items: map[string]bool{
					"hello, world": true,
				},
			},
			args: args{
				value: "hello, world",
			},
			wantAlreadyPresent: true,
			resultingSetSize:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := &stringSet{
				items: tt.fields.items,
			}
			if gotAlreadyPresent := set.Add(tt.args.value); gotAlreadyPresent != tt.wantAlreadyPresent {
				t.Errorf("Add() = %v, want %v", gotAlreadyPresent, tt.wantAlreadyPresent)
			}

			if tt.resultingSetSize != len(set.items) {
				t.Errorf("Expected resulting set to have size %d but got %d", tt.resultingSetSize, len(set.items))
			}
		})
	}
}

func Test_stringSet_Contains(t *testing.T) {
	type fields struct {
		items map[string]bool
	}
	type args struct {
		value string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		wantContained bool
	}{
		{
			name: "check for value in empty set",
			fields: fields{
				items: make(map[string]bool),
			},
			args: args{
				value: "hello, world",
			},
			wantContained: false,
		},
		{
			name: "check for actually contained value",
			fields: fields{
				items: map[string]bool{
					"hello, world": true,
				},
			},
			args: args{
				value: "hello, world",
			},
			wantContained: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := &stringSet{
				items: tt.fields.items,
			}
			if gotContained := set.Contains(tt.args.value); gotContained != tt.wantContained {
				t.Errorf("Contains() = %v, want %v", gotContained, tt.wantContained)
			}
		})
	}
}
