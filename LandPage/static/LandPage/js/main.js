$(document).ready(() => {
    let i = 0;
    let timer = setInterval(() => {
        i++;
        if (i == 50)
            clearInterval(timer);
        else {
            if ($(`.show${i}`).hasClass('showNavbar'))
                $(`.show${i}`).css('display', 'flex').hide().fadeIn();
            else
                $(`.show${i}`).fadeIn(700);
        }
    }, 200);
});

//$span.css('display', 'inline').hide().fadeIn();