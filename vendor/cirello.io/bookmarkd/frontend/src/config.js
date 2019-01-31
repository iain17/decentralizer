const config = {
  'local': {
    fetchCreds: 'include',
    websocket: 'ws://localhost:8080/ws',
    http: 'http://localhost:8080'
  },
  'production': {
    fetchCreds: 'same-origin',
    websocket: 'wss://bookmarkd.cirello.io/ws',
    http: 'https://bookmarkd.cirello.io'
  }
}

var configuration = function () {
  switch (window.location.hostname) {
    case 'localhost':
      return config['local']
    default:
      return config['production']
  }
}

module.exports = configuration
