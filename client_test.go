package tei_test

import (
	"os"
	"testing"
	"time"

	"github.com/gage-technologies/tei-go"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:22123"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	emb, err := c.Embed("Hello world!", false)
	assert.NoError(t, err)
	assert.NotNil(t, emb)
	assert.Equal(t, 1, len(emb))
	assert.Equal(t, 1024, len(emb[0]))

	info, err := c.Info()
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "WhereIsAI/UAE-Large-V1", info.ModelID)
}

func TestEmbedAll(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:22123"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	emb, err := c.EmbedAll([]string{"Deep learning transforms everything."}, false)
	assert.NoError(t, err)
	assert.NotNil(t, emb)
	assert.NotEmpty(t, emb) // Ensure not empty since the actual content depends on the model
	assert.Equal(t, 1, len(emb))
	assert.Equal(t, 7, len(emb[0]))
	assert.Equal(t, 1024, len(emb[0][0]))
}

func TestEmbedSparse(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:22126"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	sparseEmb, err := c.EmbedSparse([]string{"Machine learning"}, false)
	assert.NoError(t, err)
	assert.NotNil(t, sparseEmb)
	assert.NotEmpty(t, sparseEmb) // Check for non-empty response
	assert.Equal(t, 1, len(sparseEmb))
	assert.Equal(t, 3, len(sparseEmb[0]))
}

func TestPredict(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:22125"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	resp, err := c.Predict("What is the capital of France?")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp) // Validate the response structure
	assert.Equal(t, 28, len(resp))
	top := tei.Prediction{}
	for _, prediction := range resp {
		prediction := prediction
		if prediction.Score > top.Score {
			top = prediction
		}
	}
	assert.Equal(t, "curiosity", top.Label)
}

func TestRerank(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:22124"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	query := "Which framework is best for deep learning?"
	texts := []string{"TensorFlow", "PyTorch", "Keras"}
	ranks, err := c.Rerank(query, texts)
	assert.NoError(t, err)
	assert.NotNil(t, ranks)
	assert.Equal(t, 3, len(ranks)) // Check the number of ranks returned
	// this model likes pytorch the most
	top := ""
	top_val := float32(0.0)
	for _, rank := range ranks {
		if rank.Score > top_val {
			top = rank.Text
			top_val = rank.Score
		}
	}
	assert.Equal(t, "PyTorch", top)
}

func TestHealthCheck(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:22123"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	healthy, err := c.HealthCheck()
	assert.NoError(t, err)
	assert.True(t, healthy) // Assuming the service is up during test
}

func TestTokenize(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:22123"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	tokens, err := c.Tokenize([]string{"Hello world!"}, true)
	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.Greater(t, len(tokens), 0) // Check that tokens are returned
	ids := make([]int, 0, len(tokens[0]))
	for _, token := range tokens[0] {
		ids = append(ids, token.ID)
	}
	assert.Equal(t, []int{101, 7592, 2088, 999, 102}, ids)
}

func TestDecode(t *testing.T) {
	host := os.Getenv("TEI_HOST")
	if host == "" {
		host = "http://localhost:22123"
	}
	c := tei.NewClient(host, nil, nil, time.Second)

	decoded, err := c.Decode([][]int{{101, 7592, 2088, 999, 102}}, true)
	assert.NoError(t, err)
	assert.NotNil(t, decoded)
	assert.Equal(t, []string{"hello world!"}, decoded)
}
