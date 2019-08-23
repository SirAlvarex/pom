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
	for index, repo := range pom.Repositories.Repository {
		if repo.Releases != nil {
			value := "true"
			repo.Releases.UpdatePolicy = &value
		}
		if repo.Snapshots != nil {
			value := "true"
			repo.Snapshots.UpdatePolicy = &value
		}
		pom.Repositories.Repository[index] = repo
	}
	rawPom, err := Marshal(pom)
	if err != nil {
		t.Error(err)
	}
	//a.Equal(examplePom, string(rawPom))
	fmt.Println(string(rawPom))
}
