//const Ms = 60 / 60000;
const Sec = 1000;
const Min = Sec * 60;
const Hr = Min * 60;
const Day = Hr * 24;
const Week = Day * 7;
const Month = 2629744 || Day * 30 || Day * 31; //Seconds
const Yr = 31556926; //seconds
const MinString = "Minute";
const SecString = "Second";
const MsString = "Millisecond";
const HourString = "Hour";
const DayString = "Day";
const timeMap = new Map();

class TimeHandler {
  Until;
  From;
  Interval;
  Stepper;
  constructor(timeInfoObject) {
    if (timeInfoObject.defaultTimeIntervalUnits) {
      this.Until = timeInfoObject.until;
      this.From = timeInfoObject.from;
      this.Interval = timeInfoObject.defaultTimeInterval;
      this.Stepper = timeInfoObject.defaultTimeIntervalUnits;
      this.Name = timeInfoObject.name;
    }
    if (!timeInfoObject.defaultTimeIntervalUnits) {
      console.log("this is what we have....", timeInfoObject);
    }
  }

  ms() {
    let t = {
      value: 60 / 60000,
      name: "Milliseconds",
      unit: "Default",
    };
    return t;
  }

  minute(v) {
    if (v) {
      let t = {
        value: Min * v,
        name: "Minute(s)",
        unit: "Seconds",
        numericVal: v,
        //Takes minutes
        toMs: (mins) => {
          if (mins == null) {
            console.log("null...v =", v);
            return v * Min * 60000;
          } else {
            return mins * 60000;
          }
        },
      };
      return t;
    } else {
      let t = {
        value: 60,
        name: "Minute(s)",
        numericVal: 1,
        unit: "Seconds",
        toMs: (mins) => {
          if (mins == null) {
            return mins * 60000;
          } else {
            return mins * 60000;
          }
        },
      };
      return t;
    }
  }
  //second function that takes a value and returns a time object of type second
  second(v) {
    console.log(v);
    if (v) {
      let t = {
        value: v,
        name: "Seconds",
        unit: "Seconds",
        numericVal: v,
        toMs: (v) => {
          return v * 1000;
        },
      };
      return t;
    } else {
      let t = {
        value: Sec,
        name: "Seconds",
        unit: "Miliseconds",
        numericVal: 1,
        toMs: (n) => {
          return n * 1000;
        },
      };
      return t;
    }
  }

  hour(v) {
    if (v !== null) {
      let h = {
        value: 3600 * v,
        name: "Hour(s)",
        units: "Seconds",
        numericVal: v,
        toMs: () => {
          return (h.value / 60) * 60000;
        },
      };
      return h;
    } else {
      let h = {
        value: Hr,
        name: "Hour(s)",
        units: "Seconds",
        numericVal: 1,
        toMs: () => {
          return ((h.value * h.numericVal) / 60) * 60000;
        },
      };
      return h;
    }
  }

  day(v) {
    if (v) {
      let d = {
        value: Day,
        name: "Day(s)",
        units: "Seconds",
        numericVal: v,
        toMs: () => {
          return d.value * d.numericVal * 1000;
        },
      };
      return d;
    } else {
      let d = {
        value: Day,
        name: "Day(s)",
        unit: "Seconds",
        numericVal: 1,
        toMs: () => {
          return d.value * d.numericVal * 1000;
        },
      };
      return d;
    }
  }
  week() {
    return {
      value: Week,
      name: "Week(s)",
      unit: "Seconds",
      numericVal: 1,
    };
  }

  //Takes ms and returns seconds
  timeInSeconds(msx) {
    let v = msx;
    return v / 1000;
  }
  //Takes secs and returns minutes
  timeInMinutes(seconds) {
    seconds *= Sec;
    //  console.log(Min)
    return seconds / Min;
  }
  //takes minutes and returns hours
  timeInHours(minutes) {
    minutes *= Min;
    // console.log(Hr)
    return minutes / Hr;
  }
  //Takes hrs and returns days
  timeInDays(hours) {
    hours *= Hr;
    return hours / Day;
  }
  //Takes days and returns weeks
  timeInWeeks(days) {
    days *= Day;
    return days / Week;
  }
  //Takes weeks and returns months
  timeInMonths(weeks) {
    weeks *= Week;
    return weeks / Month;
  }
  //Takes months and returns years
  timeInYears(months) {
    months *= Month;
    return months / Yr;
  }
  //Returns the difference in milliseconds between the start and end time
  range() {
    return this.Until - this.From;
  }

  //**Returns a TimeStruct Object that will be used for time animations */
  makeTimeStruct() {
    //If we get back null or undefined values, we default the time stepper to 1 hour
    if (this.Stepper == "" || this.Stepper == undefined) {
      this.Stepper = "Hours";
    }
    if (this.Interval == "" || this.Interval == undefined) {
      this.Interval = 1;
    }
    //If  these fields are undefined, it's likely coming from IDPGIS

    if (this.From == undefined || this.Until == undefined) {
      //possibly idp_validtime
      //Set the range
      // this.From=this.timeInfoObject.timeExtent[0]
      // this.Until=this.timeInfoObject.timeExtent[1]
      //Get the difference
      let tmpVar = this.timeInHours(this.timeInMinutes(this.timeInSeconds()));
      // console.log(tmpVar)
      //Create new date range
      let newUntil = new Date().getTime() + 3 * Hr * 1000;
      this.Until = newUntil;
      this.From = new Date().getTime() - 4 * Hr * 1000;
      let difff = new Date().getTime() - this.From;
      // console.log("diff in hours...",difff/Hr)
    }

    if (this.Stepper.indexOf("Minutes") > -1) {
      //  console.log("yes minutes")
      let m = this.minute(this.Interval);

      return {
        from: this.From,
        until: this.Until,
        intv: this.Interval,
        stepper: this.Stepper,
        timeElement: m,
        minutes: this.timeInMinutes(this.timeInSeconds(this.range())),
        hours: this.timeInHours(
          this.timeInMinutes(this.timeInSeconds(this.range()))
        ),
      };
    } else if (this.Stepper.indexOf("Hours") > -1) {
      //    console.log("Creating Hour Stepper")
      let h = this.hour(this.Interval);
      //   console.log(range() )
      return {
        from: this.From,
        until: this.Until,
        intv: this.Interval,
        stepper: this.Stepper,
        timeElement: h,
        coverageMinutes: this.timeInMinutes(this.timeInSeconds(this.range())),
        coverageHours: this.timeInHours(
          this.timeInMinutes(this.timeInSeconds(this.range()))
        ),
      };
    } else if (this.Stepper.indexOf("Days") > -1) {
      // console.log("yes days")
      let d = this.day(this.Interval);
      // console.log(d)

      return {
        from: this.From,
        until: this.Until,
        intv: this.Interval,
        stepper: this.Stepper,
        timeElement: d,
        coverageMinutes: this.timeInMinutes(this.timeInSeconds(this.range())),
        coverageHours: this.timeInHours(
          this.timeInMinutes(this.timeInSeconds(this.range()))
        ),
      };
    }
  }
}

function startAnimation(timeStruct, src, sTime, xTimeFromNow) {
  //Object to store time and image data
  let picElem = {};
  //Variable to update stepper
  var update = timeStruct.intv;
  var canvas = document.getElementsByTagName("canvas")[0];

  if (timeStruct.timeElement.name == "Hour(s)") {
    sTime.setHours(sTime.getHours() + update);
    console.log(sTime, new Date(timeStruct.until));

    // console.log(timeMap)
  } else if (timeStruct.timeElement.name == "Minute(s)") {
    sTime.setMinutes(sTime.getMinutes() + update);
    console.log(sTime, update);
  }
  //Make sure the request time doesn't exceed the range
  if (sTime.getTime() > timeStruct.until) {
    console.log(
      "Time Map is ready ..about to start from the beginning......",
      timeMap
    );
    begin = new Date(timeStruct.from);
    //Once time reaches end, stop
    stop();
    //Then create GIF
    let y = 0;
    for (let timeStamp of timeMap.keys()) {
      // alert(timeStamp); // cucumber, tomatoes, onion
      picElemArray[y].time = timeStamp;
      y++;
    }
    gifshot.createGIF(
      {
        images: getPics(picElemArray),
        gifWidth: 500,
        gifHeight: 500,
        fontWeight: "bold",
        frameDuration: 1.9, //10=1sec,
        text:
          "START: " +
          new Date(timeStruct.from).toISOString() +
          "  END:" +
          new Date(picElemArray[picElemArray.length - 1].time).toISOString(),
        fontSize: "15px",
        fontColor: "#0d0c35 ",
      },
      function (obj) {
        if (!obj.error) {
          var image = obj.image,
            animatedImage = document.createElement("img");
          animatedImage.src = image;
          animatedImage.width = 500;
          animatedImage.height = 500;
          document.body.appendChild(animatedImage);
        }
      }
    );
  }
  if (sTime.getTime < threeHoursAgo().getTime()) {
    begin = threeHoursAgo();
  }

  src.updateParams({ TIME: sTime.getTime() });

  canvas.toBlob((blob) => {
    var newImg = document.createElement("img");
    url = URL.createObjectURL(blob);
    newImg.onload = () => {
      URL.revokeObjectURL(url);
    };
    newImg.src = canvas.toDataURL("image/png", 1.0);
    picElem.time = sTime;
    picElem.img = newImg.src;
    picElemArray.push(picElem);
    timeMap.set(sTime.getTime(), picElem.img);
  });

  console.log();
}

async function getTimeInfoObject2(u) {
  let initialTimeObject = await fetch(u + fmtString);
  let resp = await initialTimeObject.json();
  var TI = {};

  if (resp.timeInfo) {
    //create object,then assign timeInfo object to it
    console.log(resp);
    Object.assign(TI, resp.timeInfo);
  }
  //if there's a name, add it to the timeinfoObject
  if (resp.name || resp.mapName) {
    TI.name = resp.name || resp.mapName;
    //
  }
  //now lets check for most recent time range
  let timeUpdates = await fetch(u + "?returnUpdates=true&f=pjson");
  let returnTimeResponse = await timeUpdates.json();
  if (returnTimeResponse.timeExtent) {
    console.log(returnTimeResponse);
    // console.log(returnTimeResponse.timeExtent)
    TI.from = returnTimeResponse.timeExtent[0];
    TI.until = returnTimeResponse.timeExtent[1];
    delete TI.timeExtent;
    return TI;
  } else {
    //console.log(returnTimeResponse)

    return TI;
  }
  //  return TI
}

//TODO:refactor so that the function only needs the layer....so create a function to add the timestruct to the layer
function stuff2(timestructure, layer) {
  stop();
  console.log(timeMap, picElemArray);
  let src = layer.getSource();
  console.log(src, "source from stuff2");
  var stopHere = xTimeFromNow(timestructure);
  animationId = setInterval(
    () => startAnimation(timestructure, src, begin, stopHere),
    1000 / 0.6
  );
  src.refresh();
}
