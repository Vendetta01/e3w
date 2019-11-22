import React from 'react'
import { Redirect } from 'react-router-dom'

class Logout extends React.Component {
	logout = () => {
		fetch('/logout')
	    .then(res => {
			if (res.status !== 200) {
				var msg = 'Error logging out: status: ' + res.status
				console.log(msg)
				alert(msg)
			}
	    })
	}

    render() {
		this.logout()
		return <Redirect to="/login" />
    }
}

export default Logout;

