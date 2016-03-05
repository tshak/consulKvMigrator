package consulKvMigrator

import consul "github.com/hashicorp/consul/api"

type consulKvClient interface {
	Get(key string, q *consul.QueryOptions) (*consul.KVPair, *consul.QueryMeta, error)
	Put(p *consul.KVPair, q *consul.WriteOptions) (*consul.WriteMeta, error)
	CAS(p *consul.KVPair, q *consul.WriteOptions) (bool, *consul.WriteMeta, error)
}
