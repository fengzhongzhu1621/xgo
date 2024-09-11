package set

import (
	"testing"

	gset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
)

func TestSetsDifference(t *testing.T) {
	tests := []struct {
		former  []string
		latter  []string
		want1   []interface{}
		want2   []interface{}
		wantErr bool
	}{
		{
			former:  []string{},
			latter:  []string{},
			want1:   []interface{}{},
			want2:   []interface{}{},
			wantErr: false,
		},
		{
			former:  []string{"a", "b", "c"},
			latter:  []string{"c", "d", "e"},
			want1:   []interface{}{"a", "b"},
			want2:   []interface{}{"d", "e"},
			wantErr: false,
		},
		{
			former:  []string{"a", "b", "c"},
			latter:  []string{"c", "b", "a"},
			want1:   []interface{}{},
			want2:   []interface{}{},
			wantErr: false,
		},
		{
			former:  []string{"a", "b", "c"},
			latter:  []string{"d", "e", "f"},
			want1:   []interface{}{"a", "b", "c"},
			want2:   []interface{}{"d", "e", "f"},
			wantErr: false,
		},
		{
			former:  []string{"a", "b", "c"},
			latter:  []string{},
			want1:   []interface{}{"a", "b", "c"},
			want2:   []interface{}{},
			wantErr: false,
		},
		{
			former:  []string{},
			latter:  []string{"a", "b", "c"},
			want1:   []interface{}{},
			want2:   []interface{}{"a", "b", "c"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		got1, got2, err := SetsDifference(tt.former, tt.latter)
		if (err != nil) != tt.wantErr {
			t.Errorf("SetsDifference(%v, %v) error = %v, wantErr %v", tt.former, tt.latter, err, tt.wantErr)
			continue
		}
		s1 := gset.NewSetFromSlice(got1)
		s2 := gset.NewSetFromSlice(got2)
		assert.True(t, s1.Equal(gset.NewSetFromSlice(tt.want1)))
		assert.True(t, s2.Equal(gset.NewSetFromSlice(tt.want2)))
	}
}
