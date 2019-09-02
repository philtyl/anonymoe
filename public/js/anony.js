'use strict';

$(document).ready(function () {
    $('.poping.up').popup();
    $('.tabular.menu .item').tab();

    registerClipBoard();
    registerInboxFeed();
});

function registerClipBoard() {
    if ($('.clipboard-btn'.length === 0)) {
        return;
    }

    const clipboard = new ClipboardJS('.clipboard-btn');
    clipboard.on('success', function (e) {
        const caller = $('#' + e.trigger.getAttribute('id'));
        e.clearSelection();
        caller.popup('destroy');
        e.trigger.setAttribute('data-content', e.trigger.getAttribute('data-success'));
        caller.popup('show');
        e.trigger.setAttribute('data-content', e.trigger.getAttribute('data-original'))
    });

    clipboard.on('error', function (e) {
        const caller = $('#' + e.trigger.getAttribute('id'));
        caller.popup('destroy');
        e.trigger.setAttribute('data-content', e.trigger.getAttribute('data-error'));
        caller.popup('show');
        e.trigger.setAttribute('data-content', e.trigger.getAttribute('data-original'))
    });
}

function registerInboxFeed() {
    if ($('#inboxfeed').length === 0) {
        return;
    }

    window.setInterval(refreshInboxFeed, 10000);
}

function refreshInboxFeed() {
    const feed = document.getElementById("inboxfeed");
    const xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            feed.innerHTML = xhr.responseText;
            jdenticon();
        }
    };

    xhr.open("GET", window.location + "/node", true);
    try {
        xhr.send();
    } catch (err) {
        err.print();
    }
}