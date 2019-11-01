package routers

import (
	"log"

	"github.com/VendettA01/e3w/auth"
	"github.com/VendettA01/e3w/conf"
	"github.com/VendettA01/e3w/e3ch"
	"github.com/VendettA01/e3w/resp"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"github.com/soyking/e3ch"
)

const (
	ETCD_USERNAME_HEADER = "X-Etcd-Username"
	ETCD_PASSWORD_HEADER = "X-Etcd-Password"
)

type e3chHandler func(*gin.Context, *client.EtcdHRCHYClient) (interface{}, error)

type groupHandler func(e3chHandler) resp.RespHandler

func withE3chGroup(e3chClt *client.EtcdHRCHYClient, config *conf.Config) groupHandler {
	return func(h e3chHandler) resp.RespHandler {
		return func(c *gin.Context) (interface{}, error) {
			clt := e3chClt
			if config.EtcdAuth {
				var err error
				username := c.Request.Header.Get(ETCD_USERNAME_HEADER)
				password := c.Request.Header.Get(ETCD_PASSWORD_HEADER)
				clt, err = e3ch.CloneE3chClient(username, password, e3chClt)
				if err != nil {
					return nil, err
				}
				defer clt.EtcdClient().Close()
			}
			return h(c, clt)
		}
	}
}

type etcdHandler func(*gin.Context, *clientv3.Client) (interface{}, error)

func etcdWrapper(h etcdHandler) e3chHandler {
	return func(c *gin.Context, e3chClt *client.EtcdHRCHYClient) (interface{}, error) {
		return h(c, e3chClt.EtcdClient())
	}
}

func InitRouters(g *gin.Engine, config *conf.Config, e3chClt *client.EtcdHRCHYClient) {
	if err := auth.InitAuthFromConf(); err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	g.Static("/public", "./static/dist")
	g.GET("/", func(c *gin.Context) {
		c.File("./static/dist/index.html")
	})

	private := g.Group("/")
	private.Use(auth.AuthRequired)

	// login route cannot be protected by withAuth
	g.POST("/login", auth.LogIn)

	// checkToken and logout are protected by withAuth, so
	// only authenticated users can use these
	private.GET("/checkToken", auth.CheckToken)
	private.GET("/logout", auth.LogOut)

	e3chGroup := withE3chGroup(e3chClt, config)

	// key/value actions
	private.GET("/kv/*key", resp.Resp(e3chGroup(getKeyHandler)))
	private.POST("/kv/*key", resp.Resp(e3chGroup(postKeyHandler)))
	private.PUT("/kv/*key", resp.Resp(e3chGroup(putKeyHandler)))
	private.DELETE("/kv/*key", resp.Resp(e3chGroup(delKeyHandler)))

	// members actions
	private.GET("/members", resp.Resp(e3chGroup(etcdWrapper(getMembersHandler))))

	// roles actions
	private.GET("/roles", resp.Resp(e3chGroup(etcdWrapper(getRolesHandler))))
	private.POST("/role", resp.Resp(e3chGroup(etcdWrapper(createRoleHandler))))
	private.GET("/role/:name", resp.Resp(e3chGroup(getRolePermsHandler)))
	private.DELETE("/role/:name", resp.Resp(e3chGroup(etcdWrapper(deleteRoleHandler))))
	private.POST("/role/:name/permission", resp.Resp(e3chGroup(createRolePermHandler)))
	private.DELETE("/role/:name/permission", resp.Resp(e3chGroup(deleteRolePermHandler)))

	// users actions
	private.GET("/users", resp.Resp(e3chGroup(etcdWrapper(getUsersHandler))))
	private.POST("/user", resp.Resp(e3chGroup(etcdWrapper(createUserHandler))))
	private.GET("/user/:name", resp.Resp(e3chGroup(etcdWrapper(getUserRolesHandler))))
	private.DELETE("/user/:name", resp.Resp(e3chGroup(etcdWrapper(deleteUserHandler))))
	private.PUT("/user/:name/password", resp.Resp(e3chGroup(etcdWrapper(setUserPasswordHandler))))
	private.PUT("/user/:name/role/:role", resp.Resp(e3chGroup(etcdWrapper(grantUserRoleHandler))))
	private.DELETE("/user/:name/role/:role", resp.Resp(e3chGroup(etcdWrapper(revokeUserRoleHandler))))
}
