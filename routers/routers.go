package routers

import (
	"github.com/VendettA01/e3w/conf"
	"github.com/VendettA01/e3w/e3ch"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	"github.com/soyking/e3ch"
)

const (
	ETCD_USERNAME_HEADER = "X-Etcd-Username"
	ETCD_PASSWORD_HEADER = "X-Etcd-Password"
)

type e3chHandler func(*gin.Context, *client.EtcdHRCHYClient) (interface{}, error)

type groupHandler func(e3chHandler) respHandler

func withE3chGroup(e3chClt *client.EtcdHRCHYClient, config *conf.Config) groupHandler {
	return func(h e3chHandler) respHandler {
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

	g.Static("/public", "./static/dist")

	private := g.Group("/")
	if conf.Conf.Auth {
		g.GET("/login", func(c *gin.Context) {
			c.File("./static/dist/login.html")
		})
		g.POST("/login", logIn)
		g.GET("/logout", logOut)
		private.Use(authRequired)
	}

	private.GET("/", func(c *gin.Context) {
		c.File("./static/dist/index.html")
	})

	e3chGroup := withE3chGroup(e3chClt, config)

	// key/value actions
	private.GET("/kv/*key", resp(e3chGroup(getKeyHandler)))
	private.POST("/kv/*key", (resp(e3chGroup(postKeyHandler))))
	private.PUT("/kv/*key", resp(e3chGroup(putKeyHandler)))
	private.DELETE("/kv/*key", resp(e3chGroup(delKeyHandler)))

	// members actions
	private.GET("/members", resp(e3chGroup(etcdWrapper(getMembersHandler))))

	// roles actions
	private.GET("/roles", resp(e3chGroup(etcdWrapper(getRolesHandler))))
	private.POST("/role", resp(e3chGroup(etcdWrapper(createRoleHandler))))
	private.GET("/role/:name", resp(e3chGroup(getRolePermsHandler)))
	private.DELETE("/role/:name", resp(e3chGroup(etcdWrapper(deleteRoleHandler))))
	private.POST("/role/:name/permission", resp(e3chGroup(createRolePermHandler)))
	private.DELETE("/role/:name/permission", resp(e3chGroup(deleteRolePermHandler)))

	// users actions
	private.GET("/users", resp(e3chGroup(etcdWrapper(getUsersHandler))))
	private.POST("/user", resp(e3chGroup(etcdWrapper(createUserHandler))))
	private.GET("/user/:name", resp(e3chGroup(etcdWrapper(getUserRolesHandler))))
	private.DELETE("/user/:name", resp(e3chGroup(etcdWrapper(deleteUserHandler))))
	private.PUT("/user/:name/password", resp(e3chGroup(etcdWrapper(setUserPasswordHandler))))
	private.PUT("/user/:name/role/:role", resp(e3chGroup(etcdWrapper(grantUserRoleHandler))))
	private.DELETE("/user/:name/role/:role", resp(e3chGroup(etcdWrapper(revokeUserRoleHandler))))
}
