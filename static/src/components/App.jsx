import React from 'react'
import { HashRouter, Link, Route, Switch } from 'react-router-dom'
import { Menu, Icon } from 'antd'
import { Box } from 'react-polymer-layout'
import KeyValue from './KeyValue'
import Members from './Members'
import Roles from './Roles'
import Users from './Users'
import Setting from "./Setting"

class App extends React.Component {
    constructor(props) {
	super(props);
	this.state = { menu: "" };
    }

    _getMenu = () => {
        let parts = window.location.hash.split("/")
        let menu = "kv"
        if (parts.length > 1) {
	    menu = parts[1]
        }
        return menu
    }

    _changeMenu = () => {
        this.setState({ menu: this._getMenu() })
    }

    componentDidMount() {
        this._changeMenu()
    }

    componentDidUpdate(prevProps, prevState) {
	if (this.state.menu !== this._getMenu()) {
	    this._changeMenu()
	}
    }

    handleClick = (e) => {
	if (e.key !== this.state.menu) {
	    this._changeMenu()
	}
    }

    render() {
        return (
	    <HashRouter>
            <Box centerJustified>
                <Box vertical style={{ width: 1000 }}>
                    <Box style={{ padding: 20, borderBottom: "1px #E6E6E6 solid" }}>
                        <Box center centerJustified onClick={() => { window.location.hash = "#/" } }
                            style={{
                                fontSize: 25, fontWeight: 700, marginRight: 20, paddingRight: 20,
                                borderStyle: "solid", borderWidth: "0px 2px 0px 0px", borderColor: "#ddd",
                                cursor: "pointer"
                            }}>
                            E·3·W
                        </Box>
                        <Menu onClick={this.handleClick}
                            selectedKeys={[this.state.menu]}
                            mode="horizontal"
                            style={{ fontWeight: 700, fontSize: 14 }}
                            >
                            <Menu.Item key="kv">
                                <Icon type="menu-fold" /><span>KEY / VALUE</span>
				<Link to="/kv" />
                            </Menu.Item>
                            <Menu.Item key="members">
                                <Icon type="tags" /><span>MEMBERS</span>
				<Link to="/members" />
                            </Menu.Item>
                            <Menu.SubMenu key="auth" title={<span><Icon type="team" />AUTH</span>}>
                                <Menu.Item key="roles">
				    ROLES
				    <Link to="/roles" />
				</Menu.Item>
                                <Menu.Item key="users">
				    USERS
				    <Link to="/users" />
				</Menu.Item>
                            </Menu.SubMenu>
                            <Menu.Item key="setting">
                                <Icon type="setting" /><span>SETTING</span>
				<Link to="/setting" />
                            </Menu.Item>
                        </Menu>
                    </Box>
                    <div style={{ paddingTop: 20 }}>
			<AppBody />
                    </div>
                </Box>
            </Box>
	    </HashRouter>
        );
    }
}


class AppBody extends React.Component {
    render() {
	return (
	    <Switch>
		<Route exact path="/kv" component={KeyValue} />
		<Route path="/kv/:path" component={KeyValue} />
		<Route path="/members" component={Members} />
		<Route path="/roles" component={Roles} />
		<Route path="/users" component={Users} />
		<Route path="/setting" component={Setting} />
	    </Switch>
	)
    }
}


export default App


