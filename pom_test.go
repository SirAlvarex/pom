package pom

import (
	"fmt"
	"testing"
)

func TestUnmarshalMarshal(t *testing.T) {
	//a := assert.New(t)
	pom, err := Unmarshal([]byte(examplePom))
	if err != nil {
		t.Error(err)
	}
	rawPom, err := Marshal(pom)
	//a.Equal(examplePom, string(rawPom))
	fmt.Println(string(rawPom))
}
