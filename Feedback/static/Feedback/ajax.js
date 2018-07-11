$("#message_title").keyup(() => {

    $.ajax({
        type: "GET",
        url: "found_titles/",

        data: {
            'message_title': $("#message_title").val(),
        },

        dataType: "text",
        cache: false,

        success: data => {
            data = JSON.parse(data);

            let result = "";
            data.sort((a, b) => (a.dis <= b.dis ? 1 : 0));

            // console.log(data);

            for (let element of data)
                result += `<h2><a href="./answer/?id=${element.id}">${element.title}</a></h2>`;

            $("div.hint_titles").html(result);

            if (!data.length) $('.hint_titles_wrapper').addClass('hidden');
            else $('.hint_titles_wrapper').removeClass('hidden');
        }
    });
});