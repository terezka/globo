
class App extends React.Component {
	constructor(props) {
    super(props)
    this.state = {
      firstDate: '',
      secondDate: '',
      text: ''
    }

    this.handleInput = this.handleInput.bind(this)
  }

  handleInput(name) {
    this.setState({[name]: this.refs[name].value})
  }

  render() {
    const { firstDate, secondDate, text } = this.state

    return (
      <div className="app">
        <h1>Globo</h1>
        <div className="form">
          <input ref="firstDate" className="form__date form__date--first" onChange={this.handleInput.bind(this, 'firstDate')} value={firstDate} />
          <input ref="secondDate" className="form__date form__date--second" onChange={this.handleInput.bind(this, 'secondDate')} value={secondDate} />
          <input ref="text" className="form__text" onChange={this.handleInput.bind(this, 'text')} value={firstDate} />
        </div>
     </div>
    )
  }
}

ReactDOM.render(<App/>, document.getElementById('content'))	
