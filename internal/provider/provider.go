// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ResumeProvider satisfies various provider interfaces.
var _ provider.Provider = &ResumeProvider{}

// ResumeProvider defines the provider implementation.
type ResumeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ResumeProviderModel describes the provider data model.
type ResumeProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *ResumeProvider) Metadata(
	ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse,
) {
	resp.TypeName = "resume"
	resp.Version = p.version
}

func (p *ResumeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Resume API Endpoint",
				Required:            true,
			},
			//"token": schema.StringAttribute{
			//	MarkdownDescription: "Resume API Token",
			//	Optional:            false,
			//},
		},
	}
}

func (p *ResumeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ResumeProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Endpoint configuration",
			"The provider cannot create a client with an unknown value.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("RESUME_ENDPOINT")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Resume API Endpoint",
			"Cannot create the Resume client without a valid endpoint.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := newClient(endpoint, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not create client",
			fmt.Sprintf("Error while creating Resume API client: %v", err),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ResumeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ResumeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ResumeProvider{
			version: version,
		}
	}
}
