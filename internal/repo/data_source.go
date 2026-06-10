package repo

import (
	"context"

	cobbler "github.com/cobbler/cobblerclient"
	clientpkg "github.com/cobbler/terraform-provider-cobbler/internal/client"
	"github.com/cobbler/terraform-provider-cobbler/internal/inherit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &RepoDataSource{}

type RepoDataSource struct {
	client cobbler.Client
}

func NewDataSource() datasource.DataSource {
	return &RepoDataSource{}
}

func (d *RepoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repo"
}

func (d *RepoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get the details of a Cobbler repo.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the repo.",
				Required:    true,
			},
			"arch": schema.StringAttribute{
				Description: "The architecture of the repo.",
				Computed:    true,
			},
			"breed": schema.StringAttribute{
				Description: "The \"breed\" of distribution.",
				Computed:    true,
			},
			"comment": schema.StringAttribute{
				Description: "Free form text description.",
				Computed:    true,
			},
			"environment": schema.MapAttribute{
				Description: "Environment variables to use during repo command execution.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"keep_updated": schema.BoolAttribute{
				Description: "Update the repo upon Cobbler sync.",
				Computed:    true,
			},
			"mirror": schema.StringAttribute{
				Description: "Address of the repo to mirror.",
				Computed:    true,
			},
			"mirror_locally": schema.BoolAttribute{
				Description: "Whether to copy the files locally or just references to the external files.",
				Computed:    true,
			},
			"apt_components": schema.ListAttribute{
				Description: "List of Apt components.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"apt_dists": schema.ListAttribute{
				Description: "List of Apt distribution names.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"rpm_list": schema.ListAttribute{
				Description: "List of specific RPMs to mirror.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"createrepo_flags": schema.SingleNestedAttribute{
				Description: "Flags to use with `createrepo`.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Computed: true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"owners": schema.SingleNestedAttribute{
				Description: "List of Owners for authz_ownership.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"proxy": schema.SingleNestedAttribute{
				Description: "Proxy to use for downloading the repo.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"value": schema.StringAttribute{
						Computed: true,
					},
					"inherited": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
		},
	}
}

func (d *RepoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	cfg, ok := req.ProviderData.(*clientpkg.Config)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			"Expected *client.Config, got unexpected type.")
		return
	}
	d.client = cfg.CobblerClient
}

func (d *RepoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data repoDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	repo, err := d.client.GetRepo(data.Name.ValueString(), false, false)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Cobbler Repo", err.Error())
		return
	}

	data.Name = types.StringValue(repo.Name)
	data.Arch = types.StringValue(repo.Arch)
	data.Breed = types.StringValue(repo.Breed)
	data.Comment = types.StringValue(repo.Comment)
	data.KeepUpdated = types.BoolValue(repo.KeepUpdated)
	data.Mirror = types.StringValue(repo.Mirror)
	data.MirrorLocally = types.BoolValue(repo.MirrorLocally)
	data.CreateRepoFlags = inherit.StringFrom(ctx, repo.CreateRepoFlags, &resp.Diagnostics)
	data.Owners = inherit.StringListFrom(ctx, repo.Owners, &resp.Diagnostics)
	data.Proxy = inherit.StringFrom(ctx, repo.Proxy, &resp.Diagnostics)

	env, d2 := types.MapValueFrom(ctx, types.StringType, repo.Environment)
	resp.Diagnostics.Append(d2...)
	data.Environment = env

	aptComponents, d2 := types.ListValueFrom(ctx, types.StringType, repo.AptComponents)
	resp.Diagnostics.Append(d2...)
	data.AptComponents = aptComponents

	aptDists, d2 := types.ListValueFrom(ctx, types.StringType, repo.AptDists)
	resp.Diagnostics.Append(d2...)
	data.AptDists = aptDists

	rpmList, d2 := types.ListValueFrom(ctx, types.StringType, repo.RpmList)
	resp.Diagnostics.Append(d2...)
	data.RpmList = rpmList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
