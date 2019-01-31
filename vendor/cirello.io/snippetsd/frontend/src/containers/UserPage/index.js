// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { Button, Col, Form, Grid, PageHeader, Row, FormControl, Panel } from 'react-bootstrap'
import { loadSnippets, saveSnippet } from './actions'
import groupBy from 'lodash/groupBy'
import moment from 'moment'

import './style.css'

class SubmitSnippetPage extends React.Component {
  constructor (props) {
    super(props)
    moment()

    this.setContent = this.setContent.bind(this)
    this.submit = this.submit.bind(this)

    this.state = {
      content: '',
      updateVisible: false
    }
  }

  componentDidMount () {
    this.props.loadSnippets()
  }

  setContent (e) {
    e.preventDefault()
    var content = e.target.value
    this.setState({ content })
  }

  submit (e) {
    e.preventDefault()
    var content = this.state.content
    this.props.saveSnippet(content)
  }

  render () {
    var snippets = groupBy(this.props.snippets, function (v) {
      return v.week_start
    })

    let groupedSnippets = {}
    for (let i in snippets) {
      groupedSnippets[i] = groupBy(snippets[i], function (v) {
        return v.user.team
      })
    }

    return (
      <Grid className='user-snippet-grid'>
        <Row>
          <Col md={12}>
            <div className='user-snippet-current-container'>
              <PageHeader> What did you do past week? </PageHeader>
              <Form onSubmit={this.submit}>
                <FormControl componentClass='textarea' className='user-snippet-content' onChange={this.setContent} />
                <div className='user-snippet-submit'><Button type='submit'>submit</Button></div>
              </Form>
            </div>
          </Col>
        </Row>
        <Row>
          <Col md={12}>
            <div className='user-snippet-past-container'>
              <PageHeader>Snippets <button>filter</button></PageHeader>
              {Object.entries(groupedSnippets).map(
              (week) => (
                <Panel key={week[0]} className='user-past-snippet'>
                  <Panel.Heading>Week starting {moment(week[0]).format('MMMM Do YYYY')}: </Panel.Heading>
                  <Panel.Body>
                    {Object.entries(week[1]).map((team) => (
                      <div key={team[0]} className='user-past-snippet'>
                        <em>{team[0]}</em>
                        {team[1].map(
                        (snippet) => (
                          <div key={snippet.user.email} className='user-snippet'>
                            <div>{snippet.user.email}: {snippet.contents || 'no snippet'} </div>
                          </div>
                        )
                      )}
                      </div>
                    ))}
                  </Panel.Body>
                </Panel>
              )
            )}
            </div>
          </Col>
        </Row>
      </Grid>
    )
  }
}

const s2p = state => ({ snippets: state.snippets.snippets })
const d2p = dispatch => bindActionCreators({
  loadSnippets,
  saveSnippet
}, dispatch)
export default connect(s2p, d2p)(SubmitSnippetPage)
