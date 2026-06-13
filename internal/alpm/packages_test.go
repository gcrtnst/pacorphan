package alpm

import "testing"

func TestCompareVersion(t *testing.T) {
	tt := []struct {
		in1  string
		in2  string
		want int
	}{
		// all similar length, no pkgrel
		{"1.5.0", "1.5.0", 0},
		{"1.5.1", "1.5.0", 1},

		// mixed length
		{"1.5.1", "1.5", 1},

		// with pkgrel, simple
		{"1.5.0-1", "1.5.0-1", 0},
		{"1.5.0-1", "1.5.0-2", -1},
		{"1.5.0-1", "1.5.1-1", -1},
		{"1.5.0-2", "1.5.1-1", -1},

		// with pkgrel, mixed lengths
		{"1.5-1", "1.5.1-1", -1},
		{"1.5-2", "1.5.1-1", -1},
		{"1.5-2", "1.5.1-2", -1},

		// mixed pkgrel inclusion
		{"1.5", "1.5-1", 0},
		{"1.5-1", "1.5", 0},
		{"1.1-1", "1.1", 0},
		{"1.0-1", "1.1", -1},
		{"1.1-1", "1.0", 1},

		// alphanumeric versions
		{"1.5b-1", "1.5-1", -1},
		{"1.5b", "1.5", -1},
		{"1.5b-1", "1.5", -1},
		{"1.5b", "1.5.1", -1},

		// from the manpage
		{"1.0a", "1.0alpha", -1},
		{"1.0alpha", "1.0b", -1},
		{"1.0b", "1.0beta", -1},
		{"1.0beta", "1.0rc", -1},
		{"1.0rc", "1.0", -1},

		// going crazy? alpha-dotted versions
		{"1.5.a", "1.5", 1},
		{"1.5.b", "1.5.a", 1},
		{"1.5.1", "1.5.b", 1},

		// alpha dots and dashes
		{"1.5.b-1", "1.5.b", 0},
		{"1.5-1", "1.5.b", -1},

		// same/similar content, differing separators
		{"2.0", "2_0", 0},
		{"2.0_a", "2_0.a", 0},
		{"2.0a", "2.0.a", -1},
		{"2___a", "2_a", 1},

		// epoch included version comparisons
		{"0:1.0", "0:1.0", 0},
		{"0:1.0", "0:1.1", -1},
		{"1:1.0", "0:1.0", 1},
		{"1:1.0", "0:1.1", 1},
		{"1:1.0", "2:1.1", -1},

		// epoch + sometimes present pkgrel
		{"1:1.0", "0:1.0-1", 1},
		{"1:1.0-1", "0:1.1-1", 1},

		// epoch included on one version
		{"0:1.0", "1.0", 0},
		{"0:1.0", "1.1", -1},
		{"0:1.1", "1.0", 1},
		{"1:1.0", "1.0", 1},
		{"1:1.0", "1.1", 1},
		{"1:1.1", "1.1", 1},
	}

	for _, tc := range tt {
		got := CompareVersion(tc.in1, tc.in2)
		if got != tc.want {
			t.Errorf("CompareVersion(%q, %q) = %d, want %d", tc.in1, tc.in2, got, tc.want)
		}
	}
}
