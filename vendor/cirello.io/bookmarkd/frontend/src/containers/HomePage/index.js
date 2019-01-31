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
import {
  Col,
  FormControl,
  FormGroup,
  Glyphicon,
  Grid,
  Label,
  Panel,
  Row,
  ToggleButton,
  ToggleButtonGroup
} from 'react-bootstrap'
import { initialDataload, deleteBookmark } from './actions'
import moment from 'moment'

import './style.css'

export class HomePage extends React.PureComponent {
  // eslint-disable-line react/prefer-stateless-function
  constructor (props) {
    super(props)

    this.state = {
      visible: -1,
      fuzzySearch: '',
      delete: -1,
      hide: {},
      viewMode: 0
    }
    this.filterBy = this.filterBy.bind(this)
    this.shouldDelete = this.shouldDelete.bind(this)
    this.delete = this.delete.bind(this)
    this.hide = this.hide.bind(this)
    this.changeViewMode = this.changeViewMode.bind(this)
  }

  componentDidMount () {
    if (!this.dependenciesLoaded()) {
      this.props.initialDataload()
    }
  }

  dependenciesLoaded () {
    return this.props.links && this.props.links.loaded
  }

  filterBy (v) {
    this.setState({ fuzzySearch: v.toLowerCase() })
  }

  shouldDelete (e, id) {
    e.preventDefault()
    this.setState({ delete: id })
  }

  delete (e, id) {
    e.preventDefault()
    this.props.deleteBookmark(id)
  }

  hide (e, id) {
    e.preventDefault()
    var hide = { ...this.state.hide }
    hide[id] = true
    this.setState({ hide: hide })
  }

  changeViewMode (mode) {
    this.setState({ viewMode: mode })
  }

  render () {
    if (!this.props.bookmarks.loaded) {
      return (
        <Grid>
          <Row>
            <Col>
              <div className='row-count'> loading... </div>
            </Col>
          </Row>
        </Grid>
      )
    }

    var bookmarks = this.props.bookmarks.bookmarks.filter(
      v => !this.state.hide[v.id]
    )
    if (this.state.fuzzySearch !== '') {
      var fuzzySearch = this.state.fuzzySearch.toLowerCase()
      bookmarks = bookmarks.filter(
        v =>
          fuzzyMatch(v.url.toLowerCase(), fuzzySearch) ||
          fuzzyMatch(v.title.toLowerCase(), fuzzySearch)
      )
      bookmarks.sort((a, b) => {
        const aRank =
          a.url.toLowerCase().includes(fuzzySearch) ||
          a.title.toLowerCase().includes(fuzzySearch)
        const bRank =
          b.url.toLowerCase().includes(fuzzySearch) ||
          b.title.toLowerCase().includes(fuzzySearch)
        if (aRank === bRank) {
          return 0
        }
        if (aRank) {
          return -1
        }
        return 1
      })
    }

    var listing = []
    if (this.state.viewMode === 0) {
      listing = bookmarks.map(v => (
        <Panel key={v.id} className='link-card'>
          <Grid fluid>
            <Row onClick={() => window.open(v.url, '_blank')}>
              <Col>
                <div className='link-card-title'>
                  <div>{v.title.trim() !== '' ? v.title.trim() : v.url}</div>
                  <div className='link-card-title-url'>
                    {v.host} - {moment(v.created_at).fromNow()}
                  </div>
                </div>
              </Col>
            </Row>
            <Row>
              <Col>
                <div
                  className='link-card-button-open'
                  onClick={() => window.open(v.url, '_blank')}
                >
                  <Glyphicon glyph='new-window' />
                </div>
                <div className='link-card-button-delete'>
                  {this.state.delete === v.id ? (
                    <span>
                      <Glyphicon
                        className='link-card-button-delete-ban-circle'
                        glyph='ban-circle'
                        onClick={e => this.shouldDelete(e, -1)}
                      />
                      &nbsp;
                      <Glyphicon
                        className='link-card-button-delete-trash'
                        glyph='trash'
                        onClick={e => this.delete(e, v.id)}
                      />
                    </span>
                  ) : (
                    <span>
                      <Glyphicon
                        glyph='volume-off'
                        onClick={e => this.hide(e, v.id)}
                      />
                      &nbsp;&nbsp;
                      <Glyphicon
                        glyph='remove'
                        onClick={e => this.shouldDelete(e, v.id)}
                      />
                    </span>
                  )}
                </div>
              </Col>
            </Row>
          </Grid>
        </Panel>
      ))
    } else if (this.state.viewMode === 1) {
      var stripHost = host =>
        host
          .toLowerCase()
          .trim()
          .replace('www.', '')
          .replace('http://www', '')
          .replace('https://www', '')
          .replace('http://', '')
          .replace('https://', '')
          .replace(/\/$/, '')

      bookmarks.sort((a, b) => {
        var left = stripHost(a.host)
        var right = stripHost(b.host)
        if (left === right) {
          return 0
        }
        return left < right ? -1 : 1
      })

      var repeatedURLs = {}
      for (let i in bookmarks) {
        let v = bookmarks[i]
        let strippedURL = stripHost(v.url)
        if (!repeatedURLs[strippedURL]) {
          repeatedURLs[strippedURL] = []
        }
        repeatedURLs[strippedURL].push(v)
      }

      for (let i in repeatedURLs) {
        repeatedURLs[i].sort((a, b) => {
          return a.last_status_check < b.last_status_check ? -1 : 1
        })
      }

      for (let i in repeatedURLs) {
        let v = repeatedURLs[i][0]
        listing.push(
          <Panel key={v.id} className='link-list'>
            <Grid fluid>
              {repeatedURLs[i].map((v, k) => (
                <Row key={v.id}>
                  <Col className={k > 0 ? 'link-list-secondary-container' : ''}>
                    <div
                      className={k > 0 ? 'link-list-title-secondary' : 'link-list-title'}
                      onClick={() => window.open(v.url, '_blank')} >
                      {repeatedURLs[i].length > 1 && k === 0 ? [<Label bsStyle='info'>repeated</Label>, ' '] : ''}
                      {v.title.trim() !== '' ? v.title.trim() : v.url}
                      <span className='link-list-title-host'>
                        {v.host} - {moment(v.created_at).fromNow()}
                      </span>
                    </div>
                    <div style={{ float: 'right' }}>
                      {this.state.delete === v.id ? (
                        <span>
                          <Glyphicon
                            className='link-card-button-delete-ban-circle'
                            glyph='ban-circle'
                            onClick={e => this.shouldDelete(e, -1)}
                          />
                          &nbsp;
                          <Glyphicon
                            className='link-card-button-delete-trash'
                            glyph='trash'
                            onClick={e => this.delete(e, v.id)}
                          />
                        </span>
                      ) : (
                        <span>
                          {' '}
                          <Glyphicon
                            glyph='remove'
                            onClick={e => this.shouldDelete(e, v.id)}
                          />{' '}
                        </span>
                      )}
                    </div>
                  </Col>
                </Row>
              ))}
            </Grid>
          </Panel>
        )
      }
    }
    return (
      <Grid>
        <Row>
          <Col>
            <div className='row-count'>
              <ToggleButtonGroup
                name='viewMode'
                type='radio'
                defaultValue={0}
                value={this.state.viewMode}
                onChange={this.changeViewMode}
              >
                <ToggleButton value={0}>card</ToggleButton>
                <ToggleButton value={1}>list</ToggleButton>
              </ToggleButtonGroup>
            </div>
          </Col>
        </Row>
        <Row>
          <Col>
            <div className='row-count'>
              {bookmarks.length +
                ' bookmark' +
                (bookmarks.length > 0 ? 's' : '')}
            </div>
          </Col>
        </Row>
        <Row>
          <Col>
            <FormGroup>
              <FormControl
                id='filter-box'
                type='text'
                label='Text'
                placeholder='search'
                className='filter-box'
                onChange={e => this.filterBy(e.target.value)}
              />
            </FormGroup>
          </Col>
        </Row>
        <Row>
          <Col>{this.props.bookmarks.loaded ? listing : ''}</Col>
        </Row>
      </Grid>
    )
  }
}

function s2p (state) {
  return {
    bookmarks: state.bookmarks ? state.bookmarks : { loaded: false }
  }
}

function d2p (dispatch) {
  return bindActionCreators(
    {
      initialDataload,
      deleteBookmark
    },
    dispatch
  )
}

// distilled from https://gist.github.com/mdwheele/7171422
function fuzzyMatch (haystack, needle) {
  var caret = 0
  for (var i = 0; i < needle.length; i++) {
    var c = needle[i]
    if (c === ' ') {
      continue
    }
    caret = haystack.indexOf(c, caret)
    if (caret === -1) {
      return false
    }
    caret++
  }
  return true
}

export default connect(
  s2p,
  d2p
)(HomePage)
