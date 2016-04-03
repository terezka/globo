
function fetchMultipolygon(precision, callback) {
  $.ajax({
    url: 'http://10.0.1.3:3000/tos2/geojson/multipolygon?precision=16',
    method: 'POST',
    contentType: 'application/json',
    data: {"type":"MultiPolygon","coordinates":[[[[12.196884155273436,55.83638561341604],[12.271728515625,55.834843217637676],[12.19001770019531,55.77348342260549],[12.113800048828125,55.75146296066621],[12.01904296875,55.65279803318956],[12.057495117187498,55.606281251302114],[12.1893310546875,55.64659898563683],[12.1893310546875,55.46017083861815],[12.28271484375,55.55660246986701],[12.420043945312498,55.61558902526749],[12.4859619140625,55.61869112567042],[12.689208984375,55.5099714998319],[12.689208984375,55.55660246986701],[12.689208984375,55.68687525596441],[12.469482421875,55.71782880151228],[12.5189208984375,55.819801652442436],[12.3870849609375,55.819801652442436],[12.282028198242188,55.86336763758299],[12.2662353515625,55.84294011297761],[12.196884155273436,55.83638561341604]]]]},
    success: data => callback(data),
    error: err => callback('Call to /multipolygon failed')
  })
}

class App extends React.Component {
	constructor(props) {
    super(props)
    this.state = {
      firstDate: '',
      secondDate: '',
      text: '',
      precision: 12,
      result: null
    }
    
    this.handleInput = this.handleInput.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
  }

  handleInput(name) {
    this.setState({[name]: this.refs[name].value})
  }

  handleSubmit() {
    const { precision } = this.state
    fetchMultipolygon(precision, result => this.setState({result}))
  }

  render() {
    const { firstDate, secondDate, text, precision, result } = this.state

    return (
      <div id="app">
        <h1>Globo</h1>
        <div className="form">
          <input ref="firstDate" placeholder="Date from" className="form__date form__date--first" onChange={this.handleInput.bind(this, 'firstDate')} value={firstDate} />
          <input ref="secondDate" placeholder="Date to" className="form__date form__date--second" onChange={this.handleInput.bind(this, 'secondDate')} value={secondDate} />
          <input ref="text" placeholder="Text" className="form__text" onChange={this.handleInput.bind(this, 'text')} value={text} />
          <input ref="precision" placeholder="Precision" className="form__precision" onChange={this.handleInput.bind(this, 'precision')} value={precision} />
          <button className="form__submit" onClick={this.handleSubmit}>Submit</button>
          <Map data={result} />
        </div>
     </div>
    )
  }
}

class Map extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <div className="map">
        {this.props.data}
      </div>
    )
  }
}

ReactDOM.render(<App/>, document.getElementById('content'))	
