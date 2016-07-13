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
  });
}

console.log("test");
$(document).ready(function() {
  $("#clock-out").click(function() {
    clockOut();
  })

  $("#clock-in").click(function() {
    clockIn();
  })

  $("#cancel-out").click(function() {
    cancelOut();
  })
  $("#save-out").click(function() {
    confirmOut();
  })
});

function clockIn() {
  confirmIn();
}

function clockOut() {
  document.getElementById("out-popup").classList.toggle("fadeIn");
  document.getElementById("out-popup").classList.toggle("fadeOut");

  document.getElementById("overlay").classList.toggle("fadeIn");
  document.getElementById("overlay").classList.toggle("fadeOut");
}

function cancelOut() {
  document.getElementById("out-popup").classList.toggle("fadeOut");
  document.getElementById("out-popup").classList.toggle("fadeIn");

  document.getElementById("overlay").classList.toggle("fadeIn");
  document.getElementById("overlay").classList.toggle("fadeOut");
}

function confirmOut() {
  document.getElementById("out-popup").classList.toggle("fadeOut");
  document.getElementById("out-popup").classList.toggle("fadeIn");

  document.getElementById("overlay").classList.toggle("fadeIn");
  document.getElementById("overlay").classList.toggle("fadeOut");

  $("#current-status").html("out");
}

function confirmIn() {
  $("#current-status").html("in");
}
