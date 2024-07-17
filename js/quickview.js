//write a function to display the quick view of the drive stats
function quickView() {
  var quickView = document.getElementById("quick-view");
  quickView.style.display = "block";
  var close = document.getElementById("close");
  close.onclick = function () {
    quickView.style.display = "none";
  };
}
