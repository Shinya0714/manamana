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

      responseValue: "",
      responseValueForSbi: "",
      responseValueForMizuho: "",
      responseValueForSmbc: "",
      responseValueForRakuten: "",

      companyNameStringList: [],

      bookBuildingStringList: [],
      bookBuildingPossibleBoolList: [],
      bookBuildingPossibleBoolListForMizuho: [],

      targetCdStringList: [],
      targetPriceStringList: [],

      responseValueForSbiBookBuilding: "",
      responseValueForMizuhoBookBuilding: "",
      
      dataList: [],

      sbiBalanceRenderingFlg: false,
      mizuhoBalanceRenderingFlg: false,
      smbcBalanceRenderingFlg: false,
      rakutenBalanceRenderingFlg: false,
      
      scheduleRenderingFlg: false,
      
      balance: 0
    }
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

        console.log('this.state.dataList:' + this.state.dataList);

        if(jsonObject != null) {

          this.setState({scheduleRenderingFlg: true})
        }

        console.log('schedule success');
      })
      .then(() => {

        alert('スケジュールの更新が完了しました。')
      })
      .catch((error) => {
        console.error('schedule err:', error);
      })
  }

  getBalance = () => {

    axios.get('/api/balance')
    .then((response) => {

      var sbiBalance = response.data.sbiBalance
      var mizuhoBalance = response.data.mizuhoBalance
      var smbcBalance = response.data.smbcBalance
      var rakutenBalance = response.data.rakutenBalance

      this.setState({responseValueForSbi: sbiBalance})
      this.setState({responseValueForMizuho: mizuhoBalance})
      this.setState({responseValueForSmbc: smbcBalance})
      this.setState({responseValueForRakuten: rakutenBalance})


      if(sbiBalance != null) {

        this.setState({sbiBalanceRenderingFlg: true})
      }

      if(mizuhoBalance != null) {

        this.setState({mizuhoBalanceRenderingFlg: true})
      }

      if(smbcBalance != null) {

        this.setState({smbcBalanceRenderingFlg: true})
      }

      if(rakutenBalance != null) {

        this.setState({rakutenBalanceRenderingFlg: true})
      }
      
      console.log('getBalance success');
    })
    .catch((error) => {
      console.error('getBalance err:', error);
    })
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
      })
      .then(() => {

        this.schedule()
      })
      .catch(console.error);
    }
  }

  smbcBookBuildingSubmit = (tickerSymbol, companyName) => {

    var result = window.confirm('【対象】\r\n' + '（' + tickerSymbol + '）' + companyName + '\r\n\r\n実行してもよろしいですか？');
    if(result) {

      axios.get('/api/smbcBookBuilding/' + tickerSymbol)
      .then((response) => {

        this.setState({responseValueForSmbcBookBuilding: response.data})
      })
      .then(() => {

        this.schedule()
      })
      .catch(console.error);
    }
  }

  checkBookoBuildingPossible(target) {

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

  update() {

    this.getBalance();
    this.schedule();
  }

  render() {
    return (
      <div>
      <div className="container">
        <div className="row mt-3">
          <div className='col-md-2 col-sm-12 text-center'>
            <img src={imgLogo} id="logo" />
            <br/>
            <input className="mt-3 mb-3 btn btn-outline-dark" type="button" value="最新の情報に更新" onClick={() => this.update()} id="updateButton" />
          </div>
          <div className='col-md-10 col-sm-12'>
            <table className="table table-hover">
              <thead>
                <tr>
                  <th scope="col">システムアクセスログ</th>
                  <td></td>
                </tr>
              </thead>
              <tbody>
                <tr scope="row">
                  <th>2021/10/17 18:19:46</th>
                  <td>スケジュールの更新が押下されました。</td>
                </tr>
                <tr scope="row">
                <th>2021/10/17 18:19:46</th>
                  <td>スケジュールの更新が押下されました。</td>
                </tr>
                <tr scope="row">
                <th>2021/10/17 18:19:46</th>
                  <td>スケジュールの更新が押下されました。</td>
                </tr>
                <tr scope="row">
                <th>2021/10/17 18:19:46</th>
                  <td>スケジュールの更新が押下されました。</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div className="contentDiv p-3">
        <table className="table">
          <thead>
            <tr>
              <th scope="col">証券会社</th>
              <th scope="col" onClick={() => this.test()}>買い付け余力</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <th className='align-middle' scope="row">SBI</th>
              <td>
                {this.state.responseValueForSbi}
                <button className="btn btn-primary" type="button" disabled style={{display: this.state.sbiBalanceRenderingFlg? 'none' : ''}}>
                  <span className="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                  &nbsp;&nbsp;Loading...
                </button>
              </td>
            </tr>
            <tr>
              <th className='align-middle' scope="row">みずほ</th>
              <td>
                {this.state.responseValueForMizuho}
                <button className="btn btn-primary" type="button" disabled style={{display: this.state.mizuhoBalanceRenderingFlg? 'none' : ''}}>
                  <span className="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                  &nbsp;&nbsp;Loading...
              </button>
              </td>
            </tr>
            <tr>
              <th className='align-middle' scope="row">SMBC</th>
              <td>
                {this.state.responseValueForSmbc}
                <button className="btn btn-primary" type="button" disabled style={{display: this.state.smbcBalanceRenderingFlg? 'none' : ''}}>
                  <span className="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                  &nbsp;&nbsp;Loading...
              </button>
              </td>
            </tr>
            <tr>
              <th className='align-middle' scope="row">楽天</th>
              <td>
                {this.state.responseValueForRakuten}
                <button className="btn btn-primary" type="button" disabled style={{display: this.state.rakutenBalanceRenderingFlg? 'none' : ''}}>
                  <span className="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                  &nbsp;&nbsp;Loading...
              </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div className="contentDiv p-3 mt-3 table-responsive-sm">
        <table className="table">
          <thead>
            <tr>
              <th scope="col" className='text-center'>ブックビルディング期間</th>
              <th scope="col" className='text-center'>証券<br/>コード</th>
              <th scope="col" className='text-center'>仮条件<br/>（円）</th>
              <th scope="col" className='text-center'>新規上場企業名</th>
              <th scope="col" className='text-center'>SBI</th>
              <th scope="col" className='text-center'>みずほ</th>
              <th scope="col" className='text-center'>SMBC</th>
              <th scope="col" className='text-center'>楽天</th>
              <th scope="col" className='text-center'>結果<br/>（SBI）</th>
              <th scope="col" className='text-center'>結果<br/>（みずほ）</th>
              <th scope="col" className='text-center'>結果<br/>（SMBC）</th>
              <th scope="col" className='text-center'>結果<br/>（楽天）</th>
            </tr>
          </thead>
          <tbody style={{display: this.state.scheduleRenderingFlg? '' : 'none'}} >
          {this.state.dataList.map((target, i) => (
            <tr className={target.BookBuildingString == '---' ? 'table-secondary' : ''} key={i}>
              <td className='text-center align-middle' key={target.BookBuildingString}>{target.BookBuildingString}</td>
              <td className='text-center align-middle' key={target.TargetCdString}>{target.TargetCdString}</td>
              <td className='text-center align-middle' key={target.TargetPriceString}>{target.TargetPriceString}</td>
              <td className='text-center align-middle' key={target.CompanyNameString}>{target.CompanyNameString}</td>
              <td className='text-center align-middle' scope="row">
                <input type="button" value="実行" disabled={(this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'kikanGai' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'false' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'kanryo') ? true: false} onClick={() => this.sbiBookBuildingSubmit(target.TargetCdString, target.CompanyNameString)}/>
              </td>
              <td className='text-center align-middle' scope="row">
                <input type="button" value="実行" disabled={(this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kikanGai' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'false' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kanryo') ? true: false} onClick={() => this.mizuhoBookBuildingSubmit(target.TargetCdString, target.CompanyNameString)}/>
              </td>
              <td className='text-center align-middle' scope="row">
                <input type="button" value="実行" disabled={(this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSmbc) == 'kikanGai' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSmbc) == 'false') ? true: false} onClick={() => this.smbcBookBuildingSubmit(target.TargetCdString, target.CompanyNameString)}/>
              </td>
              <td className='text-center align-middle' scope="row">
                <input type="button" value="実行" disabled={(this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kikanGai' || this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'false') ? true: false} onClick={() => this.mizuhoBookBuildingSubmit(target.TargetCdString, target.CompanyNameString)}/>
              </td>
              <td className='text-center align-middle' style={{display: target.BookBuildingString == '---'? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: target.BookBuildingString == '---'? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: target.BookBuildingString == '---'? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: target.BookBuildingString == '---'? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'kikanGai'? '' : 'none'}}>
                <img src={kikanGai} id="statusImageForKikanGai" />
              </td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForSbi) == 'kanryo' || this.state.responseValueForSbiBookBuilding == 'ブックビルディングのお申し込みを受付いたしました。')? '' : 'none'}}>
                <img src={kanryo} id="statusImageForKanryo" />
              </td>
              <td className='text-center align-middle' style={{display: this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kikanGai'? '' : 'none'}}>
                <img src={kikanGai} id="statusImageForKikanGai" />
              </td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kanryo' || this.state.responseValueForMizuhoBookBuilding == 'ブックビルディングのお申し込みを受付いたしました。')? '' : 'none'}}>
                <img src={kanryo} id="statusImageForKanryo" />
              </td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'false')? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'true')? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kikanGai'? '' : 'none'}}>
                <img src={kikanGai} id="statusImageForKikanGai" />
              </td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kanryo' || this.state.responseValueForMizuhoBookBuilding == 'ブックビルディングのお申し込みを受付いたしました。')? '' : 'none'}}>
                <img src={kanryo} id="statusImageForKanryo" />
              </td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'false')? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'true')? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kikanGai'? '' : 'none'}}>
                <img src={kikanGai} id="statusImageForKikanGai" />
              </td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'kanryo' || this.state.responseValueForMizuhoBookBuilding == 'ブックビルディングのお申し込みを受付いたしました。')? '' : 'none'}}>
                <img src={kanryo} id="statusImageForKanryo" />
              </td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'false')? '' : 'none'}}></td>
              <td className='text-center align-middle' style={{display: (this.checkBookoBuildingPossible(target.BookBuildingPossibleBoolStringForMizuho) == 'true')? '' : 'none'}}></td>
            </tr>
          ))}
          </tbody>
        </table>
        <p>
        <button className="btn btn-primary" type="button" disabled style={{display: this.state.scheduleRenderingFlg? 'none' : ''}}>
          <span className="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
          &nbsp;&nbsp;Loading...
        </button>
        </p>
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