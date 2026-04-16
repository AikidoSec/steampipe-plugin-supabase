package supabase

import (
	"context"
	"fmt"

	"github.com/supabase/cli/pkg/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableSupabaseSecurityAdvisor(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "supabase_security_advisor",
		Description: "Supabase Security Advisor",
		List: &plugin.ListConfig{
			ParentHydrate: listSupabaseProjects,
			Hydrate:       listSupabaseSecurityAdvisors,
		},
		Columns: []*plugin.Column{
			{Name: "cache_key", Type: proto.ColumnType_STRING, Description: "The unique cache key of the advisor."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The advisor name."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "The advisor description."},
			{Name: "detail", Type: proto.ColumnType_STRING, Description: "The advisor detail text."},
			{Name: "remediation", Type: proto.ColumnType_STRING, Description: "The advisor remediation guidance."},
			{Name: "level", Type: proto.ColumnType_STRING, Description: "The severity level of the advisor."},
			{Name: "facing", Type: proto.ColumnType_STRING, Description: "Whether the advisor affects external or internal surfaces."},
			{Name: "categories", Type: proto.ColumnType_JSON, Description: "The advisor categories."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Additional advisor metadata."},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The ID of the project."},
			{Name: "raw", Type: proto.ColumnType_JSON, Description: "The raw advisor payload.", Transform: transform.FromField("Advisor")},

			{
				Name:        "organization_slug",
				Type:        proto.ColumnType_STRING,
				Description: "The organization slug.",
				Hydrate:     getOrganizationSlugForProjectIDFromSecurityAdvisor,
				Transform:   transform.FromValue(),
			},

			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromValue().Transform(generateSecurityAdvisorAKA).Transform(transform.EnsureStringArray),
			},
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Title"),
			},
		},
	}
}

type SecurityAdvisorMetadata struct {
	Entity      *string                                         `json:"entity,omitempty"`
	FkeyColumns *[]float32                                      `json:"fkey_columns,omitempty"`
	FkeyName    *string                                         `json:"fkey_name,omitempty"`
	Name        *string                                         `json:"name,omitempty"`
	Schema      *string                                         `json:"schema,omitempty"`
	Type        *api.V1ProjectAdvisorsResponseLintsMetadataType `json:"type,omitempty"`
}

type SecurityAdvisor struct {
	CacheKey    string                                         `json:"cache_key"`
	Categories  []api.V1ProjectAdvisorsResponseLintsCategories `json:"categories"`
	Description string                                         `json:"description"`
	Detail      string                                         `json:"detail"`
	Facing      api.V1ProjectAdvisorsResponseLintsFacing       `json:"facing"`
	Level       api.V1ProjectAdvisorsResponseLintsLevel        `json:"level"`
	Metadata    *SecurityAdvisorMetadata                       `json:"metadata,omitempty"`
	Name        api.V1ProjectAdvisorsResponseLintsName         `json:"name"`
	Remediation string                                         `json:"remediation"`
	Title       string                                         `json:"title"`
}

type ProjectSecurityAdvisor struct {
	SecurityAdvisor
	Advisor   SecurityAdvisor
	ProjectId string
}

func listSupabaseSecurityAdvisors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	project := h.Item.(api.V1ProjectWithDatabaseResponse)

	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_security_advisor.listSupabaseSecurityAdvisors", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetSecurityAdvisorsWithResponse(ctx, project.Id, nil)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_security_advisor.listSupabaseSecurityAdvisors", "query_error", err)
		return nil, err
	}

	if resp.JSON200 == nil {
		return nil, nil
	}

	for _, advisor := range resp.JSON200.Lints {
		item := SecurityAdvisor{
			CacheKey:    advisor.CacheKey,
			Categories:  advisor.Categories,
			Description: advisor.Description,
			Detail:      advisor.Detail,
			Facing:      advisor.Facing,
			Level:       advisor.Level,
			Name:        advisor.Name,
			Remediation: advisor.Remediation,
			Title:       advisor.Title,
		}

		if advisor.Metadata != nil {
			item.Metadata = &SecurityAdvisorMetadata{
				Entity:      advisor.Metadata.Entity,
				FkeyColumns: advisor.Metadata.FkeyColumns,
				FkeyName:    advisor.Metadata.FkeyName,
				Name:        advisor.Metadata.Name,
				Schema:      advisor.Metadata.Schema,
				Type:        advisor.Metadata.Type,
			}
		}

		d.StreamListItem(ctx, ProjectSecurityAdvisor{
			SecurityAdvisor: item,
			Advisor:         item,
			ProjectId:       project.Id,
		})

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getOrganizationSlugForProjectIDFromSecurityAdvisor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getOrganizationSlugForProjectID(ctx, d, h.Item.(ProjectSecurityAdvisor).ProjectId)
}

func generateSecurityAdvisorAKA(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	advisor, ok := d.Value.(ProjectSecurityAdvisor)
	if !ok {
		return nil, fmt.Errorf("could not cast value to transform")
	}

	return fmt.Sprintf("%s/security-advisor/%s", advisor.ProjectId, advisor.CacheKey), nil
}
