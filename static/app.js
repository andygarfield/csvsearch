"use strict";

var input = document.getElementById("searchterm");
var header;
getHeader();

input.addEventListener("keyup", function(event) {
  event.preventDefault();
  if (event.keyCode === 13) {
    document.getElementById("submit").click();
  }
});

function submit() {
  var inputval = input.value;
  var xhttp = new XMLHttpRequest();
  xhttp.open("POST", "/search/");
  xhttp.onload = function () {
    populateTable(this.responseText);
  }
  xhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  xhttp.send("search=" + encodeURIComponent(inputval));
}

function getHeader() {
  var xhttp = new XMLHttpRequest();
  xhttp.open("GET", "/getheader/");
  xhttp.onload = function() {
    header = JSON.parse(this.responseText);
  }
  xhttp.send();
}

function populateTable(res) {
  var tableDiv = document.getElementById("table");
  tableDiv.innerHTML = ""

  console.log(res.responseText);
  var r = JSON.parse(res);
  if (r === "not found") {
    var alert = document.createElement("div");
    alert.setAttribute("class", "alert alert-info");
    alert.setAttribute("role", "alert");
    alert.innerText = "Value not found";
    table.appendChild(alert);
  } else {
    renderTable(tableDiv, r);
  }
}

function renderTable(tableDiv, rows) {
  // root elements
  var tableEl = document.createElement("table");
  tableEl.setAttribute("class", "table");
  tableDiv.appendChild(tableEl);

  // header element
  var headerEl = document.createElement("thead");
  headerEl.setAttribute("class", "thead-dark");
  for (var i = 0; i < header.length; i++) {
    var fieldEl = document.createElement("th");
    fieldEl.setAttribute("scope", "col");
    fieldEl.innerText = header[i];
    headerEl.appendChild(fieldEl);
  }
  tableEl.appendChild(headerEl);

  // data elements
  for (var i = 0; i < rows.length; i++) {
    var row = document.createElement("tr");
    for (var j = 0; j < rows[i].length; j++) {
      var val = document.createElement("td");
      val.innerText = rows[i][j];
      row.appendChild(val);
    }
    tableEl.appendChild(row);
  }
}

