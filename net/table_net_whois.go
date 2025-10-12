package net

import (
	"context"

	"github.com/likexian/whois"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableNetWhois(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "net_whois",
		Description: "WHOIS information for a domain.",
		List: &plugin.ListConfig{
			Hydrate:    tableNetWhoisList,
			KeyColumns: plugin.SingleColumn("domain"),
		},
		Columns: []*plugin.Column{
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "Domain name to query WHOIS information for."},
			{Name: "raw", Type: proto.ColumnType_STRING, Description: "Raw WHOIS response text."},
		},
	}
}

type tableNetWhoisRow struct {
	Domain string `json:"domain"`
	Raw    string `json:"raw"`
}

func tableNetWhoisList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	domain := d.EqualsQualString("domain")

	if domain == "" {
		return nil, nil
	}

	// Perform WHOIS lookup
	result, err := whois.Whois(domain)
	if err != nil {
		plugin.Logger(ctx).Error("net_whois.tableNetWhoisList", "WHOIS lookup failed", err)
		return nil, err
	}

	row := tableNetWhoisRow{
		Domain: domain,
		Raw:    result,
	}

	d.StreamListItem(ctx, row)

	return nil, nil
}
