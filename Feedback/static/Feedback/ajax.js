$("#message_title").keyup( function () {

    $.ajax({
        type: "GET",
        url: "found_titles/",

        data: {
            'message_title': $("#message_title").val(),
        },

        dataType: "text",
        cache: false,

        success: function (data)
        {
            var titles = data.split("~");
            var result = "";

            for(var i = 0; i < titles.length; i++)
                if(titles[i] != "" && titles[i] != " ")
                    result += "<h5>Title " + (i + 1).toString() + ": " + titles [i] + "</h5>";

            $("div.hint_titles").html(result);
        }
    });
});
