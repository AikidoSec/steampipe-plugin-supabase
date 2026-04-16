package supabase

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const pluginName = "steampipe-plugin-supabase"

// Plugin creates this (supabase) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             pluginName,
		DefaultTransform: transform.FromCamel().Transform(transform.NullIfZeroValue),
		DefaultGetConfig: &plugin.GetConfig{},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		TableMap: map[string]*plugin.Table{
			"supabase_function":         tableSupabaseFunction(ctx),
			"supabase_organization":     tableSupabaseOrganization(ctx),
			"supabase_project":          tableSupabaseProject(ctx),
			"supabase_secret":           tableSupabaseSecret(ctx),
			"supabase_api_key":          tableSupabaseAPIKey(ctx),
			"supabase_bucket":           tableSupabaseBucket(ctx),
			"supabase_security_advisor": tableSupabaseSecurityAdvisor(ctx),
		},
	}

	return p
}
