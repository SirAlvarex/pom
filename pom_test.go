package pom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalMarshal(t *testing.T) {
	a := assert.New(t)
	pom, err := Unmarshal([]byte(examplePom))
	a.NoError(err, "Error unmarshalling test data")
	rawPom, err := Marshal(pom)
	a.NoError(err, "Error marshalling test data")
	a.NotEmpty(rawPom, "Pom was empty")
	a.Equal(reMarshaledPom, string(rawPom), "Remarshaled pom is not correct")
}

func TestUpdatedFieldPersists(t *testing.T) {
	a := assert.New(t)
	pom, err := Unmarshal([]byte(examplePom))
	a.NoError(err, "Error unmarshalling test data")
	name := "Test"
	pom.SetName(name)
	rawPom, err := Marshal(pom)
	a.NoError(err, "Error marshalling test data")
	a.NotEmpty(rawPom, "Pom was empty")
	pom, err = Unmarshal([]byte(rawPom))
	a.NoError(err, "Error unmarshalling test data")
	a.NotNil(pom.Name, "Pom Name was not persisted")
	getName, _ := pom.GetName()
	a.Equal(name, getName, "Name was not %s when we remarshaled", name)
}

func TestUpdatedSequence(t *testing.T) {
	a := assert.New(t)
	pom, err := Unmarshal([]byte(examplePom))
	a.NoError(err, "Error unmarshalling test data")
	testName := "This is a Test"
	if repos, ok := pom.GetRepositories(); ok {
		for index, repo := range repos.GetRepository() {
			repo.SetName(testName)
			repos.UpdateRepository(repo, index)
		}
	}
	rawPom, err := Marshal(pom)
	a.NoError(err, "Error marshalling test data")
	a.NotEmpty(rawPom, "Pom was empty")
	pom, err = Unmarshal([]byte(rawPom))
	a.NoError(err, "Error unmarshalling test data")
	a.NotNil(pom.Repositories, "No repos unmarshaled")
	repos, _ := pom.GetRepositories()
	a.Equal(1, len(repos.GetRepository()), "Not enough repos")
	for _, repo := range repos.GetRepository() {
		name, _ := repo.GetName()
		a.Equal(testName, name, "Name does not match")
	}
}
