import React, { Component } from 'react';
import {Route} from 'react-router-dom';
import MyFooter from './Footer/MyFooter';
import Home from './Home/Home';
import Navbar from './Navbar/Navbar'
import TrendData from './TrendData/TrendData';

export default class Main extends Component {
    render() {
        return (
            <div>
                <Route path="/" component={Navbar}/>
                <Route path="/" component={Home} exact={true}/>
                <Route path="/trendData" component={TrendData}/>
                <Route path="/" component={MyFooter} exact={true}/>
            </div>
        )
    }
}
