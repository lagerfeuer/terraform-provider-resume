package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"net/http"
)

var (
	_ datasource.DataSource              = &infoDataSource{}
	_ datasource.DataSourceWithConfigure = &infoDataSource{}
)

func NewInfoDataSource() datasource.DataSource {
	return &infoDataSource{}
}

type infoDataSource struct {
	client *client
}

type infoDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Version     types.String `tfsdk:"version"`
	Environment types.String `tfsdk:"environment"`
}

type infoDataSourceJson struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
}

func (d *infoDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *client, got: %T. Please report this to the provider developer.",
				req.ProviderData,
			),
		)
	}
	d.client = client
}

func (d *infoDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_info"
}

func (d *infoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Resume API instance.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The Resume API version.",
				Computed:    true,
			},
			"environment": schema.StringAttribute{
				Description: "The environment for the Resume API isntance.",
				Computed:    true,
			},
		},
	}
}

func (d *infoDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse,
) {
	var state infoDataSourceModel

	r, err := d.client.Get(ctx, "/info")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read /info",
			err.Error(),
		)
		return
	}

	if r.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"HTTP Error",
			fmt.Sprintf("Expected HTTP status code %d, but got: %d", http.StatusOK, r.StatusCode),
		)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read /info response",
			err.Error(),
		)
		return
	}

	var data infoDataSourceJson
	err = json.Unmarshal(body, &data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse /info response",
			err.Error(),
		)
		return
	}

	// TODO write custom Unmarshal function for infoDataSourceModel so it can read from JSON
	// See https://github.com/hashicorp/terraform-plugin-framework/issues/205
	// and https://developer.hashicorp.com/terraform/plugin/framework/handling-data/custom-types
	state.Id = types.StringValue("info")
	state.Name = types.StringValue(data.Name)
	state.Version = types.StringValue(data.Version)
	state.Environment = types.StringValue(data.Environment)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
