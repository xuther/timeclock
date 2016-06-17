var curStatus = "out"

function clockIn() {
  var time = new Date();
  console.log(time);

  var toSend = {
    "Time": time
  }

  $.ajax({
    url: "/api/clockin",
    type: "POST",
    data: JSON.stringify(toSend)
  })
}
