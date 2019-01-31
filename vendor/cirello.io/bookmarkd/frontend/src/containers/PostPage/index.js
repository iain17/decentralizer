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
import { push } from 'react-router-redux'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import {
Button,
Col,
ControlLabel,
FormControl,
FormGroup,
Grid,
HelpBlock,
PageHeader,
Row
} from 'react-bootstrap'
import { loadBookmark, newBookmark } from './actions'

export class PostPage extends React.PureComponent { // eslint-disable-line react/prefer-stateless-function
  constructor (props) {
    super(props)

    this.state = {
      submitting: false,
      bookmark: {title: '', url: ''}
    }
    this.submit = this.submit.bind(this)
    this.update = this.update.bind(this)
  }
  componentDidMount () {
    var url = this.hasPredefinedURL()
    if (url && this.state.bookmark.url === '') {
      loadBookmark(url, (b) => {
        this.setState({bookmark: b})
      })
    }
  }

  hasPredefinedURL () {
    const params = new URLSearchParams(window.location.search)
    return params.get('url')
  }

  submit (e) {
    e.preventDefault()
    this.setState({
      submitting: true
    }, () => {
      this.props.newBookmark(this.state.bookmark)
    })
  }

  update (k, v) {
    var bookmark = {...this.state.bookmark}
    bookmark[k] = v
    this.setState({bookmark})
  }
  render () {
    if (this.state.submitting) {
      return (
        <div>
          <Grid>
            <Row>
              <Col>
              storing new bookmark...
              </Col>
            </Row>
          </Grid>
        </div>
      )
    }
    return (
      <div>
        <Grid>
          <Row>
            <Col>
              <PageHeader>
                New Link
              </PageHeader>
            </Col>
          </Row>
          <Row>
            <Col>
              {
                this.hasPredefinedURL() && this.state.bookmark.url === ''
                ? 'loading...'
                : ''
              }
              <form onSubmit={(e) => this.submit(e)}>
                <FieldGroup
                  id='formControlsText'
                  type='text'
                  placeholder='Title'
                  onChange={(e) => this.update('title', e.target.value)}
                  value={this.state.bookmark.title} />
                <FieldGroup
                  id='formControlsText'
                  type='text'
                  placeholder='URL'
                  onChange={(e) => this.update('url', e.target.value)}
                  value={this.state.bookmark.url} />
                <Button type='submit'>Add</Button>
              </form>
            </Col>
          </Row>
        </Grid>
      </div>
    )
  }
}

function FieldGroup ({ id, label, help, ...props }) {
  return (
    <FormGroup controlId={id}>
      <ControlLabel>{label}</ControlLabel>
      <FormControl {...props} />
      {help && <HelpBlock>{help}</HelpBlock>}
    </FormGroup>
  )
}

function s2p (state) {
  return {
    bookmarks: state.bookmarks ? state.bookmarks : []
  }
}

function d2p (dispatch) {
  return bindActionCreators(
    {
      newBookmark,
      changePage: (target) => push(target)
    }, dispatch)
}

export default connect(s2p, d2p)(PostPage)
