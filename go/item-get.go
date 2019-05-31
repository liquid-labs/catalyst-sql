package catsql

import (
  "context"

  "github.com/go-pg/pg/urlvalues"

  "github.com/Liquid-Labs/catalyst-core-model/go"
  "github.com/Liquid-Labs/go-rest/rest"
)

type PageRequest struct {
  Page         int
  ItemsPerPage int
}

// GetItem will attempt to retrieve an Entity by it's public ID, if available,
// and will otherwise fall back to the internal ID.
func GetItem(item *core.Entity, baseQuery *orm.Query, accessRoute AccessRoute, ctx context.Context) rest.RestError {
  if baseQuery == nil {
    return rest.BadRequestError(`Request does not resolve to a base query. Contact customer support if you believe this is a bug.`)
  }

  query := baseQuery.Context(ctx)
  if item.GetID().Valid {
    query = baseQuery.Where(`e.id=?`, item.GetID().Int64)
  } else if item.GetPubID().Valid {
    query = baseQuery.Where(`e.pub_id=?`, item.GetPubID().String)
  } else {
    return rest.BadRequestError(`Entity model does not provide a valid ID for retrieval.`, nil)
  }

  query = authorizedModel(query, accessRoute, core.StdAuthorizationGet, ctx)

  if err := query.Select(); err != nil {
    return rest.ServerError(fmt.Sprintf(`Problem retrieving entity '%s'.`, item.GetPubID().String), err)
  }
  return nil
}

func ListItems(baseQuery *orm.Query, accessRoute AccessRoute, pageRequest PageRequest, ctx context.Context) (int, rest.RestError) {
  if baseQuery == nil {
    return nil, rest.BadRequestError(`Request does not resolve to a base query. Contact customer support if you believe this is a bug.`)
  }

  pager := urlvalues.Pager{ Limit: pageRequest.ItemsPerPage, MaxLimit: 100 }
  pager.SetPage(pageRequest.Page)

  query := baseQuery.
    Context(ctx).
    Apply(pager.Pagination)

  query = authorizedModel(query, accessRoute, core.StdAuthorizationGet, ctx)

  if count, err := query.SelectAndCount(); err != nil {
    return nil, rest.ServerError(fmt.Sprintf(`Problem retrieving entity '%s'.`, e.GetPubID().String), err)
  } else {
    return count, nil
  }
}
