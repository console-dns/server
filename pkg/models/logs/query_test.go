package logs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryAfter(t *testing.T) {
	now := time.Now()
	data := &Meta{
		CreatedAt: now,
	}
	assert.Equal(t, QueryAfter(now.Add(+time.Hour))(data), STOP)
	assert.Equal(t, QueryAfter(now.Add(-time.Hour))(data), OK)
}

func TestQueryAuthor(t *testing.T) {
}

func TestQueryBefore(t *testing.T) {
	now := time.Now()
	data := &Meta{
		CreatedAt: now,
	}
	assert.Equal(t, QueryBefore(now.Add(+time.Hour))(data), OK)
	assert.Equal(t, QueryBefore(now.Add(-time.Hour))(data), STOP)
}

func TestQueryGroup(t *testing.T) {
}

func TestQueryIpAddr(t *testing.T) {
}
