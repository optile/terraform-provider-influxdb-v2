package influxdbv2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/lancey-energy-storage/influxdb-client-go"
)

func ResourceBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceBucketCreate,
		Delete: resourceBucketDelete,
		Read:   resourceBucketRead,
		Update: resourceBucketUpdate,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"org_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retention_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"every_seconds": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: false,
							Default:  "expire",
						},
					},
				},
			},
			"rp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBucketCreate(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(*influxdb.Client)

	if d.Get("name") == "" {
		return errors.New("a name is required")
	}

	rr := d.Get("retention_rules")
	retentionRules, err := SetRetentionRules(rr)
	if err != nil {
		return err
	}

	result, err := influx.CreateBucket(d.Get("description").(string), d.Get("name").(string), d.Get("org_id").(string), retentionRules, d.Get("rp").(string))
	if err != nil {
		return err
	}

	return resourceBucketRead(d, meta)
}

func resourceCreateBucketDelete(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(*influxdb.Client)
	err := influx.DeleteABucket(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceBucketRead(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(*influxdb.Client)
	result, err := influx.GetBucketByID(d.Id())
	if err != nil {
		return fmt.Errorf("error getting bucket: %v", err)
	}
	d.Set("type", result.Type)
	d.Set("created_at", result.CreatedAt)
	d.Set("updated_at", result.UpdatedAt)

	d.SetId(result.Id)
	return nil
}

func resourceCreateBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	influx := meta.(*influxdb.Client)

	if d.Get("name") == "" {
		return errors.New("a name is required")
	}

	retentionRules, err := SetRetentionRules(d.Get("retention_rules"))
	if err != nil {
		return err
	}

	labels, err := SetLabels(d.Get("labels"))
	if err != nil {
		return err
	}

	_, err = influx.UpdateABucket(d.Id(), d.Get("description").(string), labels, d.Get("name").(string), d.Get("org_id").(string), retentionRules, d.Get("rp").(string))
	if err != nil {
		return err
	}

	return nil
}
