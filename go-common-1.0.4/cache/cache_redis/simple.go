package cache_redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type Option struct {
	Addrs      []string
	MasterName string
	Username   string
	Password   string
	DB         int
}

func SimpleCluster(opts Option) redis.UniversalClient {
	options := &redis.UniversalOptions{
		Addrs:      opts.Addrs,
		MasterName: opts.MasterName,
		Username:   opts.Username,
		Password:   opts.Password,
		DB:         opts.DB,
	}

	if opts.MasterName != "" {
		return redis.NewFailoverClient(options.Failover())
	} else if len(opts.Addrs) > 1 {
		return redis.NewClusterClient(options.Cluster())
	}
	simpleClient := redis.NewClient(options.Simple())
	ctx, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()
	info := simpleClient.Info(ctx, "cluster")
	if info.Err() != nil {
		return simpleClient
	}
	if !strings.Contains(info.String(), "cluster_enabled:1") {
		return simpleClient
	}

	nodes := simpleClient.ClusterNodes(context.Background())
	if nodes.Err() != nil {
		return simpleClient
	}
	_ = simpleClient.Close()
	nodesContent := nodes.String()
	nodesContent = strings.TrimPrefix(nodesContent, "cluster nodes: ")
	nodesContent = strings.TrimSpace(nodesContent)
	lines := strings.SplitN(nodesContent, "\n", -1)
	nodeAddrs := make([]string, 0, len(lines))
	for _, line := range lines {
		nodeAddrs = append(nodeAddrs, readAddr(line))
	}
	options.Addrs = nodeAddrs
	return redis.NewClusterClient(options.Cluster())

}
func readAddr(line string) string {
	fields := strings.Fields(line)
	addr := fields[1]

	index := strings.Index(addr, "@")
	if index > 0 {
		addr = addr[:index]
	}
	return addr

}
