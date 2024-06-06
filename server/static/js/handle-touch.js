card = document.querySelector("a.card");

// Register touch event handlers
card.addEventListener("touchstart", process_touchstart, false);
// card.addEventListener("touchcancel", process_touchcancel, false);
card.addEventListener("touchend", process_touchend, false);

let startingX;

function process_touchstart(evt) {
    startingX = evt.touches[0].clientX;
    console.log("start : " + startingX + 'px')
};

function process_touchend(evt) {
    let change = startingX - evt.changedTouches[0].clientX;
    let threshold = screen.width * 0.1;
    // let threshold = 20;
    if (Math.abs(change) < threshold) {
        console.log(`cancel because : ${Math.round(change/screen.width*100)}% < ${threshold/screen.width*100}%`)
        return
    } else if (change < 0) {
        move_right()
        console.log(`move right : ${Math.round(change/screen.width*100)}% < -${threshold/screen.width*100}%`)
    } else {
        move_left()
        console.log(`move left : ${Math.round(change/screen.width*100)}% > ${threshold/screen.width*100}%`)
    }
}