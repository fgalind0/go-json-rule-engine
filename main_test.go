package main

import (
	"testing"
)

func TestStringEqualsNoMatch(t *testing.T) {
	var genericCondition Condition
	example := StringCondition{
		Subject:  "fern",
		Operator: "StringEquals",
		Value:    "fernando",
	}
	genericCondition = &example
	expect := false
	received, err := genericCondition.Check()
	if expect != received || err != nil {
		t.Logf("[TestStringEquals] failed due to receiving: %v but expecting: %v ", received, expect)
		t.Fail()
	}
}
func TestStringEqualsMatch(t *testing.T) {
	var genericCondition Condition
	example := StringCondition{
		Subject:  "fernando",
		Operator: "StringEquals",
		Value:    "fernando",
	}
	genericCondition = &example
	expect := true
	received, err := genericCondition.Check()
	if expect != received || err != nil {
		t.Logf("[TestStringEquals] failed due to receiving: %v but expecting: %v ", received, expect)
		t.Fail()
	}
}
