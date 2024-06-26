package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/myscribae/myscribae-sdk-go/provider"
	"github.com/myscribae/myscribae-sdk-go/utilities"
)

var _ datasource.DataSource = (*scriptDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*scriptDataSource)(nil)

type scriptDataSource struct {
	terraformProvider *myScribaeProvider
	myscribaeProvider *provider.Provider
	script            *provider.Script
}

func newScriptDataSource() datasource.DataSource {
	return &scriptDataSource{}
}

func (e *scriptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "myscribae_script"
}

func (e *scriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	prov, ok := req.ProviderData.(*myScribaeProvider)
	if !ok {
		resp.Diagnostics.AddError("invalid provider data", "expected *myScribaeProvider")
		return
	}

	e.terraformProvider = prov
}

func (e *scriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the script",
				Computed:    true,
				Required:    false,
			},
			"provider_id": schema.StringAttribute{
				Description: "The provider id of the script",
				Required:    true,
			},
			"script_group_id": schema.StringAttribute{
				Description: "The script group id",
				Required:    true,
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the script",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the script",
				Computed:    true,
			},
			"recurrence": schema.StringAttribute{
				Description: "The recurrence of the script",
				Computed:    true,
			},
			"price_in_cents": schema.Int64Attribute{
				Description: "The price in cents of the script",
				Computed:    true,
			},
			"sla_sec": schema.Int64Attribute{
				Description: "The sla in seconds of the script",
				Computed:    true,
			},
			"token_lifetime_sec": schema.Int64Attribute{
				Description: "The token lifetime in seconds of the script",
				Computed:    true,
			},
			"public": schema.BoolAttribute{
				Description: "Is the script public",
				Computed:    true,
			},
		},
	}
}

func (e *scriptDataSource) MakeClient(ctx context.Context, providerId string, scriptGroupId string, altId string) error {
	providerUuid, err := uuid.Parse(providerId)
	if err != nil {
		return err
	}
	e.myscribaeProvider = &provider.Provider{
		Uuid:   providerUuid,
		Client: e.terraformProvider.Client,
	}

	scriptGroupAltID, err := utilities.NewAltUuid(scriptGroupId)
	if err != nil {
		return err
	}

	e.script, err = e.myscribaeProvider.Script(scriptGroupAltID, altId)

	return err
}

func (e *scriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &scriptResourceData{}
	if diags := req.Config.Get(ctx, data); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, data.ProviderID.ValueString(), data.ScriptGroupID.ValueString(), data.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError("error making client", err.Error())
		return
	}

	profile, err := e.script.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError("error reading script", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &scriptResourceData{
		Id:               basetypes.NewStringValue(profile.Uuid.String()),
		Uuid:             basetypes.NewStringValue(profile.Uuid.String()),
		ProviderID:       data.ProviderID,
		ScriptGroupID:    data.ScriptGroupID,
		AltID:            basetypes.NewStringValue(profile.AltID),
		Name:             basetypes.NewStringValue(profile.Name),
		Description:      basetypes.NewStringValue(profile.Description),
		Recurrence:       basetypes.NewStringValue(profile.Recurrence),
		PriceInCents:     basetypes.NewInt64Value(int64(profile.PriceInCents)),
		SlaSec:           basetypes.NewInt64Value(int64(profile.SlaSec)),
		TokenLifetimeSec: basetypes.NewInt64Value(int64(profile.TokenLifetimeSec)),
		Public:           basetypes.NewBoolValue(profile.Public),
	})
	resp.Diagnostics.Append(diags...)
}
