package openstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/amphorae"
)


func dataSourceAmphoraeVrrpIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAmphoraeVrrpIpV2Read,

		Schema: map[string]*schema.Schema{
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ips": {
			    Type:   schema.TypeList,
			    Computed:   true,
			    Elem:   &schema.Schema{Type: schema.TypeString},
			}
		},
	}
}

func dataSourceAmphoraeVrrpIpV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	
	listOpts := amphorae.ListOpts{}

	if v, ok := d.GetOk("loadbalancer_id"); ok {
		listOpts.LoadbalancerID = v.(string)
	}


	allPages, err := ports.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to list Ports: %s", err)
	}

	var allLbs []Amphora

	err = amphorae.ExtractAmphorae(allPages, &allLbs)
	if err != nil {
		return fmt.Errorf("Unable to retrieve amphorea: %s", err)
	}

	if len(allLbs) == 0 {
		return fmt.Errorf("No amphorea found")
	}

    var allIps []String
    
    for _, ip := range allLbs {
        allIps = append(allIps, ip)
    }

	d.SetId(listOpts.LoadbalancerID)

	d.Set("ips", allIps)
	
	return nil
}
