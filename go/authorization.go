package catsql

import (
  "context"
  "time"

  "firebase.google.com/go/auth"

  "github.com/Liquid-Labs/catalyst-core-model/go"
  "github.com/go-pg/pg"
  "github.com/go-pg/pg/orm"
)

type AuthorizationOptions struct {
  PublicOnly    bool
  OwnershipOnly bool
  MultiTarget   bool
  LimitQuery    *orm.Query
}

- take in first limiting query, then join that with azn query
- support user being authenticated but looking at public stuff
- public stuff is either public or private.
- All access is either through a public grant, ownership, or a direct grant. No groups. (Though we will still implement a 'team' function, but the teams are really implicit reductions of members with the same grant structure)
- Flips logic here and in the 'Grant' function.
- Shared systems disable group-grants and negative grants.
- Private systems enable group-grants and negative grants.
- Can we do public and self-check first and preempt the grant check?

func AuthorizedSingleModel(e *core.Entity, ctx context.Context) *orm.Query {
  auth.Token authToken = ctx.Value(core.AuthTokenKey).(*auth.Token)
  // TODO: is it necessary to check the time? And if so, where? I'm thinking in
  // the base handler where the token is bound to the context. Our rule for
  // short lived requests thus being: it's authorized if action starts
  // authorized.
  if authToken == nil { // limited to public actions
    return db.Model(e).
      // Check for a public grant for non-authorized user.
      Join(`azn_grants AS azn_grount ON azn_grant.target=e.id`).
      Where(`azn_grant.subject IS NULL AND azn_grant.DENY=FALSE AND azn_grant.action=?`, model.StdActions.EntityGet)
  } else if auth.HasClaim(`root`, authToken.Claims)  {
    // No auth check if user claims root.
    return db.Model(e)
  } else {
    // 1) Select recursive group membership, with depth.
    // 2) Select auth grant subjects against public, self, and groups with that
    //    precedence order.
    // 3) select grant assertion (which may be a positive or negative grant)
    //   from highest precedence.
    // TODO: make this more efficient with cached groups;
    // a) select the group memberships when the user logs in
    // b) add groups and 'last updated' stamp info as a signed cookie
    // c) cookie is passed pack to server where it is checked and either used or
    //    refreshed.
    authID := authToken.UID

    return db.Model(e).
      With(`azn_entities`, db.Model(e).
        With(`filtered_groups`,
          db.Model(e).
          ColumnExp(`COALESCE(group.depth, -1) AS depth`).
          With(`WITH RECURSIVE groups(id, depth)`, // WithRecursive(`groups(group, depth)`,
            `SELECT group AS id, 0 FROM azn_group_members agm1 JOIN users u ON agm1.member=u.id WHERE u.auth_id=?
            UNION
            SELECT group AS id, depth+1 FROM azn_group_members agm2 WHERE agm2.member=agm1.group`,
            authID).
          Join(`azn_grants AS azn_grant ON azn_grant.target=?`, e.GetID()).
          Join(`LEFT OUTER JOIN groups AS group ON azn_grant.subject=group.id OR azn_grount.subject IS NULL`).
          Where(`azn_grants.action=?`, model.StdActions.EntityGet).
          GroupBy(`e.id, depth`).
          Having("MIN(depth)")).
        Table(`filtered_groups`).
        Where(`deny=FALSE`)).
      Table(`azn_entities`)
  }

}

func AuthorizedModel(e *core.Entity, ctx context.Context) *orm.Query {
  auth.Token authToken = ctx.Value(core.AuthTokenKey).(*auth.Token)
  // TODO: is it necessary to check the time? And if so, where? I'm thinking in
  // the base handler where the token is bound to the context. Our rule for
  // short lived requests thus being: it's authorized if action starts
  // authorized.
  if authToken == nil { // limited to public actions
    return db.Model(e).
      // Check for a public grant for non-authorized user.
      Join(`azn_grants AS azn_grount ON azn_grant.target=e.id`).
      Where(`azn_grant.subject IS NULL AND azn_grant.DENY=FALSE AND azn_grant.action=?`, model.StdActions.EntityGet)
  } else if auth.HasClaim(`root`, authToken.Claims)  {
    // No auth check if user claims root.
    return db.Model(e)
  } else {
    // 1) Select recursive group membership, with depth.
    // 2) Select auth grant subjects against public, self, and groups with that
    //    precedence order.
    // 3) select grant assertion (which may be a positive or negative grant)
    //   from highest precedence.
    // TODO: make this more efficient with cached groups;
    // a) select the group memberships when the user logs in
    // b) add groups and 'last updated' stamp info as a signed cookie
    // c) cookie is passed pack to server where it is checked and either used or
    //    refreshed.
    authID := authToken.UID

    return db.Model(e).
      With(`azn_entities`, db.Model(e).
        With(`filtered_groups`,
          db.Model(e).
          ColumnExp(`COALESCE(group.depth, -1) AS depth`).
          With(`WITH RECURSIVE groups(id, depth)`, // WithRecursive(`groups(group, depth)`,
            `SELECT group AS id, 0 FROM azn_group_members agm1 JOIN users u ON agm1.member=u.id WHERE u.auth_id=?
            UNION
            SELECT group AS id, depth+1 FROM azn_group_members agm2 WHERE agm2.member=agm1.group`,
            authID).
          Join(`azn_grants AS azn_grant ON azn_grant.target=e.id`).
          Join(`LEFT OUTER JOIN groups AS group ON azn_grant.subject=group.id OR azn_grount.subject IS NULL`).
          Where(`azn_grants.action=?`, model.StdActions.EntityGet).
          GroupBy(`e.id, depth`).
          Having("MIN(depth)")).
        Table(`filtered_groups`).
        Where(`deny=FALSE`)).
      Table(`azn_entities`)
  }

}
