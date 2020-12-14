import React, { Component } from "react";
import axios from 'axios';
 
class Create extends Component {

  state ={
    url: '',
    short: '',
    showDiv: false
  }

  

  handleChange = event => {
    this.setState({ url: event.target.value });
  }

  handleSubmit = event => {
    event.preventDefault();

    const headers = {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    }

    const link = {
      "Original_url": this.state.url
    };

    //Replace with cp api
    axios.post(`https://hzb8691qnb.execute-api.us-east-1.amazonaws.com/prod/links/create`, link, headers )
      .then(res => {
        console.log(res);
        console.log(res.data);
        this.setState({ showDiv: true, short: res.data.Redirect_url});
      })
  }

  render() {
    const { url, short, showDiv } = this.state;
    const styleObj = {
      marginRight: "7px",
      marginLeft: "7px",
      minHeight: "40px",
      width: "75%",
      fontSize: "medium",
      textIndent: "10px"
    }

    const buttonObj = {
      minHeight: "45px",
      fontSize: "medium",
      width: "80px"
    }
    return (
      <div>
        <h2>Create a short url</h2>
        <p/>
        <p/>
        <p/>
        <div className="header">
          <form onSubmit={this.handleSubmit} style={{marginTop: "10px"}}>
            <label style={{fontSize: "20px"}}> Original URL:   
            <input placeholder="Enter any url" type="text" onChange={this.handleChange} 
            style={styleObj} name="url"/>
            </label>
            <button type="submit" style={buttonObj}>Create</button>
          </form>
        </div>
        <p/>
        <p/>
        <p/>
        {showDiv && (
          <div>
            <p>New URL: {short}</p>
          </div>)}

      </div>
    );
  }
}
 
export default Create;