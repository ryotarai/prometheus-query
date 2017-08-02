package main

import "testing"

func TestStringMapToString(t *testing.T) {
	s := stringMapToString(map[string]string{
		"a": "1",
		"d": "2",
		"b": "3",
		"e": "4",
		"c": "5",
	})

	if s != "a:1,b:3,c:5,d:2,e:4" {
		t.Error(s)
	}
}
