#! /usr/bin/env node

require("dotenv").config({ silent: true });
const path = require("path");
const express = require("express");
const publicPath = path.join(__dirname, "build");
//const app = require("https-localhost")()
// app is an express app, do what you usually do with express
const app = express();
const bodyParser = require("body-parser");
app.use(bodyParser.json({ limit: "10mb", extended: true }));
app.use(bodyParser.urlencoded({ limit: "10mb", extended: true }));
// app.use(bodyParser.raw());
app.use(express.static(publicPath));

//Host react application on root url
app.get("*", (req, res) => {
  res.sendFile(path.join(publicPath, "index.html"));
});

// app.get("/ping", (req, res) => {
//   res.json({
//     Test: "Privacy Policy"
//   });
// });

//Run application
const port = process.env.PORT || 3500;

app.listen(port, () => {
  console.log("Server running on port: %d", port);
});