package main

import (
	"bytes"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestUnmarshalYAML(t *testing.T) {
	buf, err := os.ReadFile("testdata/AddVolume.yml")
	if err != nil {
		t.Fatal(err)
	}
	p := Patterns{}

	if err := yaml.Unmarshal(buf, &p); err != nil {
		t.Error(err)
	}
}

func TestBash(t *testing.T) {
	buf, err := os.ReadFile("testdata/AddVolume.yml")
	if err != nil {
		t.Fatal(err)
	}
	p := Patterns{}
	if err := yaml.Unmarshal(buf, &p); err != nil {
		t.Error(err)
	}

	b := ToBash(p)
	out := &bytes.Buffer{}
	if err = bashtmpl.Execute(out, b); err != nil {
		t.Fatal(err)
	}
	println("OUT", out.String())
}
