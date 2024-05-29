window.htmx.defineExtension("stream", {
  onEvent: function (name, evt) {
    if (name === "htmx:beforeRequest") {
      var element = evt.detail.elt;
      if (evt.detail.requestConfig.target) {
        element["__target"] = evt.detail.requestConfig.target;
        element = evt.detail.requestConfig.target;
        element.innerHTML = "";
      }

      var xhr = evt.detail.xhr;

      var lastLength = 0;
      xhr.addEventListener("readystatechange", function () {
        if (xhr.readyState === 2 || xhr.readyState === 3) {
          var newText = xhr.responseText.substring(lastLength);
          element["__streamedChars"] = lastLength;
          lastLength = xhr.responseText.length;
          element.innerHTML += newText;
          var indicator = document.querySelector(".htmx-indicator");
          if (indicator) {
            indicator.style.display = "none";
          }
        }
      });
    }
    return true;
  },
  transformResponse: function (text, _xhr, elt) {
    var lastLength = elt["__streamedChars"];
    var target = elt["__target"];
    if (target) {
      lastLength = target["__streamedChars"];
    }

    if (lastLength) {
      var newText = text.substring(lastLength);
      return newText;
    }

    return text;
  },
});


document.addEventListener("DOMContentLoaded", function () {
  const themeController = document.querySelector("#themeswitch");

  function setTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme);
    themeController.checked = theme === "dracula";
  }
  const savedTheme = localStorage.getItem("theme");
  if (savedTheme) {
    document.documentElement.setAttribute("data-theme", savedTheme);
    themeController.checked = savedTheme === "dracula";
  } else {
    const systemPrefersDark = window.matchMedia(
      "(prefers-color-scheme: dark)"
    ).matches;
    const defaultTheme = systemPrefersDark ? "dracula" : "nord";
    console.log(
      "No saved theme found, setting theme based on system preference:",
      defaultTheme
    );
    setTheme(defaultTheme);
  }

  themeController.addEventListener("change", function () {
    const theme = this.checked ? "dracula" : "nord";
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme); // Save the theme in local storage
  });

  window
    .matchMedia("(prefers-color-scheme: dark)")
    .addEventListener("change", (e) => {
      const newColorScheme = e.matches ? "dracula" : "nord";
      console.log("System theme changed:", newColorScheme);
      setTheme(newColorScheme);
    });
});
