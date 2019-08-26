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
}

func TestUpdatedFieldPersists(t *testing.T) {
	a := assert.New(t)
	pom, err := Unmarshal([]byte(examplePom))
	a.NoError(err, "Error unmarshalling test data")
	name := "Test"
	pom.Name = &name
	pom.GetRepositories()
	rawPom, err := Marshal(pom)
	a.NoError(err, "Error marshalling test data")
	a.NotEmpty(rawPom, "Pom was empty")
	pom, err = Unmarshal([]byte(rawPom))
	a.NoError(err, "Error unmarshalling test data")
	a.NotNil(pom.Name, "Pom Name was not persisted")
	a.Equal(name, *pom.Name, "Name was not %s when we remarshaled", name)
}
