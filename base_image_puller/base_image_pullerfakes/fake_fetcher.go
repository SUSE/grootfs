// Code generated by counterfeiter. DO NOT EDIT.
package base_image_pullerfakes

import (
	"io"
	"sync"

	"code.cloudfoundry.org/grootfs/base_image_puller"
	"code.cloudfoundry.org/grootfs/groot"
	"code.cloudfoundry.org/lager"
)

type FakeFetcher struct {
	BaseImageInfoStub        func(logger lager.Logger) (groot.BaseImageInfo, error)
	baseImageInfoMutex       sync.RWMutex
	baseImageInfoArgsForCall []struct {
		logger lager.Logger
	}
	baseImageInfoReturns struct {
		result1 groot.BaseImageInfo
		result2 error
	}
	baseImageInfoReturnsOnCall map[int]struct {
		result1 groot.BaseImageInfo
		result2 error
	}
	StreamBlobStub        func(logger lager.Logger, layerInfo groot.LayerInfo) (io.ReadCloser, int64, error)
	streamBlobMutex       sync.RWMutex
	streamBlobArgsForCall []struct {
		logger    lager.Logger
		layerInfo groot.LayerInfo
	}
	streamBlobReturns struct {
		result1 io.ReadCloser
		result2 int64
		result3 error
	}
	streamBlobReturnsOnCall map[int]struct {
		result1 io.ReadCloser
		result2 int64
		result3 error
	}
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	closeReturns     struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeFetcher) BaseImageInfo(logger lager.Logger) (groot.BaseImageInfo, error) {
	fake.baseImageInfoMutex.Lock()
	ret, specificReturn := fake.baseImageInfoReturnsOnCall[len(fake.baseImageInfoArgsForCall)]
	fake.baseImageInfoArgsForCall = append(fake.baseImageInfoArgsForCall, struct {
		logger lager.Logger
	}{logger})
	fake.recordInvocation("BaseImageInfo", []interface{}{logger})
	fake.baseImageInfoMutex.Unlock()
	if fake.BaseImageInfoStub != nil {
		return fake.BaseImageInfoStub(logger)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.baseImageInfoReturns.result1, fake.baseImageInfoReturns.result2
}

func (fake *FakeFetcher) BaseImageInfoCallCount() int {
	fake.baseImageInfoMutex.RLock()
	defer fake.baseImageInfoMutex.RUnlock()
	return len(fake.baseImageInfoArgsForCall)
}

func (fake *FakeFetcher) BaseImageInfoArgsForCall(i int) lager.Logger {
	fake.baseImageInfoMutex.RLock()
	defer fake.baseImageInfoMutex.RUnlock()
	return fake.baseImageInfoArgsForCall[i].logger
}

func (fake *FakeFetcher) BaseImageInfoReturns(result1 groot.BaseImageInfo, result2 error) {
	fake.BaseImageInfoStub = nil
	fake.baseImageInfoReturns = struct {
		result1 groot.BaseImageInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeFetcher) BaseImageInfoReturnsOnCall(i int, result1 groot.BaseImageInfo, result2 error) {
	fake.BaseImageInfoStub = nil
	if fake.baseImageInfoReturnsOnCall == nil {
		fake.baseImageInfoReturnsOnCall = make(map[int]struct {
			result1 groot.BaseImageInfo
			result2 error
		})
	}
	fake.baseImageInfoReturnsOnCall[i] = struct {
		result1 groot.BaseImageInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeFetcher) StreamBlob(logger lager.Logger, layerInfo groot.LayerInfo) (io.ReadCloser, int64, error) {
	fake.streamBlobMutex.Lock()
	ret, specificReturn := fake.streamBlobReturnsOnCall[len(fake.streamBlobArgsForCall)]
	fake.streamBlobArgsForCall = append(fake.streamBlobArgsForCall, struct {
		logger    lager.Logger
		layerInfo groot.LayerInfo
	}{logger, layerInfo})
	fake.recordInvocation("StreamBlob", []interface{}{logger, layerInfo})
	fake.streamBlobMutex.Unlock()
	if fake.StreamBlobStub != nil {
		return fake.StreamBlobStub(logger, layerInfo)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.streamBlobReturns.result1, fake.streamBlobReturns.result2, fake.streamBlobReturns.result3
}

func (fake *FakeFetcher) StreamBlobCallCount() int {
	fake.streamBlobMutex.RLock()
	defer fake.streamBlobMutex.RUnlock()
	return len(fake.streamBlobArgsForCall)
}

func (fake *FakeFetcher) StreamBlobArgsForCall(i int) (lager.Logger, groot.LayerInfo) {
	fake.streamBlobMutex.RLock()
	defer fake.streamBlobMutex.RUnlock()
	return fake.streamBlobArgsForCall[i].logger, fake.streamBlobArgsForCall[i].layerInfo
}

func (fake *FakeFetcher) StreamBlobReturns(result1 io.ReadCloser, result2 int64, result3 error) {
	fake.StreamBlobStub = nil
	fake.streamBlobReturns = struct {
		result1 io.ReadCloser
		result2 int64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeFetcher) StreamBlobReturnsOnCall(i int, result1 io.ReadCloser, result2 int64, result3 error) {
	fake.StreamBlobStub = nil
	if fake.streamBlobReturnsOnCall == nil {
		fake.streamBlobReturnsOnCall = make(map[int]struct {
			result1 io.ReadCloser
			result2 int64
			result3 error
		})
	}
	fake.streamBlobReturnsOnCall[i] = struct {
		result1 io.ReadCloser
		result2 int64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeFetcher) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.closeReturns.result1
}

func (fake *FakeFetcher) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeFetcher) CloseReturns(result1 error) {
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeFetcher) CloseReturnsOnCall(i int, result1 error) {
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeFetcher) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.baseImageInfoMutex.RLock()
	defer fake.baseImageInfoMutex.RUnlock()
	fake.streamBlobMutex.RLock()
	defer fake.streamBlobMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeFetcher) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ base_image_puller.Fetcher = new(FakeFetcher)
