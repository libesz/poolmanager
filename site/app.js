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
        var url = window.location.href+'login';
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                if (xhr.status !== 202) {
                    snackbar.labelText = "Error: Unauthorized";
                    snackbar.open()
                    return
                }
                location.reload();
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
}

if (pageFunction == "default") {
    const logoutButton = document.querySelector('.logout-button');
    logoutButton.addEventListener('click', () => {location.href = '/logout'})

    const switches = document.querySelectorAll('.mdc-switch');
    switches.forEach(s => {
        const mDCSwitch = new MDCSwitch(s)
        const native = s.querySelector('.mdc-switch__native-control')
        native.addEventListener('change', () => {
            //console.log(`Value of ${s.attributes['controller'].value} key ${s.attributes['key'].value} changed to ${native.attributes['aria-checked'].value}`)
            updateConfigItem(s.attributes['controller'].value, "toggle", s.attributes['key'].value, native.attributes['aria-checked'].value, function(resetValue){mDCSwitch.checked = resetValue})
        })
        mDCSwitch.checked = native.attributes['aria-checked'].value
    });

    const sliders = document.querySelectorAll('.mdc-slider')
    sliders.forEach(slider => {
        const mDCSlider = new MDCSlider(slider)
        mDCSlider.listen('MDCSlider:change', () => {
            //console.log(`Value of ${slider.attributes['controller'].value} key ${slider.attributes['key'].value} changed to ${mDCSlider.value}`)
            updateConfigItem(slider.attributes['controller'].value, "range", slider.attributes['key'].value, mDCSlider.value, function(resetValue){mDCSlider.value = resetValue; slider.querySelector('.slider-value').textContent = mDCSlider.value})
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
        var url = window.location.href;
        xhr.open("POST", url, true);
        xhr.setRequestHeader("Content-Type", "application/json");
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4 && xhr.status !== 200) {
                var json = JSON.parse(xhr.responseText);
                //console.log("Error: " + json.error);
                snackbar.labelText = "Error: " + json.error;
                snackbar.open()
                cbOnError(json.origValue)
            }
        };
        var data = JSON.stringify({"controller": controller, "type": type, "key": key, "value": value.toString()});
        xhr.send(data);
    }

}