package resources

import (
	"context"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/client"
	"github.com/arvancloud/terraform-provider-arvan/internal/api/iaas"
	"github.com/arvancloud/terraform-provider-arvan/internal/service/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceVolumeAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeAttachCreate,
		ReadContext:   helper.DummyResourceAction,
		UpdateContext: helper.DummyResourceAction,
		DeleteContext: resourceVolumeAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "region code",
				ValidateFunc: validation.StringInSlice(iaas.AvailableRegions, false),
			},
			"abrak_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of abrak",
				ValidateFunc: validation.IsUUID,
			},
			"volume_uuid": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "UUID of volume",
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceVolumeAttachCreate(ctx context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "could not parse region",
		})
		return errors
	}

	// volumeAttachmentOpts Options
	volumeAttachmentOpts := &iaas.VolumeAttachmentOpts{
		ServerId: data.Get("abrak_uuid").(string),
		VolumeId: data.Get("volume_uuid").(string),
	}

	err := c.Volume.Attach(region, volumeAttachmentOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(volumeAttachmentOpts.VolumeId)
	return errors
}

func resourceVolumeAttachDelete(_ context.Context, data *schema.ResourceData, meta any) (errors diag.Diagnostics) {
	c := meta.(*client.Client).IaaS

	region, ok := data.Get("region").(string)
	if !ok {
		errors = append(errors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "can not get region",
		})
		return errors
	}

	// volumeAttachmentOpts Options
	volumeAttachmentOpts := &iaas.VolumeAttachmentOpts{
		ServerId: data.Get("abrak_uuid").(string),
		VolumeId: data.Get("volume_uuid").(string),
	}

	err := c.Volume.Detach(region, volumeAttachmentOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	return errors
}
