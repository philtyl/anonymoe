'use strict';

$(document).ready(function () {
    $('.poping.up').popup();
    $('.tabular.menu .item').tab();

    // Clipboard JS
    var clipboard = new ClipboardJS('.clipboard');
    clipboard.on('success', function (e) {
        e.clearSelection();
        $('#' + e.trigger.getAttribute('id')).popup('destroy');
        e.trigger.setAttribute('data-content', e.trigger.getAttribute('data-success'));
        $('#' + e.trigger.getAttribute('id')).popup('show');
        e.trigger.setAttribute('data-content', e.trigger.getAttribute('data-original'))
    });

    clipboard.on('error', function (e) {
        $('#' + e.trigger.getAttribute('id')).popup('destroy');
        e.trigger.setAttribute('data-content', e.trigger.getAttribute('data-error'));
        $('#' + e.trigger.getAttribute('id')).popup('show');
        e.trigger.setAttribute('data-content', e.trigger.getAttribute('data-original'))
    });
});