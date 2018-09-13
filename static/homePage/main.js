$(window).on('scroll', () => {
    let s = $(window).scrollTop(),
        d = $(document).height(),
        c = $(window).height();

    let scrolledArea = (s / (d - c));
    // console.log(s, d, c, scrolledArea);
    $("#main").css("opacity", scrolledArea);
    $("body").css("background-color", `rgb(${242 + scrolledArea * 13}, ${252 + scrolledArea * 3}, ${255 - scrolledArea * 13})`);
});

$('')