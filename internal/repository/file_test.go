package repository

import (
	"github.com/mingalevme/avito/internal/model"
	log "github.com/mingalevme/gologger"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func CreateTestFileRepository(t *testing.T) *FileRepository {
	filename := strings.TrimRight(os.TempDir(), string(os.PathSeparator)) + string(os.PathSeparator) + "mingalevme-avito-" + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	r, err := NewFileRepository(filename, log.NewStdoutLogger(log.LevelDebug))
	assert.NoError(t, err)
	return r
}

func TestFileRepository(t *testing.T) {
	r := CreateTestFileRepository(t)
	items, err := r.GetAll()
	assert.NoError(t, err)
	assert.Len(t, items, 0)
	item := model.Item{
		ID:    1,
		URL:   "https://example.com/1",
		Title: "Test",
		Price: 100,
	}
	err = r.Add(item)
	assert.NoError(t, err)
	items, err = r.GetAll()
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	item, err = r.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, item.ID)
	err = r.Clean()
	assert.NoError(t, err)
}
