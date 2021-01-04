import React, { Component } from 'react';
import axios from 'axios';
import { toast } from 'react-toastify';
import { Bar } from 'react-chartjs-2';
import 'react-toastify/dist/ReactToastify.css';
const { REACT_APP_BACKENDURL } = process.env;
const { REACT_APP_BASEURL } = process.env;


toast.configure();
export default class TrendData extends Component {
    constructor(props) {
        super(props);
        this.state = {
            trendData: [],
            selectedLink: {},
            selectedIndex: 0,
            finalData: [],
        }

        this.selectedLinkChange = this.selectedLinkChange.bind(this);
    }

    selectedLinkChange = (e) => {
        e.preventDefault();
        debugger;
        let newSelectedLink = this.state.trendData[e.currentTarget.id];

        this.setState({
            selectedLink : newSelectedLink,
            selectedIndex: parseInt(e.currentTarget.id),
        });
    }

    processData = (trendData)=> {
        let finalData = [];
        let minDate = new Date();
        minDate.setDate(minDate.getDate()-10);
        let maxDate = new Date();
        for(let i = 0; i< trendData.length; i++){
            let currDate = new Date(trendData[i].DateCreated);
            let formatted_date = this.getFormattedDate(trendData[i].DateCreated);
            let foundData = finalData.find(ele => ele.DateCreated === formatted_date);
            if(currDate>=minDate && currDate<= maxDate){
                if(foundData){
                    foundData.Count++;
                    foundData.Clicks += trendData[i].Clicks;
                    continue;
                }
                finalData.push({ 
                    DateCreated: formatted_date,
                    Clicks: trendData[i].Clicks,
                    date: currDate,
                    Count: 1
                });
            }
        }

        finalData.sort((a, b)=>{
            if(a.date > b.date){
                return 1;
            } else {
                return -1;
            }
        })

        return finalData;
    }

    getFormattedDate = (date) => {
        const months = ["JAN", "FEB", "MAR","APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC"];
        let current_datetime = new Date(date);
        let formatted_date = current_datetime.getDate() + " "
            + months[current_datetime.getMonth()]
            + " " + current_datetime.getFullYear()
        return formatted_date
    }

    componentDidMount(){
        axios.defaults.withCredentials = true;
        axios.get(REACT_APP_BACKENDURL+'trendData')
        .then(response => {
            console.log("Response : ", response);
            console.log("Status Code : ", response.status);
            if (response.status === 200) {
                let finalData = this.processData(response.data.Data);
                this.setState({
                    trendData: response.data.Data,
                    selectedLink: response.data.Data[0],
                    selectedIndex: 0,
                    finalData,
                });
            } else {
                this.setState({
                    trendData: this.state.trendData,
                    selectedLink: {},
                    selectedIndex: 0,
                });
                toast.error(response.statusText, {
                    position: toast.POSITION.TOP_CENTER,
                    autoClose: 2000,
                    pauseOnFocusLoss: false,
                });
            }
        }).catch(error => {
            this.setState({
                trendData: this.state.trendData,
                selectedLink: {},
                selectedIndex: 0,
            });
            toast.error(error.toString(), {
                position: toast.POSITION.TOP_CENTER,
                autoClose: 2000,
                pauseOnFocusLoss: false,
            });
        });
    }

    render() {
        let linkCount = 0;
        if(this.state.trendData && this.state.trendData.length > 0){
            linkCount = this.state.trendData.length;
        }


        let allLinks = null;
        if(this.state.trendData && this.state.trendData.length> 0){
            allLinks = this.state.trendData.map((item, index) => {
                const formatted_date = this.getFormattedDate(item.DateCreated);
                return (
                    <button onClick={this.selectedLinkChange} id={index} class={index===this.state.selectedIndex?"bitlink-item--ACTIVE":"bitlink-item--MAIN"}>
                            <label class="bitlink-item--created-date" datetime="2020-11-27">{formatted_date}</label>
                            <div class="bitlink-item--title">{item.Url}</div>
                            <div>
                                <div class="bitlink--MAIN" tabindex="-1" title={REACT_APP_BASEURL+item.Id}>{REACT_APP_BASEURL}
                                    <span class="bitlink--hash">{item.Id}</span>
                                </div>
                            </div>
                    </button>
                )
            })
        }


        let linkDetails = null;
        if(this.state.selectedLink){
            const formatted_date = this.getFormattedDate(this.state.selectedLink.DateCreated);
            linkDetails = (
                <div class="item-detail--MAIN open">
                    <div class="bitlink-detail">
                        <div>
                            <label class="item-detail--created-date" datetime="2020-11-27">CREATED: {formatted_date}</label>
                            <div>  Original URL: &nbsp;<a class="item-detail--url" href={"https://"+this.state.selectedLink.Url} target="_blank">{this.state.selectedLink.Url}</a></div>
                        </div>
                        <div class="bitlink--copy-wrapper">
                            <div class="bitlink--copy-interface">
                                Dwarfed URL: &nbsp;   
                                <a href={REACT_APP_BASEURL+this.state.selectedLink.Id} class="bitlink--copyable-text" target="_blank">
                                    <div class="bitlink--MAIN" tabindex="-1" title={REACT_APP_BASEURL+this.state.selectedLink.Id}>{REACT_APP_BASEURL}<span class="bitlink--hash">{this.state.selectedLink.Id}</span></div>
                                </a>
                            </div>
                        </div>
                        <div>
                            <div class="info-wrapper--MAIN">
                                <div class="item-detail--click-stats-wrapper">
                                    <div class="info-wrapper--user-clicks">
                                        <div class="info-wrapper--header"><span class="info-wrapper--clicks-text">{this.state.selectedLink.Clicks}</span><span class="icon clicks-icon"></span></div>
                                        <div class="item-detail--selected-day">total click</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                )
        }

        let bar = null;
        if(this.state.finalData && this.state.finalData.length > 0){
            let data = {};

            let count = this.state.finalData.map(ele => {
                return ele.Count
            })
            data.labels = this.state.finalData.map( ele => {
                return ele.DateCreated
            })

            const options = {
                scales: {
                  yAxes: [
                    {
                      ticks: {
                        beginAtZero: true,
                      },
                    },
                  ],
                },
              }
            data.datasets= [
                {
                    label: "Data Date v/s Link Counts",
                    backgroundColor: 'rgba(255,99,132,0.2)',
                    borderColor: 'rgba(255,99,132,1)',
                    borderWidth: 1,
                    hoverBackgroundColor: 'rgba(255,99,132,0.4)',
                    hoverBorderColor: 'rgba(255,99,132,1)',
                    data: count
                }
            ]
            bar = (
                <Bar
                    data={data}
                    width={100}
                    height={50}
                    options={{
                        maintainAspectRatio: false,
                        scales: {
                            yAxes: [
                              {
                                ticks: {
                                  beginAtZero: true,
                                },
                              },
                            ],
                          },
                    }}
                />
            )
        }
        return (
            <div>
                <div className="dataCardTrendData">
                        {bar}
                </div>
                        <div class="list--MAIN">
                            <span class="list--total">{linkCount} Links</span>
                            <div className="jobListContainer">
                                {allLinks}
                            </div>
                        </div>
                        {linkDetails}
            </div>
        )
    }
}
