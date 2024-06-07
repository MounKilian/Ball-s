let strike = document.querySelector("#strike")
let miss = document.querySelector("#miss")
let ismoving = false;

const animationOptions = {
    duration: 1000,
    iterations: 1,
    easing: 'ease-in-out'
}

const move_right = () => {
    if (ismoving) {
        return;
    }
    ismoving = true;

    ball.style.display = "none";
    card.style.backgroundColor = "green";
    card.animate(
        { backgroundPosition: '-100vw', transform: "translate3d(80%, 0, 0)" },
        animationOptions
    ).addEventListener("finish", () => animEnd(true))
}

const move_left = () => {
    if (ismoving) {
        return;
    }
    ismoving = true;
    ball.style.display = "none";
    card.style.backgroundColor = "red";
    card.animate(
        { backgroundPosition: '100vw', transform: "translate3d(-80%, 0, 0)" },
        animationOptions
    ).addEventListener("finish", () => animEnd(false))
}

strike.addEventListener("click", move_right)
miss.addEventListener("click", move_left)