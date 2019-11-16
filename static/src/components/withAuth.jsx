import React from 'react';
import { Redirect } from 'react-router-dom';

function withAuth(ComponentToProtect) {
    return class extends React.Component {
	constructor() {
	    super()
	    this.state = { loading: true, redirect: false , url: "" }
	}

	/*static getDerivedStateFromProps(props, state) {
	    var new_url = props.location.pathname
	    if (state.url != new_url) {
		fetch('/checkToken')
		.then(res => {
		    if (res.status === 200) {
			return { loading: false, redirect: false, url: new_url }
		    } else {
			const error = new Error(res.error)
			throw error
		    }
		})
		.catch(err => {
		    console.log(err)
		    return { loading: false, redirect: true, url: new_url }
		})
		//return { loading: false, redirect: true, url: new_url }
	    }
	}*/

	componentDidMount() {
	    fetch('/checkToken')
	    .then(res => {
		if (res.status === 200) {
		    this.setState({ loading: false })
		} else {
		    const error = new Error(res.error)
		    throw error
		}
	    })
	    .catch(err => {
		console.log(err)
		this.setState({ loading: false, redirect: true })
	    })
	}

	render() {
	    const { loading, redirect } = this.state;
	    if (loading) {
		return null
	    }
	    if (redirect) {
		return <Redirect to="/login" />
	    }
	    return (
		<React.Fragment>
		    <ComponentToProtect {...this.props} />
		</React.Fragment>
	    )
	}
    }
}

export default withAuth;

