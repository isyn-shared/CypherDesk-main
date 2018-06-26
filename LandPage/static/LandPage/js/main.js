$(document).ready(() => {
    let i = 0;
    let timer = setInterval(() => {
        i++;
        if (i == 5)
            clearInterval(timer);
        else 
            $(`.show${i}`).fadeIn(700);
    }, 200);    
});
