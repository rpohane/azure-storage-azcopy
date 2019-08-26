// Copyright © 2017 Microsoft <wastore@microsoft.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/Azure/azure-storage-azcopy/common"
)

// the interceptor gathers/saves the job part orders for validation
type interceptor struct {
	transfers   []common.CopyTransfer
	lastRequest interface{}
}

func (i *interceptor) intercept(cmd common.RpcCmd, request interface{}, response interface{}) {
	switch cmd {
	case common.ERpcCmd.CopyJobPartOrder():
		// cache the transfers
		copyRequest := *request.(*common.CopyJobPartOrderRequest)
		i.transfers = append(i.transfers, copyRequest.Transfers...)
		i.lastRequest = request

		// mock the result
		*(response.(*common.CopyJobPartOrderResponse)) = common.CopyJobPartOrderResponse{JobStarted: true}
	case common.ERpcCmd.ListJobs():
	case common.ERpcCmd.ListJobSummary():
	case common.ERpcCmd.ListJobTransfers():
	case common.ERpcCmd.PauseJob():
	case common.ERpcCmd.CancelJob():
	case common.ERpcCmd.ResumeJob():
	case common.ERpcCmd.GetJobFromTo():
		fallthrough
	default:
		panic("RPC mock not implemented")
	}
}

func (i *interceptor) init() {
	// mock out the lifecycle manager so that it can no longer terminate the application
	glcm = mockedLifecycleManager{}
}

func (i *interceptor) reset() {
	i.transfers = make([]common.CopyTransfer, 0)
	i.lastRequest = nil
}

// this lifecycle manager substitute does not perform any action
type mockedLifecycleManager struct{}

func (mockedLifecycleManager) Progress(common.OutputBuilder)                            {}
func (mockedLifecycleManager) Init(common.OutputBuilder)                                {}
func (mockedLifecycleManager) Info(msg string)                                          { fmt.Println(msg) }
func (mockedLifecycleManager) Prompt(string) string                                     { return "" }
func (mockedLifecycleManager) Exit(common.OutputBuilder, common.ExitCode)               {}
func (mockedLifecycleManager) Error(string)                                             {}
func (mockedLifecycleManager) SurrenderControl()                                        {}
func (mockedLifecycleManager) InitiateProgressReporting(common.WorkController, bool)    {}
func (mockedLifecycleManager) GetEnvironmentVariable(common.EnvironmentVariable) string { return "" }
func (mockedLifecycleManager) SetOutputFormat(common.OutputFormat)                      {}

type dummyProcessor struct {
	record []storedObject
}

func (d *dummyProcessor) process(storedObject storedObject) (err error) {
	d.record = append(d.record, storedObject)
	return
}
