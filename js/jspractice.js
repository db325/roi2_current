function init() {
  var runCount = 0;
  return (check = () => {
    if (runCount > 0) {
      // {  C: runCount };
      return "already run " + runCount;
    } else {
      runCount++;
      return "initializing......";
    }
  });
}

const arr = [1, 2, 3, 4];
for (let i = 0; i < arr.length; i++) {
  setTimeout(() => {
    let y = arr;
    console.log("i am at position " + y[i]);
  }, 3000);
}

function test() {
  var arr = [1, 2, 3, 4, 5];

  for (let i = 0; i < arr.length; i++) {
    (() => {
      setTimeout((a = arr) => {
        console.log("i am at posistion " + a[i], ":>>" + arr);
      }, 3000);
    })(i);
  }
}
