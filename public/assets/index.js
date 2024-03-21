var conn;
window.onload = function () {
  var msg = document.getElementById("msg");
  var log = document.getElementById("log");

  const idLSKey = "id"

  function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
      log.scrollTop = log.scrollHeight - log.clientHeight;
    }
  }

  document.getElementById("form").onsubmit = function () {
    if (!conn) {
      return false;
    }
    if (!msg.value) {
      return false;
    }
    message = {
      type: "chat",
      data: {
        message: msg.value
      }
    }

    conn.send(JSON.stringify(message));
    msg.value = "";
    return false;
  };


  document.getElementById("nameChange").onsubmit = function () {
    if (!conn) {
      return false;
    }
    var nameInput = document.getElementById("name")
    if (!nameInput.value) {
      return false;
    }
    let message = {
      type: "name",
      data: nameInput.value
    }

    conn.send(JSON.stringify(message));
    msg.value = "";
    return false;
  };

  if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/ws");

    conn.onopen = (evt) => {
      let initPacket = {
        type: "id",
        data: localStorage.getItem(idLSKey)
      }
      conn.send(JSON.stringify(initPacket))
    }

    conn.onclose = function (evt) {
      var item = document.createElement("div");
      item.innerHTML = "<b>Connection closed.</b>";
      appendLog(item);
    };

    conn.onmessage = function (evt) {
      var msg = JSON.parse(evt.data)

      switch (msg.type) {
        case "chat":
          handleChat(msg.data)
          break;
        case "id":
          handleId(msg.data)
          break;
        default:
          console.log(msg)
      }
    };

    const handleId = (data) => {
      localStorage.setItem(idLSKey, data)
    }

    const handleChat = (data) => {
      var item = document.createElement("div");

      var name = document.createElement("div")
      name.className = "name"
      name.innerText = data.sender;
      var message = document.createElement("div")
      message.className = "message"
      message.innerText = data.message;

      item.appendChild(name)
      item.appendChild(message)
      item.className = "messageBox"
      appendLog(item);
    }

  } else {
    var item = document.createElement("div");
    item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    appendLog(item);
  }
};