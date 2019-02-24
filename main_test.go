package main

import "testing"

func TestMakePhraseNot(t *testing.T) {

	actualResult := makePhrase(false, "foo")

	var expectedResult = "not foo"

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func TestMakePhrase(t *testing.T) {

	actualResult := makePhrase(true, "foo")

	var expectedResult = "foo"

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}
