import React from 'react';
import axios from 'axios';
import imgLogo from './images/logo.png';
import './css/style.css'

class App extends React.Component {

  constructor() {

    super();

    this.state = {

      selectedDiv: 1,
      responseValue: "",
      responseValueForSbi: "",
      responseValueForMizuho: "",
      companyNameStringList: [],
      bookBuildingStringList: [],
      bookBuildingPossibleBoolList: [],
      bookBuildingPossibleBoolListForMizuho: [],
      targetCdStringList: [],
      responseValueForSbiBookBuilding: "",
      responseValueForMizuhoBookBuilding: "",

      balance: 0
    }

    this.getBalance();
    this.schedule();
  }

  getOwner = () => {

    axios.get('/owner')
      .then((response) => {

        this.setState({balance: response.data})
      })
      .catch(console.error);
  }

  schedule = () => {

    axios.get('/api/schedule')
      .then((response) => {

        this.setState({companyNameStringList: response.data.split('&')[0].split(',')})
        this.setState({bookBuildingStringList: response.data.split('&')[1].split(',')})
        this.setState({bookBuildingPossibleBoolList: response.data.split('&')[2].split(',')})
        this.setState({bookBuildingPossibleBoolListForMizuho: response.data.split('&')[3].split(',')})
        this.setState({targetCdStringList: response.data.split('&')[4].split(',')})
      })
      .catch(console.error);
  }

  getBalance = () => {

    axios.get('/api/sbiBalance')
      .then((response) => {

        this.setState({responseValueForSbi: response.data})
      })
      .catch(console.error);

      axios.get('/api/mizuhoBalance')
      .then((response) => {

        this.setState({responseValueForMizuho: response.data})
      })
      .catch(console.error);
  }

  sbiBookBuildingSubmit = (tickerSymbol, companyName) => {

    var result = window.confirm('【対象】\r\n' + '（' + tickerSymbol + '）' + companyName + '\r\n\r\n実行してもよろしいですか？');
  
    if(result) {

      axios.get('/api/sbiBookBuilding/' + tickerSymbol)
      .then((response) => {

        this.setState({responseValueForSbiBookBuilding: response.data})

        this.schedule();
      })
      .catch(console.error);
    }
  }

  mizuhoBookBuildingSubmit = (tickerSymbol, companyName) => {

    var result = window.confirm('【対象】\r\n' + '（' + tickerSymbol + '）' + companyName + '\r\n\r\n実行してもよろしいですか？');
  
    if(result) {

      axios.get('/api/mizuhoBookBuilding/' + tickerSymbol)
      .then((response) => {

        this.setState({responseValueForMizuhoBookBuilding: response.data})

        this.schedule();
      })
      .catch(console.error);
    }
  }

  checkBookoBuildingPossible(target) {

    var boolean = false;

    if(target == 'false') {

      boolean = true;
    }

    return boolean;
  }

  returnFullDate(year, month ,day) {

    var year = year;
    var month = month;
    var day = day;
    month = ('0' + month).slice(-2);
    day = ('0' + day).slice(-2);
    var fullDate = year + month + day;

    return fullDate;
  }

  render() {
    return (

      <div>
        <div className="container">
          <img src={imgLogo} className="mt-3 mb-3" id="logo" />
          <div className="contentDiv p-3">
          <table className="table" style={{display: this.state.companyNameStringList.length != 0? '' : 'none'}}>
            <thead>
              <tr>
                <th scope="col"></th>
                <th scope="col" onClick={() => this.test()}>買い付け余力</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <th scope="row">SBI</th>
                <td>{this.state.responseValueForSbi}</td>
              </tr>
              <tr>
                <th scope="row">みずほ</th>
                <td>{this.state.responseValueForMizuho}</td>
              </tr>
            </tbody>
          </table>
          <p>{this.state.scheduleListString}</p>
          <table className="table" style={{display: this.state.companyNameStringList.length != 0? '' : 'none'}}>
            <thead>
              <tr>
                <th scope="col">ブックビルディング期間</th>
                <th scope="col">証券コード</th>
                <th scope="col">新規上場企業名</th>
                <th scope="col">SBI</th>
                <th scope="col">みずほ</th>
                <th scope="col">結果</th>
              </tr>
            </thead>
            <tbody>
            {this.state.companyNameStringList.map((companyName, i) => (
              <tr>
                <td key={this.state.bookBuildingStringList[i]}>{this.state.bookBuildingStringList[i]}</td>
                <td key={this.state.targetCdStringList[i]}>{this.state.targetCdStringList[i]}</td>
                <td key={companyName}>{companyName}</td>
                <td scope="row"><input type="button" value="実行" disabled={this.checkBookoBuildingPossible(this.state.bookBuildingPossibleBoolList[i])} onClick={() => this.sbiBookBuildingSubmit(this.state.targetCdStringList[i])}/></td>
                <td scope="row"><input type="button" value="実行" disabled={this.checkBookoBuildingPossible(this.state.bookBuildingPossibleBoolListForMizuho[i])} onClick={() => this.mizuhoBookBuildingSubmit(this.state.targetCdStringList[i])}/></td>
                <td>{this.state.responseValueForMizuhoBookBuilding}</td>
              </tr>
            ))}
            </tbody>
          </table>
          <button class="btn btn-primary" type="button" disabled style={{display: this.state.companyNameStringList.length != 0? 'none' : ''}}>
            <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
            &nbsp;&nbsp;Loading...
          </button>
          </div>
        </div>
        <footer>
          <p>©️MANAMANA</p>
        </footer>
      </div>
    )
  }
}

export default App;