import React from 'react'
import { Box } from 'react-polymer-layout'
import { Breadcrumb } from 'antd'
import { KVList } from './request'
import KeyValueCreate from './KeyValueCreate'
import KeyValueItem from './KeyValueItem'
import KeyValueSetting from './KeyValueSetting'
import { CommonPanel } from './utils'

class KeyValue extends React.Component {
    // states:
    // - dir: the full path of current dir, eg. / or /abc/def
    // - menus: components of Breadcrumb, including path (to another dir, using in url hash) and name
    // - list: the key under the dir, get from api

    constructor(props) {
	super(props)
	this.state = { dir: "", menus: [], list: [], setting: false, currentKey: "" }
    }

    _isRoot = () => {
        return this.state.dir === "/"
    }

    _parseList = (list) => {
        list = list || []
        // sorted dir and normal kv
        list.sort((l1, l2) => { return l1.is_dir === l2.is_dir ? l1.key > l2.key : l1.is_dir ? -1 : 1 })
        // trim prefix of dir, get the relative path, +1 for /
        let prefixLen = this.state.dir.length + (this._isRoot() ? 0 : 1)
        list.forEach(l => {
            l.key = l.key.slice(prefixLen)
        })
        this.setState({ list: list })
    }

    // dir should be / or /abc/def
    _ParseDir = (dir) => {
        let menus = [{ path: "/", name: "ROOT" }]
        if (dir !== "/") {
            let keys = dir.split("/")
            for (let i = 1; i < keys.length; i++) {
                // get the full path of every component
                menus.push({ path: keys.slice(0, i + 1).join("/"), name: keys[i] })
            }
        }
        KVList(dir, this._parseList)
        return { dir: dir, menus: menus }
    }

    // list current dir and using KeyValueSetting
    _fetch = (dir) => {
        this.setState(this._ParseDir(dir))
        this.setState({ setting: false })
    }

    // change url
    _redirect = (dir) => {
	this.props.history.push("/kv" + dir)
    }

    _fullKey = (subKey) => (this._isRoot() ? "/" : this.state.dir + "/") + subKey;

    // callback for clicking KeyValueItem to enter a new dir
    _enter = (subKey) => {
	console.log("KeyValue: _enter():")
	console.log(subKey)
	console.log(this._fullKey(subKey))
        this._redirect(this._fullKey(subKey))
    }

    // callback for clicking KeyValueItem to set the kv
    _set = (subKey) => {
        let list = this.state.list
        list.forEach(l => {
            if (l.key === subKey) { l.selected = true } else { l.selected = false }
        })
        this.setState({ setting: true, currentKey: this._fullKey(subKey), list: list })
    }

    // call back for clicking KeyValueItem again
    _unset = (subKey) => {
        let list = this.state.list
        list.forEach(l => {
            l.selected = false
        })
        this.setState({ setting: false, list: list })
    }

    // callback for deleting a key in KeyValueItem
    _delete = () => {
        this._fetch(this.state.dir)
    }

    // callback for creating kv or dir
    _update = () => {
        this._fetch(this.state.dir)
    }

    // callback for delete currentDir and enter previous dir
    _back = () => {
        let menus = this.state.menus
        let targetPath = (menus[menus.length - 2] || menus[0]).path
        this._redirect(targetPath)
    }

    // refresh the page with new path in url
    _refresh = (path) => {
	this._fetch(path)
    }

    _getDirFromProps = () => {
	let _path ="/"
	if (typeof this.props.match.url !== "undefined") {
	    _path = "/" + this.props.location.pathname.replace(/^\/kv/, "")
	    _path = _path.replace(/^\/+/, "/")
	}
	return _path
    }

    _setDirFromProps = () => {
	let _path = this._getDirFromProps()
	if (_path !== this.state.dir) {
	    this._refresh(_path)
	}
    }

    componentDidMount() {
	this._setDirFromProps()
    }

    componentDidUpdate(prevProps, prevState) {
	this._setDirFromProps()
    }

    render() {
        let currentKey = this.state.currentKey
        return (
            <Box vertical>
                <Box style={{ paddingBottom: 15 }}>
                    <Breadcrumb>
                        {
                            this.state.menus.map(
                                m => (<Breadcrumb.Item onClick={() => this._redirect(m.path) } href={"#/kv" + m.path} key={m.path}>{m.name}</Breadcrumb.Item>)
                            )
                        }
                    </Breadcrumb>
                </Box>
                <Box>
                    <Box vertical style={{ width: 400, paddingRight: 20 }}>
                        <Box vertical >
                            {
                                this.state.list.map(
                                    l => (<KeyValueItem key={l.key} enter={this._enter} set={this._set} unset={this._unset} info={l} />)
                                )
                            }
                        </Box>
                    </Box>

                    {this.state.setting ?
                        (<CommonPanel hint={currentKey} color="#cce4f6"><KeyValueSetting currentKey={currentKey} delete={this._delete} /></CommonPanel>) :
                        (<CommonPanel hint="CREATE"><KeyValueCreate update={this._update} back={this._back} dir={this.state.dir} fullKey={this._fullKey} /></CommonPanel>) }
                </Box>
            </Box >
        )
    }
}

export default KeyValue

