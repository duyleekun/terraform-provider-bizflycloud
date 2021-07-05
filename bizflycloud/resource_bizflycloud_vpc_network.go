package bizflycloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/bizflycloud/gobizfly"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBizFlyCloudVPCNetwork() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBizFlyCloudVPCNetworkCreate,
		Read:          resourceBizFlyCloudVPCNetworkRead,
		Update:        resourceBizFlyCloudVPCNetworkUpdate,
		Delete:        resourceBizFlyCloudVPCNetworkDelete,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceBizFlyCloudVPCNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	log.Println("[DEBUG] creating vpc network")
	cvp := &gobizfly.CreateVPCPayload{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CIDR:        d.Get("cidr").(string),
		IsDefault:   d.Get("is_default").(bool),
	}
	log.Printf("[DEBUG] Create vpc network configuration: %#v\n", cvp)
	network, err := client.VPC.Create(context.Background(), cvp)
	if err != nil {
		return fmt.Errorf("Error creating vpc network: %v", err)
	}
	log.Println("[DEBUG] set id " + network.ID)
	d.SetId(network.ID)
	return resourceBizFlyCloudVPCNetworkRead(d, meta)
}

func resourceBizFlyCloudVPCNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	log.Printf("[DEBUG] vpc network ID %s", d.Id())
	network, err := client.VPC.Get(context.Background(), d.Id())
	if err != nil {
		if errors.Is(err, gobizfly.ErrNotFound) {
			log.Printf("[WARN] vpc network id %s is not found", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieved vpc network: %v", err)
	}
	_ = d.Set("name", network.Name)
	_ = d.Set("description", network.Description)
	_ = d.Set("is_default", network.IsDefault)
	return nil
}

func resourceBizFlyCloudVPCNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	log.Println("[DEBUG] update vpc network")
	var uvp gobizfly.UpdateVPCPayload
	networkChanged := false
	// uvp := &gobizfly.UpdateVPCPayload{
	// 	Name:        d.Get("name").(string),
	// 	Description: d.Get("description").(string),
	// 	CIDR:        d.Get("cidr").(string),
	// 	IsDefault:   d.Get("is_default").(bool),
	// }
	if d.HasChange("name") {
		networkChanged = true
		name := d.Get("name").(string)
		uvp.Name = name
	}
	if d.HasChange("description") {
		networkChanged = true
		desc := d.Get("description").(string)
		uvp.Description = desc
	}
	if networkChanged {
		_, err := client.VPC.Update(context.Background(), d.Id(), &uvp)
		if err != nil {
			return fmt.Errorf("Error when update vpc network: %s, %v", d.Id(), err)
		}
	}
	return nil
}

func resourceBizFlyCloudVPCNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).gobizflyClient()
	err := client.VPC.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error delete vpc network: %v", err)
	}
	return nil
}
