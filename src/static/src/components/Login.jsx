import React from 'react'
import { Redirect } from 'react-router-dom'

class Login extends React.Component {
    constructor(props) {
		super(props)
		this.state = { username: '', password: '', loggedIn: false }
    }

    handleInputChange = (event) => {
		const {value, name } = event.target
		this.setState({ [name]: value })
    }

    onSubmit = (event) => {
		event.preventDefault()
		fetch('/login', {
	    	method: 'POST',
	    	body: JSON.stringify(this.state),
	    	headers: {
			'Content-Type': 'application/json'
	    	}
		})
		.then(res => {
	    	if (res.status === 200) {
				this.setState({loggedIn: true})
	    	} else {
				this.setState({loggedIn: false})
				const error = new Error(res.error)
				throw error
	    	}
		})
		.catch(err => {
	    	console.error(err)
	    	alert('Error logging in please try again')
		})
    }

    render() {
		if (this.state.loggedIn) {
			return <Redirect to ="/" />
		} else {
			return (
	    		<form onSubmit={this.onSubmit}>
				<h1>Login</h1>
				<input
		    		type="text"
		    		name="username"
		    		placeholder="Enter username"
		    		value={this.state.name}
		    		onChange={this.handleInputChange}
		    		required
				/>
				<br />
				<input
		    		type="password"
		    		name="password"
		    		placeholder="Enter password"
		    		value={this.state.password}
		    		onChange={this.handleInputChange}
		    		required
				/>
				<br />
				<input type="submit" value="Login"/>
	    		</form>
			)
		}
    }
}

export default Login;

