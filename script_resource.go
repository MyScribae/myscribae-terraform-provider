package main

import (
	"context"

	sdk "github.com/Pritch009/myscribae-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = (*scriptResource)(nil)

type scriptResource struct {
	provider myScribaeProvider
}

type scriptResourceData struct {
	Uuid             types.String `tfsdk:"uuid"`
	ScriptGroupUuid  types.String `tfsdk:"script_group_uuid"`
	AltID            types.String `tfsdk:"alt_id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Recurrence       types.String `tfsdk:"recurrence"`
	PriceInCents     types.Int64  `tfsdk:"price_in_cents"`
	SlaSec           types.Int64  `tfsdk:"sla_sec"`
	TokenLifetimeSec types.Int64  `tfsdk:"token_lifetime_sec"`
	Public           types.Bool   `tfsdk:"public"`
}

func newScriptResource() resource.Resource {
	return &scriptResource{}
}

func (e *scriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "script"
}

func (e *scriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"script_group_uuid": schema.StringAttribute{
				Description: "The script group uuid",
				Required:    true,
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the script",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the script",
				Required:    true,
			},
			"recurrence": schema.StringAttribute{
				Description: "The recurrence of the script",
				Required:    true,
			},
			"price_in_cents": schema.NumberAttribute{
				Description: "The price in cents of the script",
				Required:    true,
			},
			"sla_sec": schema.NumberAttribute{
				Description: "The SLA in seconds of the script",
				Required:    true,
			},
			"token_lifetime_sec": schema.NumberAttribute{
				Description: "The token lifetime in seconds of the script",
				Required:    true,
			},
			"public": schema.BoolAttribute{
				Description: "Is the script public",
				Required:    true,
			},
		},
	}
}

func (e *scriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := scriptResourceData{}
	diags := req.Plan.Get(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	s := e.provider.Client.Script(
		data.ScriptGroupUuid.ValueString(),
		data.AltID.ValueString(),
	)
	resultUuid, err := s.Create(ctx, sdk.ScriptInput{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		Recurrence:       data.Recurrence.ValueString(),
		PriceInCents:     int(data.PriceInCents.ValueInt64()),
		SlaSec:           int(data.SlaSec.ValueInt64()),
		TokenLifetimeSec: int(data.TokenLifetimeSec.ValueInt64()),
		Public:           data.Public.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create script",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &scriptResourceData{
		Uuid:             basetypes.NewStringValue(resultUuid.String()),
		ScriptGroupUuid:  data.ScriptGroupUuid,
		AltID:            data.AltID,
		Name:             data.Name,
		Description:      data.Description,
		PriceInCents:     data.PriceInCents,
		SlaSec:           data.SlaSec,
		TokenLifetimeSec: data.TokenLifetimeSec,
		Public:           data.Public,
	})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := scriptResourceData{}
	diags := req.State.Get(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	s := e.provider.Client.Script(
		data.ScriptGroupUuid.ValueString(),
		data.AltID.ValueString(),
	)
	profile, err := s.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get script profile",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &scriptResourceData{
		Uuid:             basetypes.NewStringValue(profile.Uuid.String()),
		ScriptGroupUuid:  data.ScriptGroupUuid,
		AltID:            data.AltID,
		Name:             basetypes.NewStringValue(profile.Name),
		Description:      basetypes.NewStringValue(profile.Description),
		Recurrence:       basetypes.NewStringValue(profile.Recurrence),
		PriceInCents:     basetypes.NewInt64Value(int64(profile.PriceInCents)),
		SlaSec:           basetypes.NewInt64Value(int64(profile.SlaSec)),
		TokenLifetimeSec: basetypes.NewInt64Value(int64(profile.TokenLifetimeSec)),
		Public:           basetypes.NewBoolValue(profile.Public),
	})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := scriptResourceData{}
	diags := req.Plan.Get(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	s := e.provider.Client.Script(
		data.ScriptGroupUuid.ValueString(),
		data.AltID.ValueString(),
	)
	resultUuid, err := s.Update(ctx, sdk.ScriptInput{
		AltID:            data.AltID.ValueString(),
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		Recurrence:       data.Recurrence.ValueString(),
		PriceInCents:     int(data.PriceInCents.ValueInt64()),
		SlaSec:           int(data.SlaSec.ValueInt64()),
		TokenLifetimeSec: int(data.TokenLifetimeSec.ValueInt64()),
		Public:           data.Public.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update script",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &scriptResourceData{
		Uuid:             basetypes.NewStringValue(resultUuid.String()),
		ScriptGroupUuid:  data.ScriptGroupUuid,
		AltID:            data.AltID,
		Name:             data.Name,
		Description:      data.Description,
		PriceInCents:     data.PriceInCents,
		SlaSec:           data.SlaSec,
		TokenLifetimeSec: data.TokenLifetimeSec,
		Public:           data.Public,
	})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := scriptResourceData{}
	diags := req.State.Get(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	s := e.provider.Client.Script(
		data.ScriptGroupUuid.ValueString(),
		data.AltID.ValueString(),
	)
	err := s.Delete(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete script",
			err.Error(),
		)
		return
	}
}
