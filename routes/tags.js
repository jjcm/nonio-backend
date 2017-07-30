var express = require('express')
var router = express.Router()
var store = require('../models/database.js')
const fortune = require('fortune')

router.get('/', function(req, res, next) {
  store.find('post', null, {match: {OC: false}}).then(results => {
    res.json(results.payload)
  })
})

router.get('/:tag', function(req, res, next){
  store.find('tag', null, {match: {name: req.params.tag}}).then(results => {
    console.log(results)
    if(results)
      res.json(results.payload)
  }).catch(err => {
    console.log(err)
  })
})

router.get('/:tag/popular', function(req, res, next){
  var range = {
    date: [ new Date(Date.now() - 497330022), new Date()]
  }
  store.find('postTag', null, {range: range}, [['post'], ['tag']]).then(results => {
    console.log(results)
    if(results)
      res.json(results.payload)
  }).catch(err => {
    console.log(err)
  })
})

router.get('/:tag/top/:range', function(req, res, next){
  getTopPostTags(req.params.range, res)
})

router.get('/:tag/top/:range/offset/:offset', function(req, res, next){
  getTopPostTags(req.params.range, res, req.params.offset)
})

function getTopPostTags(timespan, res, offset) {
  if(!offset) offset = 0
  var msInDay = 24 * 60 * 60 * 1000
  switch(timespan){
    case 'day': 
      timespan = msInDay
      break
    case 'week':
      timespan = msInDay * 7
      break
    case 'month':
      timespan = msInDay * 31
      break
    case 'year':
      timespan = msInDay * 365
      break
    default:
      timespan = msInDay * 7
      break
  }

  var range = {
    date: [ new Date(Date.now() - timespan), new Date() ]
  }

  var options = {
    range: range,
    sort: {upvotes: false},
    offset: offset,
    limit: 100
  }

  store.find('postTag', null, options, [['post'], ['tag']]).then(results => {
    console.log(results)
    if(results)
      res.json(results.payload)
  }).catch(err => {
    console.log(err)
  })
}

module.exports = router
