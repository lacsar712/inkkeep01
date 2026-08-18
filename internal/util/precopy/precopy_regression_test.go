package precopy

import "testing"

// TestClonePrefixNoAliasRegression guards against ClonePrefix returning a slice
// that shares its backing array with src: every element of the returned prefix
// must be writable without touching src.
func TestClonePrefixNoAliasRegression(t *testing.T) {
	src := []byte("abcdef")
	got := ClonePrefix(src, 4)
	if string(got) != "abcd" {
		t.Fatalf("prefix=%q, want %q", got, "abcd")
	}
	for i := range got {
		got[i] = 'Z'
		if src[i] != byte("abcdef"[i]) {
			t.Fatalf("ClonePrefix aliased src at index %d: src=%q", i, src)
		}
	}
}

func TestClonePrefixBounds(t *testing.T) {
	cases := []struct {
		name string
		src  []byte
		n    int
		want string
	}{
		{"negative clamps to zero", []byte("abc"), -1, ""},
		{"zero returns empty", []byte("abc"), 0, ""},
		{"full length", []byte("abc"), 3, "abc"},
		{"over length clamps", []byte("abc"), 10, "abc"},
		{"partial", []byte("abcdef"), 3, "abc"},
		{"empty source", []byte{}, 0, ""},
		{"empty source with n", []byte{}, 5, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ClonePrefix(c.src, c.n)
			if string(got) != c.want {
				t.Fatalf("ClonePrefix(%q, %d) = %q, want %q", c.src, c.n, got, c.want)
			}
		})
	}
}

func TestClonePrefixEmptyIndependent(t *testing.T) {
	src := []byte("abc")
	got := ClonePrefix(src, 0)
	if len(got) != 0 {
		t.Fatalf("len(got)=%d, want 0", len(got))
	}
	// Even an empty result must not share storage that could later be mutated
	// through grow; cap should not exceed the requested length.
	if cap(got) != 0 {
		t.Fatalf("cap(got)=%d, want 0 (no aliasing of src backing array)", cap(got))
	}
}

func TestCloneAllNoAliasRegression(t *testing.T) {
	src := []byte{1, 2, 3, 4}
	got := CloneAll(src)
	if len(got) != len(src) {
		t.Fatalf("len(got)=%d, want %d", len(got), len(src))
	}
	for i := range got {
		got[i] = got[i] + 100
		if src[i] == got[i] {
			t.Fatalf("CloneAll aliased src at %d: src=%v", i, src)
		}
	}
}
