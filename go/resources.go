package catsql

import (
  "context"

  "github.com/go-pg/pg"
  "github.com/go-pg/pg/urlvalues"

  "github.com/Liquid-Labs/catalyst-core-model/go/resources/entities"
  "github.com/Liquid-Labs/catalyst-core-model/go/resources/authorizations"
  "github.com/Liquid-Labs/go-rest/rest"
)

type PageRequest struct {
  Page         int
  ItemsPerPage int
}

// CreateItem will check user permissions via the indicated accessRoute and
// create a new record of the provided Entity (sub-type) where authorized.
func CreateItem(item *entities.Entity, accessRoute AccessRoute, ctx context.Context) rest.RestError {
  if item == nil {
    return rest.BadRequestError(`Entity for creation cannot be nil.`, nil)
  }

  if authResponse, restErr := CheckResourceAuthorization(e, `create`); restErr != nil {
    return restErr
  } else if !authResponse.Granted {
    // TODO: get helper to get us the name... method reciever for Entity?
    return rest.AuthorizationError(`User not authorized to create resource.`)
  } else {
    if err := db.Model(e).Create(); err != nil {
      return rest.ServerError(`Problem creating resource.`, err)
    } else {
      return nil
    }
  }
}

// GetItem will attempt to retrieve an Entity by either the public or internal
// ID. Which to use is determined by the 'id' type, which must be either a
// string (for public ID) or int64 (for internal ID). The base query is
// typically just 'db.Model(yourStruct)', where the struct used must embed
// Entity. GetItem adds hte necessary authorization checks to the provided
// base query.
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

// ListItems retrieves the set of items selected by the base query to which the
// user has the necessary access rights according to the access route selected.
// The base query may be as simble as:
//
// list = &make([]FinalClass)
// query = db.Model(list)
//
// or it may include additional filter clauses.
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
