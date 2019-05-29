package catsql

import (
  "context"

  "github.com/Liquid-Labs/catalyst-core-model/go"
  "github.com/Liquid-Labs/go-rest/rest"
)

// GetItem will attempt to retrieve an Entity by it's public ID, if available,
// and will otherwise fall back to the internal ID.
func GetItem(e *core.Entity, ctx context.Context) (interface{}, rest.RestError) {
  if (e.GetPubID().IsValid()) {
    query := AuthorizedModel(e, ctx).
      Where("pub_id=?", e.GetPubID().String)
    if err := query.Select(queryParams...); err != nil {
      return nil, rest.ServerError(fmt.Sprintf(`Problem selecting entity based on public ID '%s'.`, e.GetPubID().String), err)
    }
    
    return e, nil
  } else if (e.GetID().IsValid()) {
    query := AuthorizedModel(e, ctx).
      Where("id=?" , e.ID().Int64)
    if err := query.Select(queryParams...); err != nil {
      return nil, rest.ServerError(fmt.Sprintf(`Could not select entity based on model: %+v`, e), err)
    }

    return e, nil
  } else {
    return nil, rest.BadRequestError(fmt.Sprintf(`Could not extract any ID from: %+v`, e), nil)
  }
}
