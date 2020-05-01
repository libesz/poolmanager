import {MDCTopAppBar} from '@material/top-app-bar';
import {MDCSlider} from '@material/slider';

// Instantiation
const topAppBarElement = document.querySelector('.mdc-top-app-bar');
const topAppBar = new MDCTopAppBar(topAppBarElement);

const sliders = document.querySelectorAll('.mdc-slider')
sliders.forEach(slider => {
    const mDCSlider = new MDCSlider(slider)
    mDCSlider.listen('MDCSlider:change', () => {
        console.log(`Value ${slider.attributes['controller'].value} changed to ${mDCSlider.value}`)
        updateConfigItem(slider.attributes['controller'].value, slider.attributes['key'].value, mDCSlider.value)
    });

    //This is to correctly render sliders. Ref: https://github.com/material-components/material-components-web/issues/1017
    window.addEventListener('load', () => {
        mDCSlider.layout()
    }, false);
});

function updateConfigItem(controller, key, value) {
    var xhr = new XMLHttpRequest();
    var url = window.location.href;
    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            var json = JSON.parse(xhr.responseText);
            console.log("Error: " + json.error);
        }
    };
    var data = JSON.stringify({"controller": controller, "key": key, "value": value.toString()});
    xhr.send(data);
}
