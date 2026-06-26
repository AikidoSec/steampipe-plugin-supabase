package supabase

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/rate_limiter"
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
		RateLimiters: []*rate_limiter.Definition{
			{
				Name:       "supabase_global",
				FillRate:   25,
				BucketSize: 25,
				Scope:      []string{rate_limiter.RateLimiterScopeConnection},
			},
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
