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

package main

import (
	"fmt"
	"time"

	ovirtsdk4 "gopkg.in/imjoey/go-ovirt.v4"
)

func main() {
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

	// Find the service that manages the collection of hosts
	hostsService := conn.SystemService().HostsService()

	// Find the host
	listResp, err := hostsService.List().Search("name=myhost").Send()
	if err != nil {
		fmt.Printf("Failed to search host list, reason: %v\n", err)
		return
	}
	hostSlice, _ := listResp.Hosts()
	host := hostSlice.Slice()[0]

	// Find the service that manages the hosts
	hostService := hostsService.HostService(host.MustId())

	// Configure the host adding a bond with two slaves, and attaching it to a network with an static IP address
	setupNetworkResp, err := hostService.SetupNetworks().
		ModifiedBondsOfAny(
			ovirtsdk4.NewHostNicBuilder().
				Name("bond0").
				Bonding(
					ovirtsdk4.NewBondingBuilder().
						OptionsOfAny(
							*ovirtsdk4.NewOptionBuilder().
								Name("mode").
								Value("1").
								MustBuild(),
							*ovirtsdk4.NewOptionBuilder().
								Name("miimon").
								Value("100").
								MustBuild()).
						SlavesOfAny(
							*ovirtsdk4.NewHostNicBuilder().
								Name("eth0").
								MustBuild(),
							*ovirtsdk4.NewHostNicBuilder().
								Name("eth1").
								MustBuild()).
						MustBuild()).
				MustBuild()).
		ModifiedNetworkAttachmentsOfAny(
			ovirtsdk4.NewNetworkAttachmentBuilder().
				Network(
					ovirtsdk4.NewNetworkBuilder().
						Name("mynetwork").
						MustBuild()).
				IpAddressAssignmentsOfAny(
					*ovirtsdk4.NewIpAddressAssignmentBuilder().
						AssignmentMethod(
							ovirtsdk4.BOOTPROTOCOL_STATIC).
						Ip(
							ovirtsdk4.NewIpBuilder().
								Address("192.168.122.100").
								Netmask("255.255.255.0").
								MustBuild()).
						MustBuild()).
				MustBuild()).
		Send()
	if err != nil {
		fmt.Printf("Failed to setup network for host-%v, reason: %v\n", host.MustId(), err)
		return
	}

	// After modifying the network configuration it is very important to make it persistent
	hostService.CommitNetConfig().Send()

}
