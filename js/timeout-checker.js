// Auth: DaJovan Bryant

/**
 * Calculates the time left until a certain time.
 * @param {number} time - The time in milliseconds.
 * @returns {number} - The remaining time in milliseconds.
 */
function timeLeft(time) {
  // ...
}

/**
 * Converts time from milliseconds to minutes.
 * @param {number} time - The time in milliseconds.
 * @returns {string} - The time in minutes as a string.
 */
function timeInMinutes(time) {
  // ...
}

/**
 * Decodes a hexadecimal string to a UTF-8 string.
 * @param {string} hexString - The hexadecimal string to decode.
 * @returns {string} - The decoded string.
 */
function decodeHexToString(hexString) {
  // ...
}

/**
 * Gets the expiration time of a cookie.
 * @param {string} cookieName - The name of the cookie.
 * @returns {string} - The expiration time as a string.
 */
function getCookieExpiration(cookieName) {
  // ...
}

/**
 * Gets the difference between a given time and the current time.
 * @param {string} time - The time to compare against in string format.
 * @returns {number} - The difference in milliseconds.
 */
function getDifference(time) {
  // ...
}

/**
 * Checks the remaining time and performs actions accordingly.
 */
function checkTime() {
  // ...
}

/**
 * Converts time from milliseconds to seconds.
 * @param {number} time - The time in milliseconds.
 * @returns {string} - The time in seconds as a string.
 */
function toSeconds(time) {
  // ...
}

/**
 * Converts time from seconds to minutes.
 * @param {number} seconds - The time in seconds.
 * @returns {string} - The time in minutes as a string.
 */
function toMinutes(seconds) {
  // ...
}

/**
 * Fetches data from the server based on center number and date/time.
 * @param {number} centernumber - The center number.
 * @param {string} dateTime - The date/time.
 */

/**
 * Gets the center number from the cookie.
 * @returns {number} - The center number.
 */
function getCenterNumber() {
  // ...
}

/**
 * Gets the value of a parameter from the URL.
 * @param {string} name - The name of the parameter.
 * @param {string} url - The URL to search for the parameter.
 * @returns {string|null} - The value of the parameter, or null if not found.
 */
function getParameterByName(name, url) {
  // ...
}

/**
 * Gets all parameters from the URL and returns them as an object.
 * @param {string} url - The URL to extract parameters from.
 * @returns {Object} - An object containing the parameters as key-value pairs.
 */
function getParamsObject(url) {
  // ...
}

/**
 * Displays a popup message to the user.
 * @param {number} timeleft - The time left until logout.
 */
function showPopUp(timeleft) {
  // ...
}
//Auth:DaJovan Bryant

function timeLeft(time) {
  //get the current time
  var currentTime = new Date().getTime();
  //get the time the user was last active
  var lastActive = new Date().getTime();
  //get the time the user will be logged out
  var timeToLogout = lastActive + time;
  //get the time left
  var remaining = timeToLogout - currentTime;
  return remaining;
}

function timeInMinutes(time) {
  return (time / 60000).toString();
}

function decodeHexToString(hexString) {
  // Convert the hexadecimal string to an array of bytes
  const byteArray = hexString
    .match(/.{1,2}/g)
    .map((byte) => parseInt(byte, 16));

  // Create a Uint8Array from the array of bytes
  const uint8Array = new Uint8Array(byteArray);

  // Decode the Uint8Array to a UTF-8 string
  const decodedString = new TextDecoder().decode(uint8Array);

  return decodedString;
}

function getCookieExpiration(cookieName) {
  const cookie = document.cookie
    .split("; ")
    .find((c) => c.startsWith(`${cookieName}=`));
  if (cookie) {
    // get all remaining values after [1]

    let wholecookie = decodeHexToString(cookie.split("=")[1]);
    let role = wholecookie.split("-")[0];
    var center = wholecookie.split("-")[1];

    let date = wholecookie.split("*");

    let cookieValue = decodeHexToString(cookie.split("=")[1]);
    //the timezone will be GMT, so ensure that the time is converted to the correct timezone
    const convertedTime = new Date(date[1]);
    console.log(convertedTime.toUTCString(), "CONVERTED TIME");

    console.log(
      cookieValue,
      "<-----cookie value",
      date,
      "date we got from the cookie"
    );
    let expTime = cookieValue.split("-")[2];

    return date[1];
  }
}

function getDifference(time) {
  // retur the difference between the time entered and the current time
  var currentTime = new Date().getTime();
  var timeToLogout = time;
  console.log(timeToLogout), " time to logout";
  //create a time object from the string
  var timeObject = new Date(timeToLogout);
  console.log(timeObject), "time object we created from the string";
  //get the time in milliseconds
  var timeInMilliseconds = timeObject.getTime();
  //get the difference
  var remaining = timeInMilliseconds - currentTime;

  return remaining;
}

let to = new TimeHandler({});

let min = to.minute();
to.From = new Date().getTime();
to.Until = new Date(getCookieExpiration("ROI")).getTime();
to.Interval = to.second().toMs(15); //check every second until the user is logged out
to.Stepper = "Minutes";
var struct = to.makeTimeStruct();

function checkTime() {
  let currentTime = new Date().getTime();
  // check to see if current time is greater than or equal to the time the user should be logged out

  let value = to.Until - currentTime;

  let minutes = to.timeInMinutes(to.timeInSeconds(value));

  console.log(`You will be logged out in ${minutes} minutes`);

  //trim the minutes to remove the decimal points
  minutes = Math.trunc(minutes);

  //if the time is less than or equal to 0.5 minutes, show a popup
  if (minutes <= 1) {
    showPopUp(minutes);
    console.log("showing popup");
  }

  if (minutes <= 0) {
    //log the user out
    clearInterval(checkTime);
    console.log("logging out user");
    window.location.href = "/logout";
    return;
  }
}

setInterval(checkTime, to.Interval);

function toSeconds(time) {
  let t = time / 1000;
  return t.toString() + " seconds";
}

function toMinutes(seconds) {
  let t = seconds / 60000;
  return t.toString() + " minutes";
}

//console.log(toMinutes(getDifference(getCookieExpiration("ROI"))));

const getDrives = (centernumber, dateTime) => {
  fetch(`/get-drives?center=${centernumber}&dateTime=${dateTime}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  })
    .then((response) => response.json())
    .then((data) => {
      if (data == null) {
        data = [];
        console.log("no data found");
      }
      //console.log(data);
      data.forEach((drive) => {
        localStorage.setItem(
          JSON.stringify(drive.dateTime),
          JSON.stringify(drive)
        );
        console.log("this is the drivedata", drive);
      });
    })
    .catch((error) => console.error("Error:", error));
};

//create a function that gets the centernumber from the cookie

function getCenterNumber() {
  const cookie = document.cookie.split("; ").find((c) => c.startsWith("ROI="));
  if (cookie) {
    // get all remaining values after [1]
    let wholecookie = decodeHexToString(cookie.split("=")[1]);
    let role = wholecookie.split("*")[0];
    var center = wholecookie.split("-")[1].split("*")[0];
    console.log("this is what we got from the cookie", center);

    return Number(center);
  }
}

//check for parameters in the url
function getParameterByName(name, url) {
  if (!url) url = window.location.href;
  console.log(url);
  name = name.replace(/[\[\]]/g, "\\$&");
  var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
    results = regex.exec(url);
  console.log(results);
  if (!results) return null;
  if (!results[2]) return "";
  return decodeURIComponent(results[2].replace(/\+/g, " "));
}

console.log("UNTIL VALUE", to.Until);
console.log("FOUND DRIVE", getParameterByName("Drive", window.location.href));

//function that gets all parameters from the url and returns them as an object
function getParams(url) {
  if (!url) url = window.location.href;
  var urlParams = new URLSearchParams(url);
  var params = {};
  for (let [key, value] of urlParams) {
    params[key] = value;
  }
  return params;
}

//create a popup that will display a message to the user saying that they will be logged out in 1 minute the it closes itself
function showPopUp(timeleft) {
  const popup = document.createElement("div");
  popup.style.width = "300px";
  popup.style.height = "150px";
  popup.style.position = "fixed";
  popup.style.top = "50%";
  popup.style.left = "50%";
  popup.style.transform = "translate(-50%, -50%)";
  popup.style.backgroundColor = "white";
  popup.style.border = "1px solid #ccc";
  popup.style.borderRadius = "5px";
  popup.style.boxShadow = "0 2px 4px rgba(0, 0, 0, 0.1)";
  popup.style.padding = "20px";
  popup.style.textAlign = "center";

  const message = document.createElement("p");
  message.style.fontSize = "18px";
  message.style.marginBottom = "10px";
  //if timeleft is less than 1 minute, display a message saying seconds instead of minutes
  if (timeleft < 1) {
    message.innerText = `You will be logged out in ${timeleft} seconds`;
  } else {
    message.innerText = `You will be logged out in ${timeleft} minutes`;
  }

  const closeButton = document.createElement("button");
  closeButton.style.backgroundColor = "#ccc";
  closeButton.style.border = "none";
  closeButton.style.color = "white";
  closeButton.style.padding = "8px 16px";
  closeButton.style.borderRadius = "4px";
  closeButton.style.cursor = "pointer";
  closeButton.innerText = "Close";

  closeButton.addEventListener("click", function () {
    document.body.removeChild(popup);
  });

  popup.appendChild(message);
  popup.appendChild(closeButton);

  document.body.appendChild(popup);

  setTimeout(function () {
    document.body.removeChild(popup);
  }, 3000);
}

console.log("params object", getParamsObject(window.location.href));

//write a function that uses the getDrives function then creats a view showing the drive's stats
//the view should have a table that shows the drive's stats

function createDriveView(data) {
  // Create a div element for the view
  let view = document.createElement("div");
  view.style.width = "250px";
  view.style.height = "250px";
  view.style.border = "1px solid #ccc";
  view.style.borderRadius = "5px";
  view.style.boxShadow = "0 2px 4px rgba(0, 0, 0, 0.1)";
  view.style.overflow = "auto";
  view.style.display = "flex";
  view.style.flexDirection = "column";
  view.style.alignItems = "center";
  view.style.justifyContent = "center";
  view.style.padding = "8px";

  // Create a header element for the datetime
  let header = document.createElement("h4");
  header.style.padding = "8px";
  header.style.margin = "0";
  header.style.fontSize = "x-small"; // Set fontsize to x-small
  //set font-family to ubuntu
  header.style.fontFamily = "Ubuntu";
  header.style.backgroundColor = "#f0f0f0";
  header.style.borderBottom = "1px solid #ccc";
  header.innerText = data.dateTime;

  // Create a list element
  let list = document.createElement("ul");
  list.style.listStyleType = "none";
  list.style.padding = "0";
  list.style.margin = "0";

  // Loop through the data object and create list items for each key-value pair
  for (let key in data) {
    if (key !== "datetime") {
      let listItem = document.createElement("li");
      listItem.style.padding = "8px";
      listItem.style.borderBottom = "1px solid #ccc";
      listItem.style.display = "flex";
      listItem.style.alignItems = "center";
      listItem.style.justifyContent = "space-between";

      let keyElement = document.createElement("span");
      keyElement.style.fontWeight = "bold";
      keyElement.classList.add(
        "badge",
        "bg-dark",
        "rounded-pill",
        "text-light"
      );
      keyElement.innerText =
        key.replace(/([A-Z])/g, " $1").toLowerCase() + ": ";
      keyElement.style.paddingRight = "8px";

      let valueElement = document.createElement("span");
      valueElement.style.color = "mediumgray"; // Set value color to medium gray
      valueElement.style.fontWeight = "bold";
      valueElement.style.fontSize = "x-small"; // Set fontsize to x-small
      valueElement.style.textAlign = "center";
      valueElement.style.padding = "4px";
      valueElement.innerText = data[key];

      listItem.appendChild(keyElement);
      listItem.appendChild(valueElement);
      list.appendChild(listItem);
    }
  }

  // Append the header and list to the view
  view.appendChild(header);
  view.appendChild(list);

  // Append the view to the document body
  document.body.appendChild(view);
}

// Example usage:
const data = {
  totalDonors: 10,
  totalUnits: 100,
  totalExpenses: 500.5,
  totalWages: 1000.75,
  totalSupplies: 200.25,
  totalTravel: 300.75,
  totalCost: 1500.5,
  costIfPurchasedRBC: 1000.25,
  costIfPurchasedFFP: 800.5,
  totalCostAvoidance: 200.75,
  totalUnitsPerHr: 50.25,
  totalUnitsPerPerson: 25.5,
  totalCostPerUnit: 15.75,
  APH: 5.25,
  ACH: 2.5,
  driveScore: "Pass",
  dateTime: "2022-01-01 12:00:00",
};

//createDriveView(data);
