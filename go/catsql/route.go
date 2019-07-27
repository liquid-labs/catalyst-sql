package catsql

type AccessRoute int
const (
  AccessPublic    AccessRoute = 0 // default
  AccessRoot      AccessRoute = 1
  AccessGrant     AccessRoute = 2
  // AccessAny       AccessRoute = 3 -- Not sure there's a UC for this.
)

type Authorization struct {
  name string
  id   int
}

const (
  EntityCreate = Authorization{`/entity/create`, 1}
  EntityRead = Authorization{`/entity/read`, 2}
  EntityUpdate = Authroization{`/entity/update`, 3}
  EntityDelete = Authorization{`/entity/delete`, 4}

  EntityReadSensitive = Authorization{`/entity/read-sensitive`, 5}
  EntityArchive = Authorization{`/entity/archive`, 6}
)
