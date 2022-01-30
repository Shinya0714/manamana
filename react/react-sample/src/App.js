import React from 'react';
import axios from 'axios';

class App extends React.Component {

  constructor() {

    super();

    this.state = {

      responseValue: "",
    }
  }

  handleSubmit = () => {

    axios.get('/api')
      .then((response) => {

        this.setState({responseValue: response.data})
      })
      .catch(console.error);
  }

  render() {
    return (

      <div>
        <input type="button" value="Submit" onClick={this.handleSubmit}/>
        <p>{this.state.responseValue}</p>
      </div>
    )
  }
}

export default App;