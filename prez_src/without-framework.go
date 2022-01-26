package main

import "testing"

func TestGetPerson(t *testing.T) {
	person, err := GetPerson("Bob")
	if err != nil {
		t.Errorf("GetPerson returned error %s", err)
	} else {
		if person.Name != "Bob" {
			t.Errorf(`Name: got=%q expected="Bob"`, person.Name)
		}
		if person.Age != 42 {
			t.Errorf("Age: got=%s expected=42", person.Age)
		}
	}
}
