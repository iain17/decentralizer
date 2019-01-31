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

import config from '../../config'
import { push } from 'react-router-redux'

var cfg = config()

export function loadBookmark (url, cb) {
  fetch(cfg.http + '/loadBookmark', {
    method: 'POST',
    body: JSON.stringify({ url }),
    credentials: 'same-origin'
  })
  .then(res => res.json())
  .catch((e) => {
    console.log('cannot load bookmark information:', e)
  })
  .then((bookmark) => {
    cb(bookmark)
  })
}

export function newBookmark (bookmark) {
  return (dispatch) => {
    fetch(cfg.http + '/newBookmark', {
      method: 'POST',
      body: JSON.stringify(bookmark),
      credentials: 'same-origin'
    })
    .then(res => res.json())
    .catch((e) => {
      console.log('cannot store bookmark information:', e)
    })
    .then((bookmark) => {
      dispatch(push('/'))
    })
  }
}
