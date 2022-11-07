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

	sub := startLocalConnection(t)
	defer sub.Close()

	t.Run("kvstore set", func(t *testing.T) {
		err = sub.KVStoreSet(id, key, val)
		assert.NoError(t, err)
	})

	t.Run("kvstore get", func(t *testing.T) {
		t.Skip()
		value, err := sub.KVStoreGet(id.PublicKey(), "key15")
		log.Printf("value: %s", value)
		assert.NoError(t, err)
	})

	t.Run("kvstore delete", func(t *testing.T) {
		t.Skip()
		err = sub.KVStoreDelete(id, key)
		assert.NoError(t, err)
	})

	t.Run("kvstore list keys", func(t *testing.T) {
		value, err := sub.KVStoreList(id)
		log.Printf("values: %s", value)
		assert.NoError(t, err)
	})
}
