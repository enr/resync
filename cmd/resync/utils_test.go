package main

import "testing"

type isInsideTestCase struct {
	p      string
	d      string
	result bool
}

var isInsideFixture = []isInsideTestCase{
	{
		p:      "",
		d:      "/tmp/a",
		result: false,
	},
	{
		p:      "/tmp/a",
		d:      "",
		result: false,
	},
	{
		p:      "/tmp/a",
		d:      "/",
		result: true,
	},
	{
		p:      "/tmp/a",
		d:      "tmp",
		result: false,
	},
	{
		p:      "/tmp/a/b/c/d",
		d:      "/tmp/a",
		result: true,
	},
	{
		p:      "/tmp/a/b",
		d:      "/tmp/a/b/c",
		result: false,
	},
	{
		p:      "/tmp/a/b",
		d:      "/tmp/a/b/",
		result: true,
	},
	{
		p:      "/tmp/a/b",
		d:      "/tmp/a/b-2",
		result: false,
	},
	{
		p:      "/tmp/a/b-2",
		d:      "/tmp/a/b",
		result: false,
	},
}

func TestIsInside(t *testing.T) {
	for _, testcase := range isInsideFixture {
		res := isInside(testcase.p, testcase.d)
		if testcase.result != res {
			t.Errorf(`%s -> %s error result %t`, testcase.p, testcase.d, res)
		}
	}
}
