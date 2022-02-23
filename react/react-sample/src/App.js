import React from 'react';
import axios from 'axios';
import imgLogo from './images/logo.png';
import './css/style.css'

class App extends React.Component {

  constructor() {

    super();

    this.state = {

      nav1Class: "active",
      nav2Class: "",
      nav3Class: "",
      selectedDiv: 1,
      responseValue: "",
      responseValueForSbi: "",
      responseValueForMizuho: "",

      balance: 0
    }
  }

  divHundling = prop => {

    switch(prop) {

      case 1:

        this.setState({nav1Class: "active"})
        this.setState({nav2Class: ""})
        this.setState({nav3Class: ""})
        this.setState({selectedDiv: 1})
        break;
      case 2:
          
        this.setState({nav1Class: ""})
        this.setState({nav2Class: "active"})
        this.setState({nav3Class: ""})
        this.setState({selectedDiv: 2})
        break;
      case 3:
      
        this.setState({nav1Class: ""})
        this.setState({nav2Class: ""})
        this.setState({nav3Class: "active"})
        this.setState({selectedDiv: 3})
        break;
    }
  }

  getOwner = () => {

    axios.get('/owner')
      .then((response) => {

        this.setState({balance: response.data})
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

  sbiBookBuildingSubmit = () => {

    axios.get('/api/sbiBookBuilding')
      .then((response) => {

        this.setState({responseValue: response.data})
      })
      .catch(console.error);
  }

  render() {
    return (

      <html>
        <head>
          <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous" />
        </head>
        <body>
          <div class="container">
            <img src={imgLogo} class="mt-3 mb-3" id="logo" />

            {/* nav */}
            <nav class="navbar navbar-expand-sm py-0">
              <a class={'divBorder nav-item nav-link ' + this.state.nav1Class} href="#" onClick={() => this.divHundling(1)}>資産状況</a>
              <a class={'divBorder nav-item nav-link ' + this.state.nav2Class} href="#" onClick={() => this.divHundling(2)}>申し込み</a>
              <a class={'divBorder nav-item nav-link ' + this.state.nav3Class} href="#" onClick={() => this.divHundling(3)}>設定</a>
            </nav>

            {/* div1 */}
            <div class="contentDiv p-3" style={{display: this.state.selectedDiv == 1? '' : 'none'}}>
              <input type="button" value="最新の情報に更新" class="m-3" onClick={() => this.getBalance()} />
              <br/>
              <table class="table">
              <thead>
                <tr>
                  <th scope="col"></th>
                  <th scope="col">買い付け余力</th>
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
            </div>

            {/* div2 */}
            <div class="contentDiv p-3" style={{display: this.state.selectedDiv == 2? '' : 'none'}}>
            <table class="table">
              <thead>
                <tr>
                  <th scope="col"></th>
                  <th scope="col">結果</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <th scope="row">SBI　<input type="button" value="実行" onClick={() => this.sbiBookBuildingSubmit()} /></th>
                  <td>{this.state.responseValue}</td>
                </tr>
              </tbody>
            </table>
            </div>

            {/* div3 */}
            <div class="contentDiv p-3" style={{display: this.state.selectedDiv == 3? '' : 'none'}}>
              設定
            </div>
          </div>
        </body>
        <footer>
          <p>©️MANAMANA</p>
        </footer>
      </html>
    )
  }
}

export default App;