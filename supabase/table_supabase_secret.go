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

func tableSupabaseSecret(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "supabase_secret",
		Description: "Supabase Secret",
		List: &plugin.ListConfig{
			ParentHydrate: listSupabaseProjects,
			Hydrate:       listSupabaseSecrets,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the secret."},
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The secret value."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the function was last modified."},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The ID of the project."},

			{
				Name:        "organization_slug",
				Type:        proto.ColumnType_STRING,
				Description: "The organization slug.",
				Hydrate:     getOrganizationSlugForProjectIDFromSecret,
				Transform:   transform.FromValue(),
			},

			// Common steampipe columns
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue().Transform(generateSecretAKA).Transform(transform.EnsureStringArray),
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

type Secret struct {
	api.SecretResponse
	ProjectId string
}

//// LIST FUNCTION

func listSupabaseSecrets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the project data
	project := h.Item.(api.V1ProjectWithDatabaseResponse)

	// Create client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_secret.listSupabaseSecrets", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1ListAllSecretsWithResponse(ctx, project.Id)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_secret.listSupabaseSecrets", "query_error", err)
		return nil, err
	}

	for _, secret := range *resp.JSON200 {
		d.StreamListItem(ctx, Secret{secret, project.Id})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOrganizationSlugForProjectIDFromSecret(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getOrganizationSlugForProjectID(ctx, d, h.Item.(Secret).ProjectId)
}

func generateSecretAKA(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	secret, ok := d.Value.(Secret)
	if !ok {
		return nil, fmt.Errorf("could not cast value to transform")
	}
	return fmt.Sprintf("%s/secret/%s", secret.ProjectId, secret.Name), nil
}
