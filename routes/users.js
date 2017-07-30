var express = require('express');
var router = express.Router();
var store = require('../models/database.js')
var passport = require('passport')
const fortune = require('fortune')
const bcrypt = require('bcrypt')
const saltRounds = 11
const emailRegex = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
var passport = require('passport')
var Strategy = require('passport-local').Strategy

passport.use(new Strategy(
  function(username, password, next) {
    console.log('ayy')
    //console.log(username, password)

    store.find('user', null, {match: {username: username}}).then(results => {
      if(results.payload.records.count == 0) return next(null, false)
      var user = results.payload.records[0]
      console.log(user)

      bcrypt.compare(password, user.password, function(err, res){
        console.log(res)
        if(res) return next(null, user)
        return next(null, false)
      })
    })
  })
)

passport.serializeUser(function(user, done) {
  console.log(user)
  done(null, user)
})

passport.deserializeUser(function(user, done) {
  done(null, user)
})

/* GET users listing. */
router.get('/', function(req, res, next) {
  console.log(req.user)
  res.send(' <form action="/users/login" method="post"> <input name="username"> <input name="password"> <input type="submit"> </form> ');
});

router.post('/register', async function(req, res, next){
  console.log(req.body)
  var errors = {}
  //basic data checks
  if(!req.body.password) errors.password = 'password required'
  else if(req.body.password.length < 6) errors.password = 'password too short'

  if(req.body.password != req.body.passwordConfirmation) errors.passwordConfirmation = 'passwords dont match'

  if(!req.body.username) errors.username = 'username required'
  else if(await user.usernameInUse(req.body.username)) errors.username = 'username already taken'
  
  if(!req.body.email) errors.email = 'email required'
  else if(!emailRegex.test(req.body.email)) errors.email = 'not a valid email'
  else if(await user.emailInUse(req.body.email)) errors.email = 'email already taken'

  if(!req.body.contribution || req.body.contribution < 2 || req.body.contribution > 10) errors.contribution = 'contribution value was invalid or missing - try refreshing the page'

  console.log(Object.keys(errors))
  if(Object.keys(errors).length != 0){
    res.headers
    res.status(406).json(errors)
  }

  else {
    console.log('fields look good')
    bcrypt.hash(req.body.password, saltRounds, function(err, hash){
      store.create('user', {
        username: req.body.username,
        password: hash,
        email: req.body.email,
        contribution: parseInt(req.body.contribution),
        accountConfirmed: false,
        balance: 0,
        subscription: new Date()
      }).then( results => {
        console.log(results)
      })
    })

    bcrypt.compare(req.body.password, '$2a$11$VtD7QR9KyRgw2WM81RFop.wLQoWgXePy74Jr4wHIMbzHWzVq5WEte', function(err, res){
      console.log(res)
    })
    res.redirect('/')
  }
})

router.post('/login', passport.authenticate('local', {failureRedirect: 'http://syd.jjcm.org/sociFrontend/login/#failure'}), function(req, res, next){
  console.log(req.user)
  res.redirect('/')
})

router.post('/checkUsername', function(req, res, next){
  res.json({taken: false})
})

router.post('/checkEmail', function(req, res, next){
  res.json({taken: false})
})

var user = {
  usernameInUse: function(username){
    return new Promise(resolve => {
      store.find('user', null, {match: {username: username}}).then(results => {
        resolve(results.payload.count != 0)
      })
    })
  },
  emailInUse: function(email){
    return new Promise(resolve => {
      store.find('user', null, {match: {email: email}}).then(results => {
        resolve(results.payload.count != 0)
      })
    })
  }
}


module.exports = router;
