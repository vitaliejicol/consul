package cachetype

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/agent/cache"
	"github.com/hashicorp/consul/agent/structs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCatalogListServices(t *testing.T) {
	rpc := TestRPC(t)
	typ := &CatalogListServices{RPC: rpc}

	// Expect the proper RPC call. This also sets the expected value
	// since that is return-by-pointer in the arguments.
	var resp *structs.IndexedServices
	rpc.On("RPC", "Catalog.ListServices", mock.Anything, mock.Anything).Return(nil).
		Run(func(args mock.Arguments) {
			req := args.Get(1).(*structs.DCSpecificRequest)
			require.Equal(t, uint64(24), req.QueryOptions.MinQueryIndex)
			require.Equal(t, 1*time.Second, req.QueryOptions.MaxQueryTime)
			require.True(t, req.AllowStale)

			reply := args.Get(2).(*structs.IndexedServices)
			reply.Services = map[string][]string{
				"foo": {"prod", "linux"},
				"bar": {"qa", "windows"},
			}
			reply.QueryMeta.Index = 48
			resp = reply
		})

	// Fetch
	resultA, err := typ.Fetch(cache.FetchOptions{
		MinIndex: 24,
		Timeout:  1 * time.Second,
	}, &structs.DCSpecificRequest{
		Datacenter: "dc1",
	})
	require.NoError(t, err)
	require.Equal(t, cache.FetchResult{
		Value: resp,
		Index: 48,
	}, resultA)

	rpc.AssertExpectations(t)
}

func TestCatalogListServices_badReqType(t *testing.T) {
	rpc := TestRPC(t)
	typ := &CatalogServices{RPC: rpc}

	// Fetch
	_, err := typ.Fetch(cache.FetchOptions{}, cache.TestRequest(
		t, cache.RequestInfo{Key: "foo", MinIndex: 64}))
	require.Error(t, err)
	require.Contains(t, err.Error(), "wrong type")
	rpc.AssertExpectations(t)
}

func TestCatalogListServices_NotModifiedResponse(t *testing.T) {
	rpc := &MockRPC{}
	typ := &CatalogListServices{RPC: rpc}

	services := map[string][]string{
		"foo": {"prod", "linux"},
		"bar": {"qa", "windows"},
	}
	rpc.On("RPC", "Catalog.ListServices", mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			req := args.Get(1).(*structs.DCSpecificRequest)
			require.Equal(t, uint64(44), req.QueryOptions.MinQueryIndex)
			require.Equal(t, 1*time.Second, req.QueryOptions.MaxQueryTime)
			require.True(t, req.AllowStale)
			require.True(t, req.AllowNotModifiedResponse)

			reply := args.Get(2).(*structs.IndexedServices)
			reply.QueryMeta.Index = 44
			reply.NotModified = true
		})

	opts := cache.FetchOptions{
		MinIndex: 44,
		Timeout:  time.Second,
		LastResult: &cache.FetchResult{
			Value: &structs.IndexedServices{
				Services:  services,
				QueryMeta: structs.QueryMeta{Index: 44},
			},
			Index: 44,
		},
	}
	actual, err := typ.Fetch(opts, &structs.DCSpecificRequest{Datacenter: "dc1"})
	require.NoError(t, err)
	expected := cache.FetchResult{
		Value: &structs.IndexedServices{
			Services:  services,
			QueryMeta: structs.QueryMeta{Index: 44},
		},
		Index: 44,
	}
	require.Equal(t, expected, actual)

	rpc.AssertExpectations(t)
}
