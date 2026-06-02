package supabase

import (
	"context"

	"github.com/supabase/cli/pkg/api"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableSupabaseProject(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "supabase_project",
		Description: "Supabase Project",
		List: &plugin.ListConfig{
			Hydrate: listSupabaseProjects,
		},
		Columns: []*plugin.Column{
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The display name of the project."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "A unique identifier of the project."},
			{Name: "organization_slug", Type: proto.ColumnType_STRING, Description: "The organization slug."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the project was created."},
			{Name: "region", Type: proto.ColumnType_STRING, Description: "The project region."},
			{Name: "database", Type: proto.ColumnType_JSON, Description: "The database information."},
			{Name: "auth_service_config", Type: proto.ColumnType_JSON, Description: "The auth service config.", Hydrate: getProjectAuthServiceConfig, Transform: transform.FromValue()},
			{Name: "sso_providers", Type: proto.ColumnType_JSON, Description: "The SSO providers.", Hydrate: getProjectSSOProviders, Transform: transform.FromValue()},
			{Name: "tpa_integrations", Type: proto.ColumnType_JSON, Description: "The third-party auth integrations.", Hydrate: getProjectTPAIntegrations, Transform: transform.FromValue()},
			{Name: "postgREST", Type: proto.ColumnType_JSON, Description: "The openapi postgREST definition of the project's database.", Hydrate: getProjectDatabaseOpenAPISpec, Transform: transform.FromValue()},
			{Name: "pooler_config", Type: proto.ColumnType_JSON, Description: "The pooler config of the project's database.", Hydrate: getProjectDatabasePoolerConfig, Transform: transform.FromValue()},
			{Name: "pg_config", Type: proto.ColumnType_JSON, Description: "The postgres config of the project's database.", Hydrate: getProjectDatabasePGConfig, Transform: transform.FromValue()},
			{Name: "pgbouncer_config", Type: proto.ColumnType_JSON, Description: "The pgbouncer config of the project's database.", Hydrate: getProjectDatabasePGBouncerConfig, Transform: transform.FromValue()},
			{Name: "ssl_config", Type: proto.ColumnType_JSON, Description: "The SSL enforcement config of the project.", Hydrate: getProjectSSLConfig, Transform: transform.FromValue()},
			{Name: "vanity_domain_config", Type: proto.ColumnType_JSON, Description: "The current vanity domain config of the project.", Hydrate: getProjectVanityDomainSettings, Transform: transform.FromValue()},
			{Name: "network_restrictions", Type: proto.ColumnType_JSON, Description: "The network restrictions of the project.", Hydrate: getProjectNetworkRestrictions, Transform: transform.FromValue()},
			{Name: "custom_hostname", Type: proto.ColumnType_JSON, Description: "The custom hostname settings of the project.", Hydrate: getProjectCustomHostname, Transform: transform.FromValue()},
			{Name: "network_bans", Type: proto.ColumnType_JSON, Description: "The current network bans of the project.", Hydrate: getProjectNetworkBans, Transform: transform.FromValue()},
			{Name: "realtime_config", Type: proto.ColumnType_JSON, Description: "The realtime configuration of the project.", Hydrate: getProjectRealtimeConfig, Transform: transform.FromValue()},
			{Name: "postgres_service_config", Type: proto.ColumnType_JSON, Description: "The postgres service config of the project's database.", Hydrate: getProjectPostgresServiceConfig, Transform: transform.FromValue()},
			{Name: "storage_config", Type: proto.ColumnType_JSON, Description: "The storage configuration of the project.", Hydrate: getProjectStorageConfig, Transform: transform.FromValue()},

			// Common steampipe columns
			{
				Name:        "akas",
				Description: "Array of globally unique identifier strings (also known as) for the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Id").Transform(transform.EnsureStringArray),
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

func listSupabaseProjects(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.listSupabaseProjects", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1ListAllProjectsWithResponse(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.listSupabaseProjects", "query_error", err)
		return nil, err
	}

	if resp.JSON200 == nil {
		plugin.Logger(ctx).Warn("supabase_project.listSupabaseProjects", "status_code", resp.StatusCode(), "body", string(resp.Body))
		return nil, nil
	}

	for _, project := range *resp.JSON200 {
		d.StreamListItem(ctx, project)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getProjectAuthServiceConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectAuthServiceConfig", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetAuthServiceConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectAuthServiceConfig", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectSSOProviders(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectSSOProviders", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1ListAllSsoProviderWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectSSOProviders", "query_error", err)
		return nil, err
	}

	if resp.JSON200 == nil {
		plugin.Logger(ctx).Warn("supabase_project.getProjectSSOProviders", "status_code", resp.StatusCode(), "body", string(resp.Body))
		return nil, nil
	}

	return resp.JSON200.Items, nil
}

func getProjectTPAIntegrations(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectTPAIntegrations", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1ListProjectTpaIntegrationsWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectTPAIntegrations", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectDatabaseOpenAPISpec(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectDatabaseOpenApiSpec", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetDatabaseOpenapiWithResponse(ctx, projectID, &api.V1GetDatabaseOpenapiParams{})
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectDatabaseOpenApiSpec", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectDatabasePoolerConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectDatabasePoolerConfig", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetPoolerConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectDatabasePoolerConfig", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectDatabasePGConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectDatabasePGConfig", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetPostgresConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectDatabasePGConfig", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectDatabasePGBouncerConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectDatabasePGBouncerConfig", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetProjectPgbouncerConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectDatabasePGBouncerConfig", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectSSLConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectSSLConfig", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetSslEnforcementConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectSSLConfig", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectVanityDomainSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectVanityDomainSettings", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetVanitySubdomainConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectVanityDomainSettings", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectNetworkRestrictions(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectNetworkRestrictions", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetNetworkRestrictionsWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectNetworkRestrictions", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectCustomHostname(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectCustomHostname", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetHostnameConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectCustomHostname", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectNetworkBans(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectNetworkBans", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1ListAllNetworkBansEnrichedWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectNetworkBans", "query_error", err)
		return nil, err
	}

	return resp.JSON201, nil
}

func getProjectRealtimeConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectRealtimeConfig", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetRealtimeConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectRealtimeConfig", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectPostgresServiceConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectPostgresServiceConfig", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetPostgrestServiceConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectPostgresServiceConfig", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}

func getProjectStorageConfig(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	projectID := h.Item.(api.V1ProjectWithDatabaseResponse).Id
	client, err := getClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectStorageConfig", "connection_error", err)
		return nil, err
	}

	resp, err := client.V1GetStorageConfigWithResponse(ctx, projectID)
	if err != nil {
		plugin.Logger(ctx).Error("supabase_project.getProjectStorageConfig", "query_error", err)
		return nil, err
	}

	return resp.JSON200, nil
}
