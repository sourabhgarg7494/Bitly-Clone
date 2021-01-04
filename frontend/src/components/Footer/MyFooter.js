import React, { Component } from 'react'
import Footer from 'rc-footer';
import 'rc-footer/assets/index.css';
import '../../App.css';
import githublogo from '../../images/githublogo.png';
import linkedinlogo from '../../images/linkedinlogo.png';

export default class MyFooter extends Component {
    render() {
        return (
            <Footer
    columns={[
      {
        icon: (
          <img src={linkedinlogo} alt="LinkedIn" />
        ),
        title: <a target="_blank" href="https://www.linkedin.com/in/sourabhgarg7494/">sourabhgarg7494</a>,
      },
      {
        icon: (
          <img src={githublogo} alt="GitHub" />
        ),
        title: <a target="_blank" href="https://github.com/sourabhgarg7494">sourabhgarg7494</a>,
      },
    ]}
    bottom="CMPE-281 Bitly Clone (Individual Project)"
    theme="light"
  />
        )
    }
}
