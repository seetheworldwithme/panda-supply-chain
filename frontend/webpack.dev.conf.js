const express = require('express')
const app = express()
const apiRoutes = express.Router()

let appDate = require('./src/assets/map/shanxi.json')

app.use('api', apiRoutes)
module.exports = {
    devServer:{
        before:app => {
            app.get('api/appDate', (req, res) => {
              res.json({
                errno: 0,
                data: appDate
              })
            })
          }
    }
}
