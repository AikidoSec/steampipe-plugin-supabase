package supabase

import (
	"context"
	"fmt"
	"reflect"

	"github.com/supabase/cli/pkg/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableSupabaseAPIKey(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "supabase_api_key",
		Description: "Supabase API key",
		List: &plugin.ListConfig{
			ParentHydrate: listSupabaseProjects,
			Hydrate:       listSupabaseProjectAPIKeys,
		},
		Columns: []*plugin.Column{
			{Name: "api_key", Type: proto.ColumnType_STRING, Description: "The API key.", Transform: transform.FromField("ApiKey").Transform(transformNullable)},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The id of the API key.", Transform: transform.FromField("Id").Transform(transformNullable)},
			{Name: "type", Type: proto.ColumnType_JSON, Description: "The type of the API key.", Transform: transform.FromField("Type").Transform(transformNullable)},
			{Name: "prefix", Type: proto.ColumnType_STRING, Description: "The API key prefix.", Transform: transform.FromField("Prefix").Transform(transformNullable)},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The API key name."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The API key description.", Transform: transform.FromField("Description").Transform(transformNullable)},
			{Name: "hash", Type: proto.ColumnType_STRING, Description: "The API key hash.", Transform: transform.FromField("Hash").Transform(transformNullable)},
			{Name: "secret_jwt_template", Type: proto.ColumnType_JSON, Description: "The API key jwt secret template.", Transform: transform.FromField("SecretJwtTemplate").Transform(transformNullable)},
			// {Name: "inserted_at", Type: proto.ColumnType_TIMESTAMP, Description: "The API key creation date.", Transform: transform.FromField("InsertedAt").Transform(transformNullable)},
			// {Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The API key last update date."},

			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The ID of the project where the function is located."},
			{
				Name:        "organization_slug",
				Type:        proto.ColumnType_STRING,
				Description: "The organization slug.",
				Hydrate:     getOrganizationSlugForProjectIDFromAPIKey,
				Transform:   transform.FromValue(),
			},

			// Common steampipe columns
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue().Transform(generateAPIKeyAKA).Transform(transform.EnsureStringArray),
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

type ProjectAPIKey struct {
	api.ApiKeyResponse
	ProjectId string
}

//// LIST FUNCTION

func listSupabaseProjectAPIKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the project data
	project := h.Item.(api.V1ProjectWithDatabaseResponse)

	// Create client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project_api_key.listSupabaseProjectAPIKeys", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetProjectApiKeysWithResponse(ctx, project.Id, &api.V1GetProjectApiKeysParams{})
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project_api_key.listSupabaseProjectAPIKeys", "query_error", err)
		return nil, err
	}

	for _, key := range *resp.JSON200 {
		d.StreamListItem(ctx, ProjectAPIKey{key, project.Id})
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOrganizationSlugForProjectIDFromAPIKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getOrganizationSlugForProjectID(ctx, d, h.Item.(ProjectAPIKey).ProjectId)
}

func generateAPIKeyAKA(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	obj, ok := d.Value.(ProjectAPIKey)
	if !ok {
		return nil, fmt.Errorf("could not cast value to transform")
	}
	return fmt.Sprintf("%s/apikey/%s", obj.ProjectId, obj.Id.MustGet()), nil
}

func transformNullable(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}

	v := reflect.ValueOf(d.Value)
	// Ensure we are actually dealing with a map
	if v.Kind() != reflect.Map {
		return d.Value, nil
	}

	// Nullable[T] is map[bool]T
	// check if the map has the key 'true' (which means a value was set).
	val := v.MapIndex(reflect.ValueOf(true))

	if val.IsValid() {
		return val.Interface(), nil
	}

	return nil, nil
}
