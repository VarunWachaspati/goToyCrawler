package main

import "testing"

func TestGetBaseDomain(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"https://golang.com", "https://golang.com"},
		{"http://golang.com", "http://golang.com"},
		{"https://www.golang.com", "https://www.golang.com"},
		{"https://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/", "https://www.datadoghq.com"},
	}
	for _, table := range tables {
		output := GetBaseDomain(table.input)
		if output != table.output {
			t.Errorf("Base Domain of (%s) was incorrect, got: %s, want: %s.", table.input, output, table.output)
		}
	}
}

func TestByteSliceToString(t *testing.T) {
	tables := []struct {
		input  []byte
		output string
	}{
		{[]byte("https://golang.com"), "https://golang.com"},
		{[]byte("http://golang.com"), "http://golang.com"},
		{[]byte("https://www.golang.com"), "https://www.golang.com"},
		{[]byte("<blockquote><p>WriteFile func(filename string, data []byte, perm os.FileMode) error</p></blockquote>"),
			"<blockquote><p>WriteFile func(filename string, data []byte, perm os.FileMode) error</p></blockquote>"},
	}
	for _, table := range tables {
		output := ByteSliceToString(table.input)
		if output != table.output {
			t.Errorf("Conversion of ByteSlice [%v] was incorrect, got: %s, want: %s.", table.input, output, table.output)
		}
	}
}

func TestFormatURL(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"https://golang.com", "https://golang.com"},
		{"http://golang.com", "http://golang.com"},
		{"golang.com", "http://golang.com"},
		{"https://www.golang.com", "https://www.golang.com"},
		{"https://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/", "https://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/"},
		{"www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/", "http://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/"},
	}
	for _, table := range tables {
		output := FormatURL(table.input)
		if output != table.output {
			t.Errorf("URL Format of (%s) was incorrect, got: %s, want: %s.", table.input, output, table.output)
		}
	}
}

func TestResolveToBaseDomain(t *testing.T) {
	tables := []struct {
		input      string
		baseDomain string
		output     string
	}{
		{"/faq", "https://www.sitemaps.org/", "https://www.sitemaps.org/faq"},
		{"faq", "https://www.sitemaps.org/", "https://www.sitemaps.org/faq"},
		{"faq.php", "https://www.sitemaps.org/", "https://www.sitemaps.org/faq.php"},
		{"#faq", "https://www.sitemaps.org/", "https://www.sitemaps.org"},
	}
	for _, table := range tables {
		output := ResolveToBaseDomain(table.input, table.baseDomain)
		if output != table.output {
			t.Errorf("Absolute URL of (%s) and base domain (%s) was incorrect, got: %s, want: %s.", table.input, table.baseDomain, output, table.output)
		}
	}
}

func TestGetHost(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"https://golang.com", "golang.com"},
		{"http://golang.com", "golang.com"},
		{"golang.com", ""},
		{"https://www.golang.com", "www.golang.com"},
		{"https://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/", "www.datadoghq.com"},
		{"http://www.datadoghq.com/blog/crossing-streams-love-letter-gos-io-reader/", "www.datadoghq.com"},
	}
	for _, table := range tables {
		output, _ := GetHost(table.input)
		if output != table.output {
			t.Errorf("URL Format of (%s) was incorrect, got: %s, want: %s.", table.input, output, table.output)
		}
	}
}

func TestSanitizeURL(t *testing.T) {
	tables := []struct {
		input  string
		output string
	}{
		{"http://www.nirantk.in/", "http://www.nirantk.in"},
		{"http://www.nirantk.in?v=iq3rho3", "http://www.nirantk.in"},
	}
	for _, table := range tables {
		output := SanitizeURL(table.input)
		if output != table.output {
			t.Errorf("URL Sanitization of (%s) was incorrect, got: %s, want: %s.", table.input, output, table.output)
		}
	}
}
func TestIsSameDomain(t *testing.T) {
	tables := []struct {
		input      string
		baseDomain string
		output     bool
	}{
		{"https://www.sitemaps.org/faq", "https://www.sitemaps.org", true},
		{"https://www.sitemap.org/faq", "https://www.sitemaps.org", false},
		{"http://www.sitemaps.org/faq", "http://www.sitemaps.org", true},
	}
	for _, table := range tables {
		output := IsSameDomain(table.input, table.baseDomain)
		if output != table.output {
			t.Errorf("Absolute URL of (%s) and base domain (%s) was incorrect, got: %t, want: %t.", table.input, table.baseDomain, output, table.output)
		}
	}
}
