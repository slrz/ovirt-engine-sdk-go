//
// Copyright (c) 2017 Joey <majunjiev@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package examples

import (
	"fmt"
	"time"

	ovirtsdk4 "gopkg.in/imjoey/go-ovirt.v4"
)

func listVMDisks() {
	inputRawURL := "https://10.1.111.229/ovirt-engine/api"

	conn, err := ovirtsdk4.NewConnectionBuilder().
		URL(inputRawURL).
		Username("admin@internal").
		Password("qwer1234").
		Insecure(true).
		Compress(true).
		Timeout(time.Second * 10).
		Build()
	if err != nil {
		fmt.Printf("Make connection failed, reason: %v\n", err)
		return
	}
	defer conn.Close()

	// To use `Must` methods, you should recover it if panics
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panics occurs, try the non-Must methods to find the reason")
		}
	}()

	vmsService := conn.SystemService().VmsService()
	// Use Must version methods to make a function call chain
	vm := vmsService.List().Search("name=test4joey").MustSend().MustVms().Slice()[0]
	vmService := vmsService.VmService(vm.MustId())
	dasService := vmService.DiskAttachmentsService()

	das := dasService.List().MustSend().MustAttachments()
	for _, da := range das.Slice() {
		disk, _ := conn.FollowLink(da.MustDisk())
		if disk, ok := disk.(*ovirtsdk4.Disk); ok {
			fmt.Printf(" name: %v\n", disk.MustName())
			fmt.Printf(" id: %v\n", disk.MustId())
			fmt.Printf(" status: %v\n", disk.MustStatus())
			fmt.Printf(" provisioned_size: %v\n", disk.MustProvisionedSize())
		}
	}
}
