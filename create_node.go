package main

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/baremetal/noauth"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
)

/* totally untested.  Just a theory, feel free to delete.

func wait_for_state(client *gophercloudServiceClient, timeout int, state string) (r err) {

	i := 0
	for i < timeout {
		i += 1
		// FIXME: I know we can do a search for the specific UUID.
		nodes.ListDetail(client, nodes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
			nodeList, err := nodes.ExtractNodes(page)
			if err != nil {
				return false, err
			}

			for _, n := range nodeList {
				if n.UUID == UUID {
					if n.ProvisionState == state {
						return nil
					}
				}

				fmt.Printf("%s %s %s %s - looking for %s\n", n.UUID, n.Name, n.PowerState, n.ProvisionState, state)
			}

			return true, nil
		})
	}
}

*/

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
		Address:         "00:cd:18:18:77:ca",
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

	fmt.Printf("\n** Setting node Manageable **\n\n", validateResult)
	changeResult := nodes.ChangeProvisionState(client, createNode.UUID,
		nodes.ProvisionStateOpts{
			Target: nodes.TargetManage,
		})
	fmt.Printf("ChangeProvisionState requested, result:%v\n", changeResult)

	//wait_for_state(client, 60, node.Manageable)
}
