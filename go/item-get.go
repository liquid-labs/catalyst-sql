package catsql

import (
  "context"

  "github.com/go-pg/pg/urlvalues"

  "github.com/Liquid-Labs/catalyst-core-model/go/resources/authorizations"
  "github.com/Liquid-Labs/go-rest/rest"
)

type PageRequest struct {
  Page         int
  ItemsPerPage int
}

// GetItem will attempt to retrieve an Entity by it's public ID, if available,
// and will otherwise fall back to the internal ID.
func GetItem(id interface{}, baseQuery *orm.Query, accessRoute AccessRoute, ctx context.Context) rest.RestError {
  if baseQuery == nil {
    return rest.BadRequestError(`Request does not resolve to a base query. Contact customer support if you believe this is a bug.`)
  }

  query := baseQuery.Context(ctx)
  switch id.(type) {
  case string:
    query = baseQuery.Where(`e.pub_id=?`, id)
  case int64:
    query = baseQuery.Where(`e.id=?`, id)
  default:
    return rest.BadRequestError(fmt.Sprintf(`Invalid identifier type '%s' supplied to 'GetItem'.`, id.(type)), nil)
  }

  query = authorizedModel(query, accessRoute, authorizations.StdAuthorizationGet, ctx)

  if err := query.Select(); err != nil {
    // Notice we don't return the ID because it may be a oddly formatted
    // internal ID, which should not be revealed.
    // TODO: we should log the info though.
    return rest.ServerError(`Problem retrieving entity.`, err)
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

  query = authorizedModel(query, accessRoute, authorizations.StdAuthorizationGet, ctx)

  if count, err := query.SelectAndCount(); err != nil {
    return nil, rest.ServerError(fmt.Sprintf(`Problem retrieving entity '%s'.`, e.GetPubID().String), err)
  } else {
    return count, nil
  }
}
