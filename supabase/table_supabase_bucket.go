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

func tableSupabaseBucket(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "supabase_bucket",
		Description: "Supabase Bucket",
		List: &plugin.ListConfig{
			ParentHydrate: listSupabaseProjects,
			Hydrate:       listSupabaseBuckets,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The id of the bucket."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the bucket."},
			{Name: "owner", Type: proto.ColumnType_STRING, Description: "The owner of the bucket."},
			{Name: "public", Type: proto.ColumnType_BOOL, Description: "Whether or not the bucket is publicly accessible."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the bucket was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the bucket was last modified."},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The ID of the project."},

			{
				Name:        "organization_id",
				Type:        proto.ColumnType_STRING,
				Description: "The organization ID.",
				Hydrate:     getOrganizationIDForProjectIDFromSecret,
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

type Bucket struct {
	api.V1StorageBucketResponse
	ProjectId string
}

//// LIST FUNCTION

func listSupabaseBuckets(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the project data
	project := h.Item.(api.V1ProjectWithDatabaseResponse)

	// Create client
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_bucket.listSupabaseBuckets", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1ListAllBucketsWithResponse(ctx, project.Id)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_bucket.listSupabaseBuckets", "query_error", err)
		return nil, err
	}

	for _, secret := range *resp.JSON200 {
		d.StreamListItem(ctx, Bucket{secret, project.Id})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOrganizationIDForProjectIDFromBucket(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getOrganizationIDForProjectID(ctx, d, h.Item.(Bucket).ProjectId)
}

func generateBucketAKA(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	secret, ok := d.Value.(Secret)
	if !ok {
		return nil, fmt.Errorf("could not cast value to transform")
	}
	return fmt.Sprintf("%s/secret/%s", secret.ProjectId, secret.Name), nil
}
