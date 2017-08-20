$(function() {
    var $textarea = $("#contents"),
        $printarea = $('#printable_contents'),
        $loading = $("#loading"),
        content = $textarea.val();

    $textarea.tabby();
    $textarea.focus();

    $printarea.text(content);

    setInterval(function() {
      if (content !== $textarea.val()) {
        content = $textarea.val();

        $loading.show();
        $.post('', { t: content }).always(function() {
          $loading.hide();
        });

        $printarea.text(content);
      }
    }, 1000);
});
