package routers

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/gin-gonic/gin"
)

// constants TODO
const (
	RoleLeader   = "leader"
	RoleFollower = "follower"

	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

// Member TODO
type Member struct {
	*etcdserverpb.Member
	Role   string `json:"role"`
	Status string `json:"status"`
	DbSize int64  `json:"db_size"`
}

func getMembersHandler(c *gin.Context, client *clientv3.Client) (interface{}, error) {
	resp, err := client.MemberList(newEtcdCtx())
	if err != nil {
		return nil, err
	}

	members := []*Member{}
	for _, member := range resp.Members {
		if len(member.ClientURLs) > 0 {
			m := &Member{Member: member, Role: RoleFollower, Status: StatusUnhealthy}
			resp, err := client.Status(newEtcdCtx(), m.ClientURLs[0])
			if err == nil {
				m.Status = StatusHealthy
				m.DbSize = resp.DbSize
				if resp.Leader == resp.Header.MemberId {
					m.Role = RoleLeader
				}
			}
			members = append(members, m)
		}
	}

	return members, nil
}
