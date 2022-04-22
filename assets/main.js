const total = { id: 0 };
const latest = { id: -1 };
const trial26 = { id: 26 };

function getSteps(trailID) {
  return fetch("https://trana.fit/steps", {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ id: parseInt(trailID) })
  })
    .then((response) => response.text())
    .then((myText) => {
      return parseInt(myText);
    });
}

function getPronation() {
  return fetch("https://trana.fit/getOutliers", {
    method: "GET",
  })
    .then((response) => response.json())
    .then((data) => {
      return data;
    });
}

function getDistance(steps) {
  return Math.round((steps / 1294) * 100) / 100;
}

function getHealth(distance) {
  return Math.round(100 - distance);
}

function getTrials() {
  return fetch("https://trana.fit/trials", {
    method: "GET"
  })
    .then((response) => response.json())
    .then((trials) => {
      return trials;
    });
}

function displayTrials(trials) {
  trials.forEach((trial) => {
    var date = new Date(trial[1]).toLocaleDateString("en-US");
    var time = new Date(trial[1]).toLocaleTimeString("en-US");

    var div = document.createElement("div");
    div.innerHTML = `
        <div class="w-9/12 m-auto text-xl font-bold">Run ${trial[0]} (${time} ${date})</div>
        <div class="w-9/12 m-auto grid grid-flow-row-dense grid-cols-4 ...">
            <div class="col-span-2 bg-gray-400 h-40 m-5 flex flex-wrap content-center boxer">
                <div class="m-auto">
                    <div class="text-center text-2xl font-bold pb-5">Steps</div>
                    <div class="text-center text-5xl" id="latestSteps">${trial[2]}</div>
                </div>
            </div>
            <div class="col-span-2 bg-gray-400 h-40 m-5 flex flex-wrap content-center boxer">
                <div class="m-auto">
                    <div class="text-center text-2xl font-bold pb-5">Distance</div>
                    <div class="text-center text-5xl" id="latestDistance">${trial[3]} Miles</div>
                </div>
            </div>
        </div>`;
    document.getElementById("trials").appendChild(div);
  });
}

function alltrials() {
  getTrials().then(function (result) {
    count = 1;
    temp = [];
    result.forEach((trial) => {
      getSteps(trial["ID"]).then(function (steps) {
        count++;
        if (steps > 100) {
          temp.push([
            trial["ID"],
            trial["StartTime"] * 1000,
            steps,
            getDistance(steps)
          ]);
        }
        if (result.length == count) {
          displayTrials(
            temp.sort(function (a, b) {
              return a[0] - b[0];
            })
          );
        }
      });
    });
  });
}

function lifetime() {
  getSteps("0").then((steps) => {
    var distance = getDistance(steps);
    document.getElementById("totalSteps").innerHTML = steps;
    document.getElementById("totalDistance").innerHTML = distance + " Miles";
    console.log(getHealth(distance));
    document.getElementById("health").innerHTML = getHealth(distance) + " Miles";
  });
  getPronation().then((data) => {
    var diff = data[2];
    if (diff > 80) {
      document.getElementById("pronation").innerHTML = "Pronated";
    }else if (diff > 60) {
      document.getElementById("pronation").innerHTML = "Mild Pronation";
    }else if (diff > 40) {
      document.getElementById("pronation").innerHTML = "Normal";
    } else {
      document.getElementById("pronation").innerHTML = "Supinated";
    }
  });
}

alltrials();
lifetime();
