$(function() {
  var $textarea = $("#contents"),
      $printarea = $('#printable_contents'),
      $loading = $("#loading"),
      $copyButton = $("#copy_button"),
      $copyLinkButton = $("#copy_link_button"),
      $pasteButton = $("#paste_button"),
      content = $textarea.val();

  async function copyContents() {
    var text = $textarea.val();

    await copyText(text);
  }

  async function copyText(text) {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(text);
      return;
    }

    var textarea = document.createElement("textarea");

    textarea.value = text;
    textarea.setAttribute("readonly", "");
    textarea.style.position = "fixed";
    textarea.style.opacity = "0";
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand("copy");
    document.body.removeChild(textarea);
  }

  async function copyCurrentLink() {
    await copyText(window.location.href);
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

  $copyLinkButton.on("click", function() {
    copyCurrentLink().catch(function() {});
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
