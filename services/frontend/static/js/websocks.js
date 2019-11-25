const websocks = (function () {

    const avatar = document.getElementById("avatar-container");

    let getSocket = function(hostName){


        if(avatar == null){
            console.error("cannot read #avatar")
            return;
        }

        const socket = new WebSocket("ws://"+hostName+"/api/websockets/stream");

        socket.onopen = function () {
            console.log("Connection opened");
            avatar.setAttribute("class", "avatar-container avatar-container-live")
        }

        socket.onmessage = function (e) {
            console.log(e.data)
            switch (e.data) {
                case "off":
                    avatar.setAttribute("class", "avatar-container avatar-container-down");
                    break;
                default:
                    avatar.setAttribute("class", "avatar-container avatar-container-live");
            }
        }

        socket.onerror = function (e) {
            console.log("error occured" + e);
        }

        socket.onclose = function () {
            console.log("Connection closed");
            avatar.setAttribute("class", "avatar-container avatar-container-down")
        }
        return socket;
    };


    return {
        GetSocket: getSocket
    }
})();