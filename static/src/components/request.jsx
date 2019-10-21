import xhr from 'xhr'
import { message } from 'antd'

function handler(callback) {
    return function (err, response) {
        if (err) {
            message.error(err);
        } else {
            if (response && response.body) {
		let resp = ""
		try {
		    resp = JSON.parse(response.body)
		} catch (e) {
		    message.error("ERROR in handler(callback): " + e.message)
		    console.log("ERROR in handler(callback): " + e.message + "\n" + e.stack)
		    return
		}

		if (callback) {
		    callback(resp.result)
		} else {
		    let _err = new Error("function callback() not defined")
		    message.error("API error: handler(): " + _err.message)
		    console.log("ERROR in handler(callback): " + _err.message + "\n" + _err.stack)
		}
            }
        }
    }
}

function withAuth(options) {
    return Object.assign(
        options || {},
        {
            "headers": {
                "X-Etcd-Username": localStorage.etcdUsername,
                "X-Etcd-Password": localStorage.etcdPassword
            }
        }
    )
}

export function KVList(path, callback) {
    xhr.get("kv" + path + "?list", withAuth(), handler(callback))
}

export function KVGet(path, callback) {
    xhr.get("kv" + path, withAuth(), handler(callback))
}

export function KVPost(path, value, callback) {
    let bodyStr = JSON.stringify({ value: value })
    xhr.post("kv" + path, withAuth({ body: bodyStr }), handler(callback))
}

export function KVPut(path, value, callback) {
    let bodyStr = JSON.stringify({ value: value })
    xhr.put("kv" + path, withAuth({ body: bodyStr }), handler(callback))
}

export function KVDelete(path, callback) {
    xhr.del("kv" + path, withAuth(), handler(callback))
}

export function MembersGet(callback) {
    xhr.get("members", withAuth(), handler(callback))
}

export function RolesAll(callback) {
    xhr.get("roles", withAuth(), handler(callback))
}

export function RolesPost(name, callback) {
    let bodyStr = JSON.stringify({ name: name })
    xhr.post("role", withAuth({ body: bodyStr }), handler(callback))
}

export function RolesGet(name, callback) {
    xhr.get("role/" + encodeURIComponent(name), withAuth(), handler(callback))
}

export function RolesDelete(name, callback) {
    xhr.del("role/" + encodeURIComponent(name), withAuth(), handler(callback))
}

export function RolesAddPerm(name, permType, key, rangeEnd, prefix, callback) {
    let bodyStr = JSON.stringify({ perm_type: permType, key: key, range_end: rangeEnd })
    xhr.post("role/" + encodeURIComponent(name) + "/permission" + (prefix ? "?prefix" : ""), withAuth({ body: bodyStr }), handler(callback))
}

export function RolesDeletePerm(name, key, rangeEnd, callback) {
    let bodyStr = JSON.stringify({ key: key, range_end: rangeEnd })
    xhr.del("role/" + encodeURIComponent(name) + "/permission", withAuth({ body: bodyStr }), handler(callback))
}

export function UsersAll(callback) {
    xhr.get("users", withAuth(), handler(callback))
}

export function UsersPost(name, callback) {
    let bodyStr = JSON.stringify({ name: name })
    xhr.post("user", withAuth({ body: bodyStr }), handler(callback))
}

export function UsersGet(name, callback) {
    xhr.get("user/" + encodeURIComponent(name), withAuth(), handler(callback))
}

export function UsersDelete(name, callback) {
    xhr.del("user/" + encodeURIComponent(name), withAuth(), handler(callback))
}

export function UsersGrantRole(name, role, callback) {
    xhr.put("user/" + encodeURIComponent(name) + "/role/" + encodeURIComponent(role), withAuth(), handler(callback))
}

export function UsersRovokeRole(name, role, callback) {
    xhr.del("user/" + encodeURIComponent(name) + "/role/" + encodeURIComponent(role), withAuth(), handler(callback))
}

export function UsersChangePassword(name, password, callback) {
    let bodyStr = JSON.stringify({ password: password })
    xhr.put("user/" + encodeURIComponent(name) + "/password", withAuth({ body: bodyStr }), handler(callback))
}

