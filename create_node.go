package main

import (
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/baremetal/noauth"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
)

func waitForState(client *gophercloud.ServiceClient, timeout int, uuid string, state string) error {
	for i := 0; i < timeout; i++ {
		node, err := nodes.Get(client, uuid).Extract()
		if err != nil {
			return err
		}
		fmt.Printf("%s %d %s - %s\n", time.Now(), i, uuid, node.ProvisionState)
		if node.UUID == uuid && node.ProvisionState == state {
			return nil
		}
		time.Sleep(time.Second * 10)
	}
	return fmt.Errorf("%s never entered state %s", uuid, state)
}

func main() {
	client, err := noauth.NewBareMetalNoAuth(noauth.EndpointOpts{
		IronicEndpoint: "http://localhost:6385/v1/",
	})
	if err != nil {
		panic(err)
	}

	client.Microversion = "1.50"

	// Example to Create Node
	createNode, err := nodes.Create(client, nodes.CreateOpts{
		Driver:        "ipmi",
		BootInterface: "pxe",
		Name:          "worker-0",
		DriverInfo: map[string]interface{}{
			"ipmi_port":      "6233",
			"ipmi_username":  "admin",
			"deploy_kernel":  "http://172.22.0.1/images/tinyipa-stable-rocky.vmlinuz",
			"ipmi_address":   "192.168.122.1",
			"deploy_ramdisk": "http://172.22.0.1/images/tinyipa-stable-rocky.gz",
			"ipmi_password":  "password",
		},
	}).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("created node: %v\n", createNode)

	// Example to Create a Port
	createPort, err := ports.Create(client, ports.CreateOpts{
		NodeUUID:        createNode.UUID,
		Address:         "00:73:49:3a:76:8e",
		PhysicalNetwork: "provisioning",
	}).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("created port: %v\n", createPort)

	// Example to Update Node
	updateNode, err := nodes.Update(client, createNode.UUID,
		nodes.UpdateOpts{
			nodes.UpdateOperation{
				Op:    nodes.AddOp,
				Path:  "/instance_info/image_source",
				Value: "http://172.22.0.1/images/redhat-coreos-maipo-latest.qcow2",
			},
			nodes.UpdateOperation{
				Op:    nodes.AddOp,
				Path:  "/instance_info/image_checksum",
				Value: "97830b21ed272a3d854615beb54cf004",
			},
		}).Extract()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		panic(err)
	}
	fmt.Printf("updated node with image source and checksum: %v\n", updateNode)

	validateResult, err := nodes.Validate(client, createNode.UUID).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("validation results:\n")
	fmt.Printf("\tboot: %v\n", validateResult.Boot)
	fmt.Printf("\tconsole: %v\n", validateResult.Console)
	fmt.Printf("\tdeploy: %v\n", validateResult.Deploy)
	fmt.Printf("\tinspect: %v\n", validateResult.Inspect)
	fmt.Printf("\tmanagement: %v\n", validateResult.Management)
	fmt.Printf("\tnetwork: %v\n", validateResult.Network)
	fmt.Printf("\tpower: %v\n", validateResult.Power)
	fmt.Printf("\traid: %v\n", validateResult.RAID)
	fmt.Printf("\trescue: %v\n", validateResult.Rescue)
	fmt.Printf("\tstorage: %v\n", validateResult.Storage)

	fmt.Printf("\n** Setting node Manageable **\n\n%v", validateResult)
	changeResult := nodes.ChangeProvisionState(client, createNode.UUID,
		nodes.ProvisionStateOpts{
			Target: nodes.TargetManage,
		})
	fmt.Printf("ChangeProvisionState requested, result:%v\n", changeResult)

	err = waitForState(client, 60, createNode.UUID, nodes.Manageable)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Done!\n")
}
