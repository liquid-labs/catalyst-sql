package catsql

import (
  "context"

  "firebase.google.com/go/auth"

  "github.com/go-pg/pg/orm"
)

// type Identifiable interface {
//  GetID() int64
//  GetPubID() string
//}

type AuthorizationResponse struct {
  Granted bool
  Cookie  interface{} // could be any JSON derived structure; string, int, float, map, or array.
}

//func CheckResourceAuthorization(i Identifiable, action string) (*AuthorizationResponse, rest.RestError) {
//  return nil, nil
//}

func hasClaim(token *auth.Token, testClaim string) {
  for _, claim := range token.Claims {
    if claim == testClaim {
      return true
    }
  }
  return false
}

func resolveAuthorization(authorization interface{}, query *orm.Query) *orm.Query {
  switch t := authorization.(type) {
  case int:
    return query.Where(`azn_grant.authorization=?`, authorization)
  case int64:
    return query.Where(`azn_grant.authorization=?`, authorization)
  case string:
    return query.
      Join(`JOIN azn_authorizations AS azn_authorization ON azn_grant.authorization=azn_authorization.id`).
      Where(`azn_authorization.name=?`, authorization)
  }
}

func authorizedModel(baseQuery *orm.Query, accessRoute /*azn.*/AccessRoute, authorization interface{}, ctx context.Context) *orm.Query {
  if accessRoute == AccessPublic {
    return authorizedPublicModel(baseQuery, authorization)
  } else {
    authToken := ctx.Value(core.AuthTokenKey).(*auth.Token)
    if authToken == nil {
      return authorizedPublicModel(baseQuery, authorization)
    }
    // else, we have an auth token
    if accessRoute == AccessRoot {
      if hasClaim(authToken, `root`) {
        return baseQuery
      } else {
        return rest.BadRequestError(`Cannot make 'root' request as non-root user.`, nil)
      }
    } else if accessRoute == AccessGrant {
      return authorizedGrantModel(e, baseQuery, authorization)
    } else {
      log.Panicf(`Unmatched 'access route' value: '%d'`, accessRoute)
    }
  }
}

func authorizedPublicModel(q *orm.Query, authorization interface{}) *orm.Query {
  q = q.
    Join("JOIN azn_grants AS azn_grant ON azn_grant.target=e.id").
    Where("read_public=TRUE OR azn_grant.subject IS NULL")
  return resolveAuthorization(authorization, q)
}

func authorizedGrantModel(q *orm.Query, authorization interface{}) *orm.Query {
  authID := authToken.UID

  query := baseQuery.
    WrapWith(`WITH RECURSIVE group(id) AS (
        SELECT agm.group AS id FROM azn_group_members agm JOIN users u ON agm.member=u.id WHERE u.auth_id=?
      UNION
        SELECT agm.group AS id FROM azn_group_members agm WHERE agm.member=group.id`,
      authID).
    Join(`JOIN group`).
    Join(`JOIN container`).
    Join(`JOIN users u ON u.auth_id=?`, authID).
    Join(`JOIN azn_grants AS azn_grant ON (azn_grant.subject IS NULL OR azn_grant.subject=u.id OR azn_grant.subject=group.id) AND (azn_grant.target=e.id OR e.containers @> ARRAY[azn_grant.target])`).
    Where(`containers @> ARRAY['s']::varchar[]`).
    Group(`entity.id`)

  return resolveAuthorization(authorization, query)
}
