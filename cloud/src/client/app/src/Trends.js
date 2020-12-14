import React, { Component } from "react";
import axios from 'axios';
 
class Trends extends Component {

	state = {
		links: []
	}

	componentDidMount() {
		//Replace with lrs api
	    axios.get(`https://3cxqoqa6ci.execute-api.us-east-1.amazonaws.com/prod/links`)
	      .then(res => {
	        const links = res.data;
	        console.log(res.data);
	        this.setState({ links });
	    })
	}

	render() {
	    return (
	      <div>
	      	<div>
	        <h2>Statistics of accessed links</h2>
	        </div>
	        <p/><p/><p/><p/>
	        <div style={{align: "center"}}>
		        <table>
		        	<thead>
		        	<tr>
					    <th>Original URL</th>
					    <th>Short URL</th>
					    <th>Access Count</th>
				    </tr>
				    </thead>
				    <tbody>
			        { this.state.links.map(link => {
			        	return(
			        		<tr key={link.Short_url}>
				        		<td>
				        			<a href={link.Original_url}>{link.Original_url}</a>
				        		</td>
				        		<td>
				        			<a href={link.Redirect_url}>{link.Redirect_url}</a>
				        		</td>
				        		<td>{link.Access_count}</td>
			        		</tr>
			        	)}
			        )}
			        </tbody>
			    </table>
			</div>
	        
	      </div>
	    );
	}
}
 
export default Trends;