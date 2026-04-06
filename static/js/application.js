$(function() {
  var $textarea = $("#contents"),
      $printarea = $('#printable_contents'),
      $loading = $("#loading"),
      $copyButton = $("#copy_button"),
      $pasteButton = $("#paste_button"),
      content = $textarea.val();

  function selectAllContents() {
    var textarea = $textarea.get(0);

    $textarea.trigger("focus");
    textarea.setSelectionRange(0, textarea.value.length);
  }

  async function copyContents() {
    var text = $textarea.val();

    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text);
      return;
    }

    // fallback for older browsers
    selectAllContents();
    document.execCommand("copy");
  }

  async function pasteContents() {
    if (navigator.clipboard && window.isSecureContext) {
      $textarea.val(await navigator.clipboard.readText());
      return;
    }

    // fallback for older browsers
    $textarea.focus();
    $textarea.setSelectionRange(0, $textarea.value.length);
    document.execCommand("paste");
  }

  $copyButton.on("click", function() {
    copyContents().catch(function() {});
  });

  $pasteButton.on("click", function() {
    pasteContents().catch(function() {});
  });

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
