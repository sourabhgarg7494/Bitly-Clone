import React, { Component } from 'react';
import axios from 'axios';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { Link } from 'react-router-dom';
const { REACT_APP_BACKENDURL } = process.env;
const { REACT_APP_BASEURL } = process.env;

toast.configure();

export default class Home extends Component {
    constructor(props) {
        super(props);
        this.state = {
            longUrl: "",
            finalLongUrl:"",
            finalShortUrl:"",
            error: null,
        }
        this.longUrlChangeHandler = this.longUrlChangeHandler.bind(this);
        this.createShortUrl = this.createShortUrl.bind(this);
    }

    longUrlChangeHandler = (e) => {
        this.setState({
            longUrl : e.target.value
        })
    }

    createShortUrl = (e) => {
        e.preventDefault();
        
        let isValid = false;
        let UrlRegex = /[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)?/gi;
        let reObj = new RegExp(UrlRegex);

        if(reObj.test(this.state.longUrl)){
            isValid = true;
        }

        if(isValid===true){
            console.log(REACT_APP_BACKENDURL);
            console.log(REACT_APP_BASEURL);
            axios.defaults.withCredentials = true;
            axios.post(REACT_APP_BACKENDURL+'create', {url: this.state.longUrl, BaseUrl : REACT_APP_BASEURL})
            .then(response => {
                console.log("Response : ", response);
                console.log("Status Code : ", response.status);
                if (response.status === 200) {
                    this.setState({
                        finalLongUrl: response.data.Url,
                        finalShortUrl: response.data.ShortUrl,
                    });
                    toast.success("URL Successfully Dwarfed", {
                        position: toast.POSITION.TOP_CENTER,
                        autoClose: 2000,
                        pauseOnFocusLoss: false,
                    });
                } else {
                    this.setState({
                        finalLongUrl: "",
                        finalShortUrl: "",
                    });
                    toast.error(response.statusText, {
                        position: toast.POSITION.TOP_CENTER,
                        autoClose: 2000,
                        pauseOnFocusLoss: false,
                    });
                }
            }).catch(error => {
                this.setState({
                    finalLongUrl: "",
                    finalShortUrl: "",
                });
                toast.error(error.toString(), {
                    position: toast.POSITION.TOP_CENTER,
                    autoClose: 2000,
                    pauseOnFocusLoss: false,
                });
            });
        } else {
            this.setState({
                finalLongUrl: "",
                finalShortUrl: "",
            });
            toast.error("Input is not Valid", {
                position: toast.POSITION.TOP_CENTER,
                autoClose: 2000,
                pauseOnFocusLoss: false,
            });
        }

    }

    render() {
        return (
            <div className="dataCard">
                <div className="row">
                    <h1>Introduction</h1>
                </div>
                <div className="row">
                <p className="CardHeading">
                    This is a project for CMPE-281. It is a clone website of bitly (Named it DwarfURL). Used Amazon web services to deploy the infrastructure of the app.
                    I have opted for the below extra credit option: <br/>
                </p>
                <p className="collegeNameDiv">    
                    <b>1. Created heroku based web app for the application.</b>
                </p>
                {/* <p className="collegeNameDiv">
                    <b>2. Deployed some part of the app on GCP.</b>
                    </p> */}
                </div>
                <div className="row">
                <input onChange={this.longUrlChangeHandler} id="shorten_url" class="shorten-input" autocomplete="off" name="url" type="text" placeholder="Put your URL here.."/>
                <button onClick={this.createShortUrl} id="shorten_btn" class="button button-primary button-large shorten-button" type="submit">Create Short Link</button>
                </div>
                <div className="row">
                    <div className="EducationCardLeft">
                        <a href={this.state.finalLongUrl}>{this.state.finalLongUrl}</a>
                    </div>
                    <div className="EducationCardRight">
                        <a href={this.state.finalShortUrl}> {this.state.finalShortUrl} </a>
                    </div>
                </div>
            </div>
        )
    }
}
