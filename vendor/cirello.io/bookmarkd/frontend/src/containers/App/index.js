// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { Route, Link, withRouter } from 'react-router-dom'
import { Nav, NavItem, Navbar } from 'react-bootstrap'
import HomePage from '../HomePage'
import PostPage from '../PostPage'
import {startWebsocket} from './actions'

class App extends React.Component {
  componentDidMount () {
    this.props.startWebsocket()
  }

  render () {
    return (
      <div>
        <Navbar>
          <Navbar.Header>
            <Navbar.Brand>
              <Link to='/'>Home</Link>
            </Navbar.Brand>
          </Navbar.Header>
          <Nav>
            <NavItem componentClass={Link} eventKey={1}
              href='/post' to='/post'>new bookmark</NavItem>
          </Nav>
        </Navbar>

        <Route exact path='/' component={HomePage} />
        <Route path='/post' component={PostPage} />
      </div>
    )
  }
}

function s2p (state) { return {} }

function d2p (dispatch) {
  return bindActionCreators({ startWebsocket }, dispatch)
}

export default withRouter(connect(s2p, d2p)(App))
