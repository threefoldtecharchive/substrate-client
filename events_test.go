package substrate

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/require"
)

func TestEventsTypes(t *testing.T) {
	require := require.New(t)
	mgr := NewManager("wss://tfchain.dev.grid.tf")
	con, meta, err := mgr.Raw()

	require.NoError(err)
	defer con.Client.Close()

	//fmt.Println(meta.Version)
	require.EqualValues(14, meta.Version)
	data := meta.AsMetadataV14

	var known EventRecords
	knownType := reflect.TypeOf(known)
	//data.FindEventNamesForEventID(eventID types.EventID)
	for _, mod := range data.Pallets {
		if !mod.HasEvents {
			continue
		}

		typ, ok := data.EfficientLookup[mod.Events.Type.Int64()]
		if !ok {
			continue
		}
		//fmt.Println("Module: ", mod.Name)
		for _, variant := range typ.Def.Variant.Variants {
			name := fmt.Sprintf("%s_%s", mod.Name, variant.Name)
			filed, ok := knownType.FieldByName(name)
			if !ok {
				t.Errorf("event %s not defined in known events", name)
				continue
			}
			//fmt.Println(" - Event: ", variant.Name)
			t.Run(name, func(t *testing.T) {
				typeValidator(t, name, filed, variant)
			})
		}
	}
}

func typeValidator(t *testing.T, name string, local reflect.StructField, remote types.Si1Variant) {
	require := require.New(t)
	//first of all, each local filed should be an SliceOf(remote) type.
	//which means

	require.True(local.Type.Kind() == reflect.Slice, "found: %+v", local.Type.Kind())
	elem := local.Type.Elem()
	// each element in that array is itself a structure, so we also must do this
	require.True(elem.Kind() == reflect.Struct, "found: %+v", elem.Kind())
}
