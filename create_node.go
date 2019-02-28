package main

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/baremetal/noauth"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
)

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
			"ipmi_password":  "admin",
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
		}).Extract()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		panic(err)
	}
	fmt.Printf("updated node with image_source: %v\n", updateNode)

	// Example to Update Node
	updateNode, err = nodes.Update(client, createNode.UUID,
		nodes.UpdateOpts{
			nodes.UpdateOperation{
				Op:    nodes.AddOp,
				Path:  "/instance_info/image_checksum",
				Value: "97830b21ed272a3d854615beb54cf004",
			},
		}).Extract()
	if err != nil {
		panic(err)
	}
	fmt.Printf("updated node with image_source: %v\n", updateNode)

	validateResult := nodes.Validate(client, createNode.UUID)
	fmt.Printf("validation returned: %v\n", validateResult)
}
