package openstack

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

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
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAmphoraeVrrpIpV2Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	lbClient, err := config.loadBalancerV2Client(GetRegion(d, config))

	listOpts := amphorae.ListOpts{}

	if v, ok := d.GetOk("loadbalancer_id"); ok {
		listOpts.LoadbalancerID = v.(string)
	}

	pages, err := amphorae.List(lbClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to list amphorea: %s", err)
	}

	lbs, err := amphorae.ExtractAmphorae(pages)

	if err != nil {
		return fmt.Errorf("Unable to retrieve amphorea: %s", err)
	}

	if len(lbs) == 0 {
		return fmt.Errorf("No amphorea found")
	}

	var allIps []string

	for _, lb := range lbs {
		allIps = append(allIps, lb.VRRPPortID)
	}

	d.SetId(listOpts.LoadbalancerID)

	d.Set("ips", allIps)

	return nil
}
