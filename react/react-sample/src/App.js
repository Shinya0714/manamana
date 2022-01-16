import React from 'react';
import axios from 'axios';

class App extends React.Component {

  handleSubmit(event) {

    axios.get('/api')
      .then((response) => { console.log(response); })
      .catch(console.error);
  }

  handleSubmit2(event) {

    axios.get('/api/py')
      .then((response) => { console.log(response); })
      .catch(console.error);
  }


  handleSubmit3(event) {

    axios.get('/api/ls')
      .then((response) => { console.log(response); })
      .catch(console.error);
  }

  render() {
    return (

      <div>
        <input type="button" value="Submit" onClick={this.handleSubmit}/>
        <input type="button" value="Submit2" onClick={this.handleSubmit2}/>
        <input type="button" value="Submit3" onClick={this.handleSubmit3}/>
      </div>
    )
  }
}

export default App;