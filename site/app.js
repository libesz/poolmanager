import {MDCTopAppBar} from '@material/top-app-bar';
import {MDCSlider} from '@material/slider';

// Instantiation
const topAppBarElement = document.querySelector('.mdc-top-app-bar');
const topAppBar = new MDCTopAppBar(topAppBarElement);

const sliders = document.querySelectorAll('.mdc-slider')
sliders.forEach(slider => {
    const mDCSlider = new MDCSlider(slider)
    mDCSlider.listen('MDCSlider:change', () => console.log(`Value changed to ${mDCSlider.value}`));

    //This is to correctly render sliders. Ref: https://github.com/material-components/material-components-web/issues/1017
    window.addEventListener('load', () => {
        mDCSlider.layout()
    }, false);
});
