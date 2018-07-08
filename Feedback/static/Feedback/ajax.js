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
            //alert(data[5]);

            let titles = data.split("~");
            let result = "";

            for (let i = 0; i < titles.length; i++)
                if (titles[i] != "" && titles[i] != " ")
                    result += `<h3>${titles[i]}</h3>`;

            $("div.hint_titles").html(result);

            if (titles.length == 1) $('.hint_titles_wrapper').addClass('hidden');
            else $('.hint_titles_wrapper').removeClass('hidden');
        }
    });
});