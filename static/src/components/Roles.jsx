import React from 'react'
import AuthPanel from './AuthPanel'
import { RolesAll, RolesPost, RolesDelete } from './request'
import RolesSetting from './RolesSetting'

class Roles extends React.Component {
    constructor(props) {
	super(props)
	this.state = { roles: [] }
    }

    _getRolesDone = (result) => {
        this.setState({ roles: result || [] })
    }

    _getRoles = () => {
        RolesAll(this._getRolesDone)
    }

    _createRoleDone = (result) => {
        this._getRoles()
    }

    _createRole = (name) => {
        RolesPost(name, this._createRoleDone)
    }

    _deleteRoleDone = (result) => {
        this._getRoles()
    }

    _deleteRole = (name) => {
        RolesDelete(name, this._deleteRoleDone)
    }

    componentDidMount() {
        this._getRoles()
    }

    componentWillReceiveProps(nextProps) {
        this._getRoles()
    }

    _setting = (name) => {
        return <RolesSetting name={name} />
    }

    render() {
        return (
            <AuthPanel title="ROLES" items={this.state.roles} create={this._createRole} setting={this._setting} delete={this._deleteRole} />
        )
    }
}

export default Roles

