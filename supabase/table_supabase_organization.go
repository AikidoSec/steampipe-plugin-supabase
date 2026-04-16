package supabase

import (
	"context"
	"fmt"

	"github.com/supabase/cli/pkg/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableSupabaseOrganization(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "supabase_organization",
		Description: "Supabase Organization",
		List: &plugin.ListConfig{
			Hydrate: listSupabaseOrganizations,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the organization."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the organization."},
			{Name: "slug", Type: proto.ColumnType_STRING, Description: "The slug of the organization."},
			{Name: "plan", Type: proto.ColumnType_STRING, Description: "Pricing plan of the organization.", Hydrate: getSupabaseOrganizationDetails, Transform: transform.FromField("Plan")},
			{Name: "opt_in_tags", Type: proto.ColumnType_JSON, Description: "The opt-in tags.", Hydrate: getSupabaseOrganizationDetails, Transform: transform.FromField("OptInTags")},
			{Name: "allowed_release_channels", Type: proto.ColumnType_JSON, Description: "The allowed release channels.", Hydrate: getSupabaseOrganizationDetails, Transform: transform.FromField("AllowedReleaseChannels")},

			// Common steampipe columns
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue().Transform(generateOrganizationAKA).Transform(transform.EnsureStringArray),
			},
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

//// LIST FUNCTION

func listSupabaseOrganizations(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_organization.listSupabaseOrganizations", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1ListAllOrganizationsWithResponse(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_organization.listSupabaseOrganizations", "query_error", err)
		return nil, err
	}

	for _, organization := range *resp.JSON200 {
		d.StreamListItem(ctx, organization)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getSupabaseOrganizationDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	orgSlug := h.Item.(api.OrganizationResponseV1).Slug
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_organization.getSupabaseOrganizationDetails", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetAnOrganizationWithResponse(ctx, orgSlug)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_organization.getSupabaseOrganizationDetails", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func generateOrganizationAKA(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	obj, ok := d.Value.(api.OrganizationResponseV1)
	if !ok {
		return nil, fmt.Errorf("could not cast value to transform")
	}
	return fmt.Sprintf("org:%s", obj.Id), nil
}
