package strkit

import "testing"

func TestToUnderLine(t *testing.T) {
	list := []struct {
		in  string
		out string
	}{
		{in: "HelloWorld", out: "hello_world"},
		{in: "appID", out: "app_id"},
		{in: "pageURL", out: "page_url"},
		{in: "domainDNS", out: "domain_dns"},
		{in: "DomainDNS", out: "domain_dns"},
		{in: "UserName", out: "user_name"},
		{in: "userName", out: "user_name"},
		{in: "user4Name", out: "user4_name"},
		{in: "URL", out: "url"},
		{in: "pageWWW", out: "page_www"},
		{in: "wwwPage", out: "www_page"},
		{in: "WWWpage", out: "wwwpage"},
		{in: "abc", out: "abc"},
		{in: "hhC", out: "hh_c"},
	}

	for _, v := range list {
		under := ToUnderLine(v.in)
		if under != v.out {
			t.Fatalf("in %s,except %s ,but got %s", v.in, v.out, under)
		} else {
			t.Logf("in %s,except %s , got %s", v.in, v.out, under)
		}
	}
}
