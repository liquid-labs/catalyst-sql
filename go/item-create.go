package catsql

import (
  "github.com/go-pg/pg"

  "github.com/Liquid-Labs/catalyst-core-model/go/resources/entities"
  "github.com/Liquid-Labs/go-rest/rest"
)

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
