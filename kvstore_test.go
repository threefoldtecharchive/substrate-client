package substrate

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetKVStore(t *testing.T) {
	key := "key*****"
	val := "val1*****"
	mnem := "route visual hundred rabbit wet crunch ice castle milk model inherit outside"
	id, err := NewIdentityFromSr25519Phrase(mnem)
	if err != nil {
		panic(err)
	}
	sub := startConnection(t)
	defer sub.Close()

	t.Run("kvstore set", func(t *testing.T) {
		t.Skip()
		err = sub.KVStoreSet(key, val, id)
		assert.NoError(t, err)
	})

	t.Run("kvstore get", func(t *testing.T) {
		t.Skip()
		value, err := sub.KVStoreGet("key15", id)
		log.Printf("value: %s", value)
		assert.NoError(t, err)
	})

	t.Run("kvstore delete", func(t *testing.T) {
		t.Skip()
		err = sub.KVSToreDelete(key, id)
		assert.NoError(t, err)
	})

	t.Run("kvstore list keys", func(t *testing.T) {
		value, err := sub.KVStoreList(id)
		log.Printf("values: %s", value)
		assert.NoError(t, err)
	})
}
