package indexer

import (
	"github.com/dtylman/pictures/conf"
	"github.com/dtylman/pictures/db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	conf.Options.SourceFolders = []string{"/home/danny/Pictures"}
	t.Log(conf.Options.SourceFolders)
	err := Start(Options{IndexLocation: false, ReIndex: false})
	assert.NoError(t, err)
	assert.True(t, IsRunning())
	err = Start(Options{IndexLocation: false, ReIndex: false})
	assert.Error(t, err)
	for IsRunning() {
		time.Sleep(time.Second)
		t.Log("Still running")
	}
	sr, err := db.QueryAll(0, 100)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sr.String())
	for _, facet := range sr.Facets {
		t.Log(facet)
	}
}
