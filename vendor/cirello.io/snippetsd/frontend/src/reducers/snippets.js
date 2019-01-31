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

export const SNIPPETS_LOADED = 'snippets/LOADED'

const initialState = {
  snippets: [
    {week_start: '2018-06-18T00:00:00Z', user: {team: 'command', email: 'kirk@domain'}, contents: 'content'},
    {week_start: '2018-06-18T00:00:00Z', user: {team: 'command', email: 'jlp@domain'}, contents: 'content'},
    {week_start: '2018-06-18T00:00:00Z', user: {team: 'science', email: 'spock@domain'}, contents: 'content'},
    {week_start: '2018-06-18T00:00:00Z', user: {team: 'medical', email: 'mccoy@domain'}, contents: ''}
  ]
}

export default (state = initialState, action) => {
  switch (action.type) {
    case SNIPPETS_LOADED:
      return {
        ...state,
        snippets: action.snippets
      }
    default:
      return state
  }
}
