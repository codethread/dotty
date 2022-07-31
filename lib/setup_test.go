package lib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	t.Skip()
	assert.NotPanics(t,
		func() {
			Setup(SetupConfig{})
		},
	)
}

func TestSetupFixture(t *testing.T) {
	t.Skip()

	config := SetupConfig{
		From: getFixtureDir(),
		To:   getFixtureTarget(),
	}

	t.Cleanup(func() {
		os.RemoveAll(config.To)
	})
}

func getFixtureTarget() string {
	panic("unimplemented")
}

func getFixtureDir() string {
	panic("unimplemented")
}
