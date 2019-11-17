package routers

import (
	"github.com/VendettA01/e3w/auth"
	"github.com/VendettA01/e3w/conf"
	"github.com/VendettA01/e3w/e3ch"
	"github.com/VendettA01/e3w/resp"
	"github.com/coreos/etcd/clientv3"
	"github.com/gin-gonic/gin"
	client "github.com/soyking/e3ch"
)

const (
	etcdUsernameHeader = "X-Etcd-Username"
	etcdPasswordHeader = "X-Etcd-Password"
)

type e3chHandler func(*gin.Context, *client.EtcdHRCHYClient) (interface{}, error)

type groupHandler func(e3chHandler) resp.RespHandler

func withE3chGroup(e3chClt *client.EtcdHRCHYClient, config *conf.Config) groupHandler {
	return func(h e3chHandler) resp.RespHandler {
		return func(c *gin.Context) (interface{}, error) {
			clt := e3chClt
			if config.EtcdConf.Auth {
				var err error
				username := c.Request.Header.Get(etcdUsernameHeader)
				password := c.Request.Header.Get(etcdPasswordHeader)
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

// InitRouters initialize all served routes
// This function sets up the routes for the REST API as well as the serving of the
// static files (for the reactive web app)
func InitRouters(g *gin.Engine, config *conf.Config, e3chClt *client.EtcdHRCHYClient,
	userAuths *auth.UserAuthentications) {
	g.Static("/public", "./static/dist")
	g.GET("/", func(c *gin.Context) {
		c.File("./static/dist/index.html")
	})

	private := g.Group("/")
	if userAuths.IsEnabled {
		private.Use(authRequired(userAuths, config))
	}

	// login route cannot be protected by withAuth
	g.POST("/login", logIn(userAuths, config))

	// checkToken and logout are protected by withAuth, so
	// only authenticated users can use these
	private.GET("/checkToken", checkToken)
	private.GET("/logout", logOut)

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
