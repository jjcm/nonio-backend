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
  store.find('postTag', null, {match: {tag: req.params.tag}}).then(results => {
    console.log(results)
    if(results)
      res.json(results.payload)
  }).catch(err => {
    console.log(err)
  })
})

module.exports = router
