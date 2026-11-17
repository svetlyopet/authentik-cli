package authentik

// JSON fixtures shaped like authentik API responses.
//
// The generated client rejects responses that miss a required property, and
// its models skip read-only properties when marshalled, so fixtures are built
// as plain maps carrying every property the client requires.

type obj = map[string]any

// merge returns base with the entries of extra added or overridden.
func merge(base, extra obj) obj {
	for k, v := range extra {
		base[k] = v
	}
	return base
}

func list(items ...obj) []obj {
	if items == nil {
		return []obj{}
	}
	return items
}

// paginated wraps results in a paginated list response. next is the next
// page number, or 0 on the last page.
func paginated(next int, results ...obj) obj {
	return obj{
		"pagination": obj{
			"next":        next,
			"previous":    0,
			"count":       len(results),
			"current":     1,
			"total_pages": 1,
			"start_index": 1,
			"end_index":   len(results),
		},
		"results":      list(results...),
		"autocomplete": obj{},
	}
}

func roleFixture(pk, name string) obj {
	return obj{"pk": pk, "name": name}
}

func groupFixture(pk, name string, extra obj) obj {
	return merge(obj{
		"pk":                  pk,
		"num_pk":              1,
		"name":                name,
		"parents_obj":         list(),
		"users_obj":           list(),
		"roles_obj":           list(),
		"inherited_roles_obj": list(),
		"children":            list(),
		"children_obj":        list(),
	}, extra)
}

func partialGroupFixture(pk, name string) obj {
	return obj{"pk": pk, "num_pk": 1, "name": name}
}

func userFixture(pk int, username, name string, extra obj) obj {
	return merge(obj{
		"pk":                   pk,
		"username":             username,
		"name":                 name,
		"date_joined":          "2026-01-01T00:00:00Z",
		"is_superuser":         false,
		"groups_obj":           list(),
		"roles_obj":            list(),
		"avatar":               "",
		"uid":                  "uid",
		"uuid":                 "00000000-0000-0000-0000-000000000000",
		"password_change_date": "2026-01-01T00:00:00Z",
		"last_updated":         "2026-01-01T00:00:00Z",
	}, extra)
}

func providerFields(pk int, name, metaModelName string) obj {
	return obj{
		"pk":                                    pk,
		"name":                                  name,
		"component":                             "",
		"assigned_application_slug":             "",
		"assigned_application_name":             "",
		"assigned_backchannel_application_slug": "",
		"assigned_backchannel_application_name": "",
		"verbose_name":                          "",
		"verbose_name_plural":                   "",
		"meta_model_name":                       metaModelName,
	}
}

func providerFixture(pk int, name, metaModelName string) obj {
	return merge(providerFields(pk, name, metaModelName), obj{
		"authorization_flow": "authz-flow",
		"invalidation_flow":  "invalidation-flow",
	})
}

func oauth2ProviderFixture(pk int, name string, extra obj) obj {
	return merge(merge(providerFields(pk, name, ProviderTypeMetaModelOIDC), obj{
		"authorization_flow": "",
		"invalidation_flow":  "",
		"redirect_uris":      list(),
	}), extra)
}

// applicationFixture builds an application; provider may be nil for an
// application without a provider.
func applicationFixture(pk, name, slug string, provider obj) obj {
	var providerObj any
	if provider != nil {
		providerObj = provider
	}
	return obj{
		"pk":                        pk,
		"pbm_uuid":                  pk,
		"name":                      name,
		"slug":                      slug,
		"provider_obj":              providerObj,
		"backchannel_providers_obj": list(),
		"launch_url":                nil,
		"meta_icon_url":             nil,
		"meta_icon_themed_urls":     nil,
	}
}

func flowFixture(pk, name, slug, designation string) obj {
	return obj{
		"pk":                        pk,
		"policybindingmodel_ptr_id": pk,
		"name":                      name,
		"slug":                      slug,
		"title":                     name,
		"designation":               designation,
		"background_url":            "",
		"background_themed_urls":    nil,
		"stages":                    list(),
		"policies":                  list(),
		"cache_count":               0,
		"export_url":                "",
	}
}
