let selectedWindow = "#status";

function showWindow(windowID) {
    if (selectedWindow != windowID) {
        if (DEBUG) console.log(windowID);


        let prevSelected = selectedWindow;
        selectedWindow = windowID;

        // This piece of code lights up navbar elements
        // In order to get <a> from window id you need to add "A" to the end
        // i.e. #status is a window, while #statusA is a corresponding <a>
        $(prevSelected + "A").removeClass('active');
        $(selectedWindow + "A").addClass('active');

        $(prevSelected).removeClass('slideIn').addClass('slideAway');
        // We need to place that "out" only when we do "slideAway".
        // If we don't remove that event, it will trigger in all different cases
        $(prevSelected).on('animationend', () =>
            $(prevSelected).addClass('out').removeClass('slideAway').off('animationend'));

        $(selectedWindow).removeClass('out').addClass('slideIn');
    }
}