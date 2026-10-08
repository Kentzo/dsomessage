package dsomessage

import (
	"cmp"
	"slices"
	"testing"
)

func TestLess(t *testing.T) {
	tests := []struct {
		in  []string
		out []string
	}{
		{
			[]string{"aaa.powerdns.de.", "bbb.powerdns.net.", "xxx.powerdns.com."},
			[]string{"xxx.powerdns.com.", "aaa.powerdns.de.", "bbb.powerdns.net."},
		},
		{
			[]string{"aaa.POWERDNS.de.", "bbb.PoweRdnS.net.", "xxx.powerdns.com."},
			[]string{"xxx.powerdns.com.", "aaa.POWERDNS.de.", "bbb.PoweRdnS.net."},
		},
		{
			[]string{"aaa.aaaa.aa.", "aa.aaa.a.", "bbb.bbbb.bb."},
			[]string{"aa.aaa.a.", "aaa.aaaa.aa.", "bbb.bbbb.bb."},
		},
		{
			[]string{"aaaaa.", "aaa.", "bbb."},
			[]string{"aaa.", "aaaaa.", "bbb."},
		},
		{
			[]string{"a.a.a.a.", "a.a.", "a.a.a."},
			[]string{"a.a.", "a.a.a.", "a.a.a.a."},
		},
		{
			[]string{"example.", "z.example.", "a.example."},
			[]string{"example.", "a.example.", "z.example."},
		},
		{
			[]string{"a.example.", "Z.a.example.", "z.example.", "yljkjljk.a.example.", "\\001.z.example.", "example.", "*.z.example.", "\\200.z.example.", "zABC.a.EXAMPLE."},
			[]string{"example.", "a.example.", "yljkjljk.a.example.", "Z.a.example.", "zABC.a.EXAMPLE.", "z.example.", "\\001.z.example.", "*.z.example.", "\\200.z.example."},
		},
		{
			// RFC3034 example.
			[]string{"a.example.", "Z.a.example.", "z.example.", "yljkjljk.a.example.", "example.", "*.z.example.", "zABC.a.EXAMPLE."},
			[]string{"example.", "a.example.", "yljkjljk.a.example.", "Z.a.example.", "zABC.a.EXAMPLE.", "z.example.", "*.z.example."},
		},
	}

Tests:
	for j, test := range tests {
		slices.SortFunc(test.in, CompareNames)
		for i := range len(test.in) {
			if test.in[i] != test.out[i] {
				t.Errorf("Test %d: expected %s, got %s", j, test.out[i], test.in[i])
				n := ""
				for k, in := range test.in {
					if k+1 == len(test.in) {
						n = "\n"
					}
					t.Logf("%s <-> %s\n%s", in, test.out[k], n)
				}
				continue Tests
			}
		}
	}
}

func TestLess_EdgeCases(t *testing.T) {
	// For every case four variants are synthesized:
	//  - a  b
	//  - a. b
	//  - a  b.
	//  - a. b.
	//
	// For each variant commutativity is tested.
	tests := []struct {
		a, b     string
		variants bool
		want     int
	}{
		{``, ``, true, 0},
		{``, `\000`, true, -1},
		{``, `\.`, true, -1},
		{`\.`, `\.`, true, 0},
		{``, `example`, true, -1},
		{`example`, `example`, true, 0},
		{`a\.example`, `a.example`, true, -1},
		{`a.example`, `a-b.example`, true, -1},
		{`a.example`, `a*.example`, true, -1},
		{`a.example`, `a\000.example`, true, -1},
		{`a.eXaMpLe`, `a.example`, true, 0},
		{`\000\0320 \"\046@*`, `\000\032\048\032\034\046\064\042`, true, 0},
		{`<=>?@ABCDE`, `\060\061\062\063\064\065\066\067\068\069`, true, 0},
		{`<=>?@ABCDE`, `\060\061\062\063\064\097\098\099\100\101`, true, 0},
		{`café.example`, `CAFÉ.example`, true, 1}, // é (\195\169) > É (\195\137)
		{`\\065.example`, `\\097.example`, true, -1},
		{``, `\`, false, 0},
		{`a\.b.example`, `a\046b.example`, true, 0},
		{`0.example`, `\0.example`, true, 0},
		{`01.example`, `\01.example`, true, 0},
		{`a\.b.example`, `a\046b.example`, true, 0},
		{`\\.example`, `\092.example`, true, 0},
	}
	for i, test := range tests {
		variants := []struct{ a, b string }{
			{test.a, test.b},
			{test.a + `.`, test.b},
			{test.a, test.b + `.`},
			{test.a + `.`, test.b + `.`},
		}
		if !test.variants {
			variants = variants[:1]
		}

		for _, variant := range variants {
			if got := CompareNames(variant.a, variant.b); cmp.Compare(got, 0) != test.want {
				t.Errorf("Test %d: expected less(%s, %s)=%d, got %d", i, variant.a, variant.b, test.want, cmp.Compare(got, 0))
			}

			if got := CompareNames(variant.b, variant.a); cmp.Compare(got, 0) != -test.want {
				t.Errorf("Test %d: expected less(%s, %s)=%d, got %d", i, variant.b, variant.a, -test.want, cmp.Compare(got, 0))
			}
		}
	}
}
