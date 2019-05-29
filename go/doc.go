// Package 'catsql' provides useful DB functions that make it easy to translate
// between REST-ful request and DB actions. catsql covers:
// * database initialization; i.e., creating the schema in an empty database
// * warming up common prepared statements
// * authorization cognizant:
//   * resource item creation
//   * browsing resource items.
//   * browsing resource items by context
//   * basic resource search
//   * item detail retrieval
//   * item update
//   * item archival (soft deletion)
//
// We "warm up" the database as a simple optimization for long-running services
// under load. This does slightly deoptimizes startup times, which is
// undesirable in an environment that's frequently scaling to zero. Parallel
// resource micro-service startup mitigates this, and we do not consider it
// an issue.
package catsql
