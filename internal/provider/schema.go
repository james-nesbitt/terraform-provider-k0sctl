package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"	
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/diag" 
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	k0sctl_v1beta1 "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1"
	k0sctl_v1beta1_cluster "github.com/k0sproject/k0sctl/pkg/apis/k0sctl.k0sproject.io/v1beta1/cluster"

	k0s_rig "github.com/k0sproject/rig"
)

const (
	k0sctl_schema_kind = "cluster"
)

func k0sctl_schema(ctx context.Context) (schema.Schema, diag.Diagnostics) {
	return schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Mirantis installation using launchpad, parametrized",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"skip_destroy": schema.BoolAttribute{
				MarkdownDescription: "Skip reset on destroy",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"skip_create": schema.BoolAttribute{
				MarkdownDescription: "Skip apply on create",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},


		Blocks: map[string]schema.Block{

			"metadata": schema.SingleNestedBlock{
				MarkdownDescription: "Metadata for the launchpad cluster",

				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Cluster name",
						Required:            true,
					},
				},
			},

			"spec": schema.SingleNestedBlock{
				MarkdownDescription: "Launchpad install specifications",

				Blocks: map[string]schema.Block{

					"k0s": schema.SingleNestedBlock{
						MarkdownDescription: "K0S installation configuration",

						Attributes: map[string]schema.Attribute{
							"version": schema.StringAttribute{
								MarkdownDescription: "K0s version to install",
								Required:            true,
							},
							"channel": schema.StringAttribute{
								MarkdownDescription: "Repository installation channel",
								Optional:            true,
								Computed:            true,
								Default:             stringdefault.StaticString("stable"),
							},

							// Needs: DynamicConfig:bool Config:dig.Mapping 
						},
					},


					"host": schema.ListNestedBlock{
						MarkdownDescription: "Individual host configuration, for each machine in the cluster",

						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},

						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"role": schema.StringAttribute{
									MarkdownDescription: "Host machine role in the cluster",
									Required:            true,
								},
							},
							Blocks: map[string]schema.Block{

								"hooks": schema.ListNestedBlock{
									MarkdownDescription: "Hook configuration for the host",

									Validators: []validator.List{
										listvalidator.SizeAtMost(1),
									},

									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{},
										Blocks: map[string]schema.Block{

											"apply": schema.ListNestedBlock{
												MarkdownDescription: "Launchpad.Apply string hooks for the host",

												Validators: []validator.List{
													listvalidator.SizeAtMost(1),
												},

												NestedObject: schema.NestedBlockObject{
													Attributes: map[string]schema.Attribute{
														"before": schema.ListAttribute{
															MarkdownDescription: "String hooks to run on hosts before the Apply operation is run.",
															ElementType:         types.StringType,
															Optional:            true,
															Computed:            true,
															Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
														},
														"after": schema.ListAttribute{
															MarkdownDescription: "String hooks to run on hosts after the Apply operation is run.",
															ElementType:         types.StringType,
															Optional:            true,
															Computed:            true,
															Default:             listdefault.StaticValue(types.ListNull(types.StringType)),
														},
													},
												},
											},
										},
									},
								},

								"ssh": schema.ListNestedBlock{
									MarkdownDescription: "SSH configuration for the host",

									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"address": schema.StringAttribute{
												MarkdownDescription: "SSH endpoint",
												Required:            true,
											},
											"key_path": schema.StringAttribute{
												MarkdownDescription: "SSH endpoint",
												Required:            true,
											},
											"user": schema.StringAttribute{
												MarkdownDescription: "SSH endpoint",
												Required:            true,
											},
											"port": schema.Int64Attribute{
												MarkdownDescription: "SSH Port",
												Optional:            true,
												Computed:            true,
												Default:             int64default.StaticInt64(22),
											},
										},
									},
								},
								"winrm": schema.ListNestedBlock{
									MarkdownDescription: "WinRM configuration for the host",

									NestedObject: schema.NestedBlockObject{
										Attributes: map[string]schema.Attribute{
											"address": schema.StringAttribute{
												MarkdownDescription: "WinRM endpoint",
												Required:            true,
											},
											"user": schema.StringAttribute{
												MarkdownDescription: "WinRM user",
												Required:            true,
											},
											"password": schema.StringAttribute{
												MarkdownDescription: "WinRM password",
												Required:            true,
											},
											"port": schema.Int64Attribute{
												MarkdownDescription: "WinRM Port",
												Optional:            true,
												Computed:            true,
												Default:             int64default.StaticInt64(5985),
											},
											"use_https": schema.BoolAttribute{
												MarkdownDescription: "If false, then no HTTP is used for winrm transport",
												Optional:            true,
												Computed:            true,
												Default:             booldefault.StaticBool(true),
											},
											"insecure": schema.BoolAttribute{
												MarkdownDescription: "If false, then no SSL certificate validation is used",
												Optional:            true,
												Computed:            true,
												Default:             booldefault.StaticBool(true),
											},
										},
									},
								},
							},
						},
					},
				},

			},
		},
	
	}, nil
}

type k0sctlSchemaModel struct {
	Id          types.String `tfsdk:"id"`
	SkipCreate  types.Bool   `tfsdk:"skip_create"`
	SkipDestroy types.Bool   `tfsdk:"skip_destroy"`

	Metadata    k0sctlSchemaClusterMetadata `tfsdk:"metadata"`
	Spec k0sctlSchemaModelSpec `tfsdk:"spec"`
}

func (ksm k0sctlSchemaModel) Cluster(ctx context.Context) (k0sctl_v1beta1.Cluster, diag.Diagnostics) {
	tflog.Info(ctx, "Creating k0sctl Cluster from schema", map[string]interface{}{})
	
	var c k0sctl_v1beta1.Cluster
	var d diag.Diagnostics

	c = k0sctl_v1beta1.Cluster{
		APIVersion: k0sctl_v1beta1.APIVersion,
		Kind: k0sctl_schema_kind,

		Metadata: &k0sctl_v1beta1.ClusterMetadata{
			Name: ksm.Metadata.Name.ValueString(),
		},

		Spec: &k0sctl_v1beta1_cluster.Spec{
			Hosts: k0sctl_v1beta1_cluster.Hosts{},
			K0s: &k0sctl_v1beta1_cluster.K0s{
				Version: ksm.Spec.K0s.Version.ValueString(),
				Channel: ksm.Spec.K0s.Channel.ValueString(),
			},
		},
	}

	for _, sh := range ls.Spec.Hosts {
		h := k0sctl_v1beta1_cluster.Host{
			Role: sh.Role.ValueString(),
			Hooks: k0sctl_v1beta1_cluster.Hooks{},
		}

		if len(host.SSH) > 0 {
			shssh := sh.SSH[0]

			mccHost.Connection = k0s_rig.Connection{
				SSH: &k0s_rig.SSH{
					Address: shssh.Address.ValueString(),
					KeyPath: shssh.KeyPath.ValueStringPointer(),
					User:    shssh.User.ValueString(),
					Port:    int(shssh.Port.ValueInt64()),
				},
			}
		} else if len(host.WinRM) > 0 {
			shwinrm := host.WinRM[0]

			mccHost.Connection = k0s_rig.Connection{
				WinRM: &k0s_rig.WinRM{
					Address:  shwinrm.Address.ValueString(),
					Password: shwinrm.Password.ValueString(),
					User:     shwinrm.User.ValueString(),
					Port:     int(shwinrm.Port.ValueInt64()),
					UseHTTPS: shwinrm.UseHTTPS.ValueBool(),
					Insecure: shwinrm.Insecure.ValueBool(),
				},
			}
		}

		if len(sh.Hooks) > 0 {
			shh := sh.Hooks[0]

			if len(sh.Apply) > 0 {
				ha := shh.Apply[0]

				hha := map[string][]string{
					"before": {},
					"after":  {},
				}
				var shab []string
				if diag := ha.Before.ElementsAs(context.Background(), &shab, true); diag == nil {
					hha["before"] = shab
				}
				var shaa []string
				if diag := ha.After.ElementsAs(context.Background(), &shaa, true); diag == nil {
					hha["after"] = shab
				}

				h.Hooks["apply"] = hha
			}
		}

		c.Spec.Hosts = append(c.Spec.Hosts, h)
	}
	
	return c, d
}

type k0sctlSchemaClusterMetadata struct {
	Name types.String `tfsdk:"name"`
}

type k0sctlSchemaModelSpec struct {
	Hosts   []k0sctlSchemaModelSpecHost `tfsdk:"host"`
	K0s     k0sctlSchemaModelSpecK0s    `tfsdk:"k0s"`
}

type k0sctlSchemaModelSpecK0s struct {
	Version           types.String `tfsdk:"version"`
	Channel           types.String `tfsdk:"channel"`
}

type k0sctlSchemaModelSpecHost struct {
	Role      types.String                              `tfsdk:"role"`
	Hooks     []k0sctlSchemaModelSpecHostHooks     `tfsdk:"hooks"`
	SSH       []k0sctlSchemaModelSpecHostSSH       `tfsdk:"ssh"`
	WinRM     []k0sctlSchemaModelSpecHostWinrm     `tfsdk:"winrm"`
}
type k0sctlSchemaModelSpecHostHooks struct {
	Apply []k0sctlSchemaModelSpecHostHookAction `tfsdk:"apply"`
}
type k0sctlSchemaModelSpecHostMCRconfigDefaultAddressPools struct {
	Base types.String `json:"base" tfsdk:"base"`
	Size types.Int64  `json:"size" tfsdk:"size"`
}
type k0sctlSchemaModelSpecHostHookAction struct {
	Before types.List `tfsdk:"before"`
	After  types.List `tfsdk:"after"`
}
type k0sctlSchemaModelSpecHostSSH struct {
	Address types.String `tfsdk:"address"`
	KeyPath types.String `tfsdk:"key_path"`
	User    types.String `tfsdk:"user"`
	Port    types.Int64  `tfsdk:"port"`
}
type k0sctlSchemaModelSpecHostWinrm struct {
	Address  types.String `tfsdk:"address"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Port     types.Int64  `tfsdk:"port"`
	UseHTTPS types.Bool   `tfsdk:"use_https"`
	Insecure types.Bool   `tfsdk:"insecure"`
}
