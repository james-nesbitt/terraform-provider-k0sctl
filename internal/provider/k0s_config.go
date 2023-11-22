package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &K0sctlConfigResource{}

type K0sctlConfigResource struct {
	testingMode bool
}

func NewK0sctlConfigResource() resource.Resource {
	return &K0sctlConfigResource{}
}

func (r *K0sctlConfigResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	tflog.Info(ctx, "k0sctl metadata run", map[string]interface{}{})
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (r *K0sctlConfigResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema, diags := k0sctl_schema(ctx)
	resp.Schema = schema
	if diags != nil {
		resp.Diagnostics.Append(diags...)
	}
}

func (r *K0sctlConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

}

func (r *K0sctlConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

func (r *K0sctlConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// k0sctl has no good way to discover existing installation, so we don't do anything
}

func (r *K0sctlConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

func (r *K0sctlConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *K0sctlConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

}
