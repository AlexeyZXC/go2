package main

import "testing"

func Test_assignInt(t *testing.T) {
	in := In{I: 1, S: "asd"}
	m := map[string]interface{}{"I": -5, "S": "qwe"}

	err := Assign(&in, m)
	if err != nil {
		t.Fatal("err: ", err)
	}

	if in.I != -5 {
		t.Fatal("not equal to 5")
	}
}

func Test_assignString(t *testing.T) {
	in := In{I: 1, S: "asd"}
	m := map[string]interface{}{"I": 5, "S": "qwe"}

	err := Assign(&in, m)
	if err != nil {
		t.Fatal("err: ", err)
	}

	if in.S != "qwe" {
		t.Fatal("not equal to qwe")
	}
}
