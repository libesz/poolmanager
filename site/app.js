import {MDCTopAppBar} from '@material/top-app-bar';
import {MDCSlider} from '@material/slider';
import {MDCSnackbar} from '@material/snackbar';
import {MDCIconButtonToggle} from '@material/icon-button';
import {MDCRipple} from '@material/ripple';


// Instantiation
const topAppBarElement = document.querySelector('.mdc-top-app-bar');
const topAppBar = new MDCTopAppBar(topAppBarElement);
const snackbar = new MDCSnackbar(document.querySelector('.mdc-snackbar'));
const iconButtonRipple = new MDCRipple(document.querySelector('.mdc-button'));

const sliders = document.querySelectorAll('.mdc-slider')
sliders.forEach(slider => {
    const mDCSlider = new MDCSlider(slider)
    mDCSlider.listen('MDCSlider:change', () => {
        console.log(`Value of ${slider.attributes['controller'].value} key ${slider.attributes['key'].value} changed to ${mDCSlider.value}`)
        updateConfigItem(slider.attributes['controller'].value, slider.attributes['key'].value, mDCSlider.value, function(resetValue){mDCSlider.value = resetValue})
    });

    //This is to correctly render sliders. Ref: https://github.com/material-components/material-components-web/issues/1017
    window.addEventListener('load', () => {
        mDCSlider.layout()
    }, false);
});

function updateConfigItem(controller, key, value, cbOnError) {
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
    var data = JSON.stringify({"controller": controller, "key": key, "value": value.toString()});
    xhr.send(data);
}
