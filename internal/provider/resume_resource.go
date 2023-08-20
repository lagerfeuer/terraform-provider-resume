package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"net/http"
	"strconv"
)

var (
	_ resource.Resource                = &resumeResource{}
	_ resource.ResourceWithConfigure   = &resumeResource{}
	_ resource.ResourceWithImportState = &resumeResource{}
)

var resumeEndpoint = "/resumes"

func NewResumeResource() resource.Resource {
	return &resumeResource{}
}

type resumeResource struct {
	client *client
}

type resumeResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Address     types.String `tfsdk:"address"`
	PhoneNumber types.String `tfsdk:"phone_number"`
	Website     types.String `tfsdk:"website"`
}

type resumeResourceJson struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Website     string `json:"website"`
}

func (r *resumeResource) Configure(
	_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse,
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
	r.client = client
}

func (r *resumeResource) Metadata(
	ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_resume"
}

func (r *resumeResource) Schema(
	ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"address": schema.StringAttribute{
				Computed: false,
				Optional: true,
			},
			"phone_number": schema.StringAttribute{
				Computed: false,
				Optional: true,
			},
			"website": schema.StringAttribute{
				Computed: false,
				Optional: true,
			},
		},
	}
}

func (r *resumeResource) Create(
	ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse,
) {
	var plan resumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data resumeResourceJson
	data.Name = plan.Name.ValueString()
	data.Address = plan.Address.ValueString()
	data.PhoneNumber = plan.PhoneNumber.ValueString()
	data.Website = plan.Website.ValueString()

	reqBodyBytes, err := json.Marshal(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not convert JSON to request body",
			err.Error(),
		)
		return
	}

	reqBody := bytes.NewReader(reqBodyBytes)
	if reqBody == nil {
		resp.Diagnostics.AddError(
			"Could not create reader from bytes",
			err.Error(),
		)
		return
	}

	httpResp, err := r.client.Post(ctx, resumeEndpoint, reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Resume",
			err.Error(),
		)
		return
	}
	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError(
			"Error updating Resume",
			fmt.Sprintf("Expected 201, got %d.", httpResp.StatusCode),
		)
		return
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read HTTP response body",
			err.Error(),
		)
		return
	}

	if err = json.Unmarshal(body, &data); err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse HTTP response body",
			err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(strconv.FormatInt(data.Id, 10))
	plan.Name = types.StringValue(data.Name)
	if data.Address == "" {
		plan.Address = types.StringNull()
	} else {
		plan.Address = types.StringValue(data.Address)
	}
	if data.PhoneNumber == "" {
		plan.PhoneNumber = types.StringNull()
	} else {
		plan.PhoneNumber = types.StringValue(data.PhoneNumber)
	}
	if data.Website == "" {
		plan.Website = types.StringNull()
	} else {
		plan.Website = types.StringValue(data.Website)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resumeResource) Read(
	ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse,
) {
	var state resumeResourceModel
	var data resumeResourceJson

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/%s", resumeEndpoint, state.Id.ValueString())
	httpResp, err := r.client.Get(ctx, url)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Resume",
			err.Error(),
		)
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error updating Resume",
			fmt.Sprintf("Expected 200, got %d.", httpResp.StatusCode),
		)
		return
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read HTTP response body",
			err.Error(),
		)
		return
	}

	if err = json.Unmarshal(body, &data); err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse HTTP response body",
			err.Error(),
		)
		return
	}

	state.Id = types.StringValue(strconv.FormatInt(data.Id, 10))
	state.Name = types.StringValue(data.Name)
	if data.Address == "" {
		state.Address = types.StringNull()
	} else {
		state.Address = types.StringValue(data.Address)
	}
	if data.PhoneNumber == "" {
		state.PhoneNumber = types.StringNull()
	} else {
		state.PhoneNumber = types.StringValue(data.PhoneNumber)
	}
	if data.Website == "" {
		state.Website = types.StringNull()
	} else {
		state.Website = types.StringValue(data.Website)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resumeResource) Update(
	ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse,
) {
	var plan resumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var data resumeResourceJson
	data.Name = plan.Name.ValueString()
	data.Address = plan.Address.ValueString()
	data.PhoneNumber = plan.PhoneNumber.ValueString()
	data.Website = plan.Website.ValueString()

	reqBodyBytes, err := json.Marshal(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not convert JSON to request body",
			err.Error(),
		)
		return
	}

	reqBody := bytes.NewReader(reqBodyBytes)
	if reqBody == nil {
		resp.Diagnostics.AddError(
			"Could not create reader from bytes",
			err.Error(),
		)
		return
	}

	url := fmt.Sprintf("%s/%s", resumeEndpoint, plan.Id.ValueString())
	httpResp, err := r.client.Patch(ctx, url, reqBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Resume",
			err.Error(),
		)
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error updating Resume",
			fmt.Sprintf("Expected 200, got %d.", httpResp.StatusCode),
		)
		return
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read HTTP response body",
			err.Error(),
		)
		return
	}

	if err = json.Unmarshal(body, &data); err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse HTTP response body",
			err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(strconv.FormatInt(data.Id, 10))
	plan.Name = types.StringValue(data.Name)
	if data.Address == "" {
		plan.Address = types.StringNull()
	} else {
		plan.Address = types.StringValue(data.Address)
	}
	if data.PhoneNumber == "" {
		plan.PhoneNumber = types.StringNull()
	} else {
		plan.PhoneNumber = types.StringValue(data.PhoneNumber)
	}
	if data.Website == "" {
		plan.Website = types.StringNull()
	} else {
		plan.Website = types.StringValue(data.Website)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resumeResource) Delete(
	ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse,
) {
	var state resumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := fmt.Sprintf("%s/%s", resumeEndpoint, state.Id.ValueString())
	httpResp, err := r.client.Delete(ctx, url)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Resume",
			err.Error(),
		)
		return
	}
	if httpResp.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error deleted Resume",
			fmt.Sprintf("Expected %d, got %d.", http.StatusNoContent, httpResp.StatusCode),
		)
		return
	}
}

func (r *resumeResource) ImportState(
	ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
