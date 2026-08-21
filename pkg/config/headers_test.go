package config

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyPreservedHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		preserve []string
		src      http.Header
		want     http.Header
	}{
		{
			name:     "empty preserve list is a no-op",
			preserve: nil,
			src: http.Header{
				"Authorization": []string{"Bearer incoming"},
				"X-Custom":      []string{"custom"},
			},
			want: http.Header{},
		},
		{
			name:     "forwards Authorization when listed",
			preserve: []string{"Authorization"},
			src: http.Header{
				"Authorization": []string{"Bearer incoming"},
				"X-Custom":      []string{"custom"},
			},
			want: http.Header{
				"Authorization": []string{"Bearer incoming"},
			},
		},
		{
			name:     "forwards custom headers when listed",
			preserve: []string{"X-Scope-OrgID", "X-Custom"},
			src: http.Header{
				"X-Scope-Orgid": []string{"tenant-a"},
				"X-Custom":      []string{"custom"},
				"Authorization": []string{"Bearer incoming"},
			},
			want: http.Header{
				"X-Scope-Orgid": []string{"tenant-a"},
				"X-Custom":      []string{"custom"},
			},
		},
		{
			name:     "blocks headers that are not listed",
			preserve: []string{"Authorization"},
			src: http.Header{
				"Authorization": []string{"Bearer incoming"},
				"Cookie":        []string{"secret"},
				"X-Custom":      []string{"custom"},
			},
			want: http.Header{
				"Authorization": []string{"Bearer incoming"},
			},
		},
		{
			name:     "matches listed names case-insensitively",
			preserve: []string{"authorization", "x-scope-orgid"},
			src: http.Header{
				"Authorization": []string{"Bearer incoming"},
				"X-Scope-Orgid": []string{"tenant-a"},
				"X-Custom":      []string{"custom"},
			},
			want: http.Header{
				"Authorization": []string{"Bearer incoming"},
				"X-Scope-Orgid": []string{"tenant-a"},
			},
		},
		{
			name:     "copies all values of a multi-value header",
			preserve: []string{"X-Custom"},
			src: http.Header{
				"X-Custom": []string{"one", "two"},
			},
			want: http.Header{
				"X-Custom": []string{"one", "two"},
			},
		},
		{
			name:     "ignores blank preserve names",
			preserve: []string{"", "Authorization"},
			src: http.Header{
				"Authorization": []string{"Bearer incoming"},
			},
			want: http.Header{
				"Authorization": []string{"Bearer incoming"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dst := make(http.Header)
			CopyPreservedHeaders(dst, tt.src, tt.preserve)
			require.Equal(t, tt.want, dst)
		})
	}
}

func TestCopyPreservedHeaders_DoesNotShareBackingArray(t *testing.T) {
	t.Parallel()

	src := http.Header{"X-Custom": []string{"original"}}
	dst := make(http.Header)
	CopyPreservedHeaders(dst, src, []string{"X-Custom"})

	src["X-Custom"][0] = "mutated"
	require.Equal(t, "original", dst.Get("X-Custom"))
}

func TestCopyPreservedHeaders_NilSafe(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		CopyPreservedHeaders(nil, http.Header{"Authorization": []string{"Bearer x"}}, []string{"Authorization"})
		CopyPreservedHeaders(make(http.Header), nil, []string{"Authorization"})
	})
}
