$(function() {
  var $textarea = $("#contents"),
      $printarea = $('#printable_contents'),
      $toast = $("#toast"),
      $copyButton = $("#copy_button"),
      $copyLinkButton = $("#copy_link_button"),
      content = $textarea.val(),
      toastTimeout,
      toastHideTimeout;

  function showToast(message) {
    clearTimeout(toastTimeout);
    clearTimeout(toastHideTimeout);

    $toast
      .text(message)
      .removeClass("hidden opacity-0")
      .addClass("opacity-100");

    toastTimeout = setTimeout(function() {
      $toast
        .removeClass("opacity-100")
        .addClass("opacity-0");

      toastHideTimeout = setTimeout(function() {
        $toast.addClass("hidden");
      }, 200);
    }, 1600);
  }

  async function copyContents() {
    var text = $textarea.val();

    await copyText(text);
    showToast("Text copied");
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
    showToast("Link copied");
  }

  $copyButton.on("click", function() {
    copyContents().catch(function() {});
  });

  $copyLinkButton.on("click", function() {
    copyCurrentLink().catch(function() {});
  });

  $printarea.text(content);

  setInterval(function() {
    if (content !== $textarea.val()) {
      content = $textarea.val();

      $.post('', { t: content });

      $printarea.text(content);
    }
  }, 1000);
});
