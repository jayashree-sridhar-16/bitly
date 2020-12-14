import React, { Component } from "react";
 
class Home extends Component {
  render() {
    return (
      <div>
        <h2>About Bitly Project</h2>
        <p>This app allows creation of shorter links. 
        Users can view statistics on the number of times short links have been accessed.
        Use create tab for creating new links.
        Use View tab for viewing already created links and their statistics</p>
      </div>
    );
  }
}
 
export default Home;