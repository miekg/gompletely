package main

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestUnmarshalYAML(t *testing.T) {
	buf, err := os.ReadFile("testdata/AddVolume.yml")
	if err != nil {
		t.Fatal(err)
	}
	d := Definition{}

	if err := yaml.Unmarshal(buf, &d); err != nil {
		t.Error(err)
	}
	fmt.Printf("+%v\n", d)
}
