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

            for (let element of data)
                result += `<h4><a href="./answer/?id=${element.id}">${element.title}</a></h4>`;

            $("div.hint_titles").html(result);

            if (!data.length) $('.hint_titles_wrapper').addClass('hidden');
            else $('.hint_titles_wrapper').removeClass('hidden');
        }
    });
});