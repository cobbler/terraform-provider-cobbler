package repo

import "github.com/hashicorp/terraform-plugin-framework/types"

type repoResourceModel struct {
	Name            types.String `tfsdk:"name"`
	AptComponents   types.List   `tfsdk:"apt_components"`
	AptDists        types.List   `tfsdk:"apt_dists"`
	Arch            types.String `tfsdk:"arch"`
	Breed           types.String `tfsdk:"breed"`
	Comment         types.String `tfsdk:"comment"`
	Environment     types.Map    `tfsdk:"environment"`
	KeepUpdated     types.Bool   `tfsdk:"keep_updated"`
	Mirror          types.String `tfsdk:"mirror"`
	MirrorLocally   types.Bool   `tfsdk:"mirror_locally"`
	RpmList         types.List   `tfsdk:"rpm_list"`
	CreateRepoFlags types.Object `tfsdk:"createrepo_flags"`
	Owners          types.Object `tfsdk:"owners"`
	Proxy           types.Object `tfsdk:"proxy"`
}
