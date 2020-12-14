import React, { Component } from "react";
import {
  Route,
  NavLink,
  HashRouter
} from "react-router-dom";
import Home from "./Home";
import Create from "./Create";
import Trends from "./Trends";
 
class App extends Component {
  
  render() {
    const styleObj = {
      textAlign: "center"
    }
    return (
      <HashRouter>
        <div>
          <h1 style={{styleObj}}>Bitly</h1>
          <ul className="header">
            <li><NavLink exact to="/">Home</NavLink></li>
            <li><NavLink to="/create">Create</NavLink></li>
            <li><NavLink to="/trends">View</NavLink></li>
          </ul>
          <div className="content">
            <Route exact path="/" component={Home}/>
            <Route path="/create" component={Create}/>
            <Route path="/trends" component={Trends}/>
             
          </div>
        </div>
      </HashRouter>
    );
  }
}
 
export default App;
