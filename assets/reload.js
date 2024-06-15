function pageReload() {
  location.reload()
}

function delayReload() {
  setTimeout(pageReload(), 200);
}

function sendLang(e) {
  const xhttp = new XMLHttpRequest();
  xhttp.onload = function () {
    // console.log("funkcija send lang event:", e.target.value)
    // document.getElementById("zadatak").innerHTML = xhttp.responseText;
    // document.getElementById("kombi_gif").style.display = "none";
    // document.getElementById("kombi_muzika").style.display = "none";
    // document.getElementById("dugme_za_zadatak1").style.color = "white";
    //document.getElementById("dugme_za_zadatak1").style.display = "none";
    
    // let selectDOM = document.getElementById("lang");
    // let arr = [...selectDOM.children]
    // console.log("delay reloaddddddddd", arr[0])
    // arr.map(c => c.value == e.target.value ? c.setAttribute('selected', 'selected') : '');
  }
  //zbog firefox-a mora async da bude false i onda se buni ali radi, ostali rade u oba slučaja bez primedbi
  let url = ""
  e.target.value == "sign_in" ? url = "" : url = e.target.value
  xhttp.open("POST", "/" + url, false);
  xhttp.send();
  delayReload()
  console.log("funkcija send lang event:", e.target.value)
}


