    function timeOutAlert(){
        timeOut = window.setTimeout(function(){
            console.log("time out..." + timeOut);

            document.getElementById("p1").innerHTML="dkdkdk" + timeOut;
            if (timeOut < 1000)
                timeOutAlert();
        },10);
    } 