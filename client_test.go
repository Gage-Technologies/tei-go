package tei_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/gage-technologies/tei-go"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:8080"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	emb, err := c.Embed("Hello world!")
	assert.NoError(t, err)
	assert.NotNil(t, emb)
	assert.Equal(t, 1, len(emb))
	assert.Equal(t, 768, len(emb[0]))

	info, err := c.Info()
	assert.NoError(t, err)
	assert.NotNil(t, info)
	b, _ := json.MarshalIndent(info, "", "  ")
	t.Log(string(b))
}
