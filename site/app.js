import {MDCTopAppBar} from '@material/top-app-bar';
import {MDCSlider} from '@material/slider';
import {MDCSnackbar} from '@material/snackbar';
import {MDCIconButtonToggle} from '@material/icon-button';
import {MDCRipple} from '@material/ripple';
import {MDCSwitch} from '@material/switch';
import {MDCTextField} from '@material/textfield';

// Instantiation
const topAppBarElement = document.querySelector('.mdc-top-app-bar');
const topAppBar = new MDCTopAppBar(topAppBarElement);
const snackbar = new MDCSnackbar(document.querySelector('.mdc-snackbar'));

if (pageFunction == "login") {
    const passwordField = new MDCTextField(document.querySelector('.mdc-text-field'));
    const passwordInput = document.querySelector('.password-input-field')
    const passwordButton = document.querySelector('.password-button');
    const mDCRipple = new MDCRipple(passwordButton)
    const sendPassword = function() {
        var xhr = new XMLHttpRequest();
        xhr.open("POST", window.location.href, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                var json = JSON.parse(xhr.responseText);
                if (xhr.status !== 202) {
                    //console.log("Error: " + json.error);
                    snackbar.labelText = "Error: " + json.error;
                    snackbar.open()
                    return
                }
                document.cookie = "token="+encodeURIComponent(json.token)+";samesite=lax";
                window.location.pathname = "/"
            }
        };
        var data = JSON.stringify({"password": passwordInput.value});
        xhr.send(data);
    }
    passwordInput.addEventListener("keyup", function(event) {
        if (event.key === "Enter") {
            sendPassword()
        }
    });
    passwordButton.addEventListener('click', sendPassword)
} else if (pageFunction == "default") {
    const logoutButton = document.querySelector('.logout-button');
    logoutButton.addEventListener('click', () => {document.cookie = "token=;samesite=lax"; location.reload()})
    var statusUpdater = setInterval(updateStatus, 6000);

    const switches = document.querySelectorAll('.mdc-switch');
    switches.forEach(s => {
        const mDCSwitch = new MDCSwitch(s)
        const native = s.querySelector('.mdc-switch__native-control')
        const checked = native.attributes['checked'] !== undefined
        //console.log(mDCSwitch.checked + ' ' + checked)
        mDCSwitch.checked = checked
        native.addEventListener('change', () => {
            //console.log(`Value of ${s.attributes['controller'].value} key ${s.attributes['key'].value} changed to ${native.attributes['aria-checked'].value}`)
            updateConfigItem(s.attributes['controller'].value, "toggle", s.attributes['key'].value, native.attributes['aria-checked'].value, function(resetValue){mDCSwitch.checked = resetValue})
        })
    });

    const sliders = document.querySelectorAll('.mdc-slider')
    sliders.forEach(slider => {
        const mDCSlider = new MDCSlider(slider)
        mDCSlider.listen('MDCSlider:change', () => {
            //console.log(`Value of ${slider.attributes['controller'].value} key ${slider.attributes['key'].value} changed to ${mDCSlider.value}`)
            updateConfigItem(slider.attributes['controller'].value, "range", slider.attributes['key'].value, mDCSlider.value, function(resetValue){mDCSlider.value = resetValue; document.getElementById('slider-value-'+slider.attributes['controller'].value+'-'+slider.attributes['key'].value).textContent = mDCSlider.value})
        });
        mDCSlider.listen('MDCSlider:input', () => {
            document.getElementById('slider-value-'+slider.attributes['controller'].value+'-'+slider.attributes['key'].value).textContent = mDCSlider.value
        });
        document.getElementById('slider-value-'+slider.attributes['controller'].value+'-'+slider.attributes['key'].value).textContent = mDCSlider.value

        //This is to correctly render sliders. Ref: https://github.com/material-components/material-components-web/issues/1017
        window.addEventListener('load', () => {
            mDCSlider.layout()
        }, false);
    });

    function updateConfigItem(controller, type, key, value, cbOnError) {
        var xhr = new XMLHttpRequest();
        var url = window.location.href+"api/config";
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.setRequestHeader("Authorization", "Bearer " + getCookie("token"))
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                if (xhr.status === 401) {
                    window.location.pathname = "/login";
                } else if (xhr.status !== 200) {
                    var json = JSON.parse(xhr.responseText);
                    //console.log("Error: " + json.error);
                    snackbar.labelText = "Error: " + json.error;
                    snackbar.open();
                    cbOnError(json.origValue)
                } else {
                    updateStatus();
                }
            }
        };
        var data = JSON.stringify({"controller": controller, "type": type, "key": key, "value": value.toString()});
        xhr.send(data);
    }

    function updateStatus() {
        var xhr = new XMLHttpRequest();
        var url = window.location.href+"api/status";
        xhr.open("GET", url, true);
        //xhr.setRequestHeader("Content-Type", "application/json");
        xhr.setRequestHeader("Authorization", "Bearer " + getCookie("token"))
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                var json = JSON.parse(xhr.responseText);
                if (xhr.status !== 200) {
                    //console.log("Error: " + json.error);
                    snackbar.labelText = "Error: " + json.error;
                    snackbar.open()
                    cbOnError(json.origValue)
                } else {
                    //console.log(json)
                    //var updateDiv = document.getElementsByClassName("status-div-to-update")[0]
                    var updateDiv = document.createElement("div")
                    updateDiv.className = "mdc-layout-grid__cell mdc-layout-grid__cell--span-4 status-div-to-update"
                    var statusHeader = document.createElement("h2")
                    statusHeader.className = "mdc-typography--headline6"
                    statusHeader.innerText = "Status"
                    updateDiv.appendChild(statusHeader)
                    Object.keys(json.inputs).forEach(function(name){
                        //console.log(name + '=' + json.inputs[name]);
                        var wrapperDiv = document.createElement("div")
                        wrapperDiv.setAttribute("class", "status-wrapper")

                        var textDiv = document.createElement("div")
                        textDiv.className = "status-text"
                        var text = document.createElement("p")
                        text.innerText = name
                        textDiv.appendChild(text)
                        wrapperDiv.appendChild(textDiv)

                        var valueDiv = document.createElement("div")
                        valueDiv.className = "status-value"
                        var value = document.createElement("p")
                        value.innerText = json.inputs[name]
                        valueDiv.appendChild(value)
                        wrapperDiv.appendChild(valueDiv)

                        updateDiv.appendChild(wrapperDiv)
                    });
                    Object.keys(json.outputs).forEach(function(name){
                        //console.log(name + '=' + json.outputs[name]);
                        var wrapperDiv = document.createElement("div")
                        wrapperDiv.setAttribute("class", "status-wrapper")

                        var textDiv = document.createElement("div")
                        textDiv.className = "status-text"
                        var text = document.createElement("p")
                        text.innerText = name
                        textDiv.appendChild(text)
                        wrapperDiv.appendChild(textDiv)

                        var valueDiv = document.createElement("div")
                        valueDiv.className = "status-value"
                        var value = document.createElement("p")
                        var valueSpan = document.createElement("span")
                        valueSpan.className = "dot"
                        if(json.outputs[name]) {
                            valueSpan.className += " green-dot"
                        }
                        value.appendChild(valueSpan)
                        valueDiv.appendChild(value)
                        wrapperDiv.appendChild(valueDiv)

                        updateDiv.appendChild(wrapperDiv)
                    });
                    document.getElementsByClassName("status-div-to-update")[0].replaceWith(updateDiv)
                }
            }
        };
        xhr.send();
    }
}

function getCookie(name) {
    var cookieArr = document.cookie.split(";");
    
    for(var i = 0; i < cookieArr.length; i++) {
        var cookiePair = cookieArr[i].split("=");
        if(name == cookiePair[0].trim()) {
            var decoded = decodeURIComponent(cookiePair[1]);
            if (decoded.length == 0) {
                return null;
            }
            return decoded;
        }
    }
    return null;
}