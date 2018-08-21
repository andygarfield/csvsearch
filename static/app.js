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
  var hasGeom = header.length != rows[0].length

  // root elements
  var tableEl = document.createElement("table");
  tableEl.setAttribute("class", "table");
  tableDiv.appendChild(tableEl);

  // header element
  var headerEl = document.createElement("thead");
  headerEl.setAttribute("class", "thead-dark");

  if (hasGeom) {
    headerEl.appendChild(createHeaderEl("Location"));
  }
  for (var i = 0; i < header.length; i++) {
    headerEl.appendChild(createHeaderEl(header[i]));
  }

  tableEl.appendChild(headerEl);

  // if it has geometry, don't loop through last field string
  var dataLen = rows[0].length;
  if (hasGeom) {
    dataLen--;
  }
  // data elements
  for (var i = 0; i < rows.length; i++) {
    var row = document.createElement("tr");

    if (hasGeom) {
      var link = document.createElement("a");
      link.setAttribute("href", rows[i][rows[0].length-1])
      link.setAttribute("target", "_blank");
      link.innerHTML = "View Map"

      var val = document.createElement("td");
      val.appendChild(link);

      row.appendChild(val);
    }

    for (var j = 0; j < dataLen; j++) {
      row.appendChild(createDataEl(rows[i][j]));
    }
    tableEl.appendChild(row);
  }
}

function createHeaderEl(name) {
    var fieldEl = document.createElement("th");
    fieldEl.setAttribute("scope", "col");
    fieldEl.innerText = name;

    return fieldEl
}

function createDataEl(name) {
  var val = document.createElement("td");
  val.innerText = name;
  
  return val
}
