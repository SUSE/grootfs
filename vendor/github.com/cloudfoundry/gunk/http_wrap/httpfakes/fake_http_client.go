// This file was generated by counterfeiter
package httpfakes

import (
	os_http "net/http"
	"sync"

	"github.com/cloudfoundry/gunk/http_wrap"
)

type FakeClient struct {
	DoStub        func(req *os_http.Request) (resp *os_http.Response, err error)
	doMutex       sync.RWMutex
	doArgsForCall []struct {
		req *os_http.Request
	}
	doReturns struct {
		result1 *os_http.Response
		result2 error
	}
}

func (fake *FakeClient) Do(req *os_http.Request) (resp *os_http.Response, err error) {
	fake.doMutex.Lock()
	fake.doArgsForCall = append(fake.doArgsForCall, struct {
		req *os_http.Request
	}{req})
	fake.doMutex.Unlock()
	if fake.DoStub != nil {
		return fake.DoStub(req)
	} else {
		return fake.doReturns.result1, fake.doReturns.result2
	}
}

func (fake *FakeClient) DoCallCount() int {
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	return len(fake.doArgsForCall)
}

func (fake *FakeClient) DoArgsForCall(i int) *os_http.Request {
	fake.doMutex.RLock()
	defer fake.doMutex.RUnlock()
	return fake.doArgsForCall[i].req
}

func (fake *FakeClient) DoReturns(result1 *os_http.Response, result2 error) {
	fake.DoStub = nil
	fake.doReturns = struct {
		result1 *os_http.Response
		result2 error
	}{result1, result2}
}

var _ http_wrap.Client = new(FakeClient)
