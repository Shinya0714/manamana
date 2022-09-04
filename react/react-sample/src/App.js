import React from 'react';
import axios from 'axios';
import imgLogo from './images/logo.png';
import kikanGai from './images/kikanGai.png';
import kanryo from './images/kanryo.png';
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
      dataList: [],

      balance: 0
    }

    // this.getBalance();
    this.schedule();
  }

  // getOwner = () => {

  //   axios.get('/owner')
  //     .then((response) => {

  //       this.setState.target({balance: response.data})
  //     })
  //     .catch(console.error);
  // }

  schedule = () => {

    axios.get('/api/schedule')
      .then((response) => {

        var jsonObject = JSON.parse(response.data.outputJson)

        this.setState({dataList: jsonObject})
        this.setState({responseValueForSbi: response.data.sbiBalance})

        console.log(response.data.outputJson);
        console.log(response.data.sbiBalance);
        console.log(response.data.mizuhoBalance);
        console.log(this.state.dataList)
        // const allResponse = {};
        // allResponse.servers = [];
        // for (const res of responses) {
        //     // 必要ならこの辺でレスポンスを加工
        //     allResponse.servers.push(res.data);
        // }

        // console.log(allResponse.servers[1]);
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

    console.log(target)

    var result = 'false';

    if(target == 'true') {

      result = 'true';
    }else if(target == 'kikanGai') {

      result = 'kikanGai';
    }else if(target == 'kanryo') {

      result = 'kanryo';
    }

    return result;
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
      {/* {this.state.dataList.map((data, i) => (
        <p>{data.CompanyNameString}</p>
      ))} */}
      <div className="container">
        <img src={imgLogo} className="mt-3 mb-3" id="logo" />
        <div className="contentDiv p-3">
        <table className="table" style={{display: this.state.dataList.length != 0? '' : 'none'}}>
          <thead>
            <tr>
              <th scope="col">証券会社</th>
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
        <table className="table" style={{display: this.state.dataList.length != 0? '' : 'none'}}>
          <thead>
            <tr>
              <th scope="col" className='text-center'>ブックビルディング期間</th>
              <th scope="col" className='text-center'>証券コード</th>
              <th scope="col" className='text-center'>新規上場企業名</th>
              <th scope="col" className='text-center'>SBI</th>
              <th scope="col" className='text-center'>みずほ</th>
              <th scope="col" className='text-center'>結果（SBI）</th>
              <th scope="col" className='text-center'>結果（みずほ）</th>
            </tr>
          </thead>
          <tbody>
          {this.state.dataList.map((target, i) => (
            <tr className={this.state.dataList[i] == '---' ? 'table-secondary' : ''}>
              <td className='text-center' key={target.BookBuildingString}>{target.BookBuildingString}</td>
              <td className='text-center' key={target.TargetCdString}>{target.TargetCdString}</td>
              <td className='text-center' key={target.CompanyNameString}>{target.CompanyNameString}</td>
              <td className='text-center' scope="row"><input type="button" value="実行" disabled={(this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'kikanGai' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'false' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'kanryo') ? true: false} onClick={() => this.sbiBookBuildingSubmit(target.TargetCdString, target.CompanyNameString)}/></td>
              <td className='text-center' scope="row"><input type="button" value="実行" disabled={(this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kikanGai' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'false') ? true: false} onClick={() => this.mizuhoBookBuildingSubmit(target.TargetCdString, target.CompanyNameString)}/></td>
              {(() => {
                if (target.BookBuildingString == '---') {
                  return <td className='text-center'></td>;
                }
              })()}
              {(() => {
                if (target.BookBuildingString == '---') {
                  return <td className='text-center'></td>;
                }
              })()}
              {(() => {
                if (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'kikanGai') {
                  return <td className='text-center'><img src={kikanGai} id="statusImageForKikanGai" /></td>;
                }
              })()}
              {(() => {
                if ((this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'kanryo' || this.state.responseValueForSbiBookBuilding == 'ブックビルディングのお申し込みを受付いたしました。')) {
                  return <td className='text-center'><img src={kanryo} id="statusImageForKanryo" /></td>;
                }
              })()}
              {(() => {
                if (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kikanGai') {
                  return <td className='text-center'><img src={kikanGai} id="statusImageForKikanGai" /></td>;
                } else {
                  return <td className='text-center'>{target.BookBuildingPossibleBoolStringForMizuho}</td>;
                }
              })()}
            </tr>
          ))}
          </tbody>
        </table>
        <button className="btn btn-primary" type="button" disabled style={{display: this.state.dataList.length != 0? 'none' : ''}}>
          <span className="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
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