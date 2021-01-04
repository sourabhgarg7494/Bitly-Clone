import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import logo from '../../images/logo.png';
import '../../App.css';

export default class Navbar extends Component {
    render() {
        return (
            <div>
                <nav className="navbar navbar-inverse">
                    <div className="container-fluid">
                        <div className="navbar-header">
                            <Link to="/"><img className="navbar-brand__logo-full" alt="DwarfURL LOGO" src={logo} /></Link>
                        </div>
                        <ul className="nav navbar-nav navbar-right">
                            <li><Link to="/">Home</Link></li>
                            <li><Link to="/trendData">Trend Data</Link></li>
                        </ul>
                    </div>
                </nav>
            </div>
        )
    }
}
