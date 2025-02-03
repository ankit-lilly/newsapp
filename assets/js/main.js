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

        if (xhr.readyState === 4) {
          var indicator = document.querySelector(".htmx-indicator");
          if (indicator) {
            indicator.style.display = indicator.getAttribute(
              "data-initial-display"
            );
          }
        }
      });

      var indicator = document.querySelector(".htmx-indicator");
      if (indicator) {
        indicator.setAttribute("data-initial-display", indicator.style.display);
        indicator.style.display = "flex";
      }
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

class ThemeManager {
  constructor() {
    this.themes = {
      LIGHT: "lemonade",
      DARK: "dracula",
    };
    this.controllers = {
      desktop: document.querySelector("#themeswitch"),
      mobile: document.querySelector("#themeswitch-mobile"),
    };
    this.init();
  }

  init() {
    this.setupEventListeners();
    this.loadAndApplyTheme();
    this.setupSystemThemeListener();
  }

  setupEventListeners() {
    Object.values(this.controllers).forEach((controller) => {
      if (controller) {
        controller.addEventListener("change", () => {
          const theme = controller.checked
            ? this.themes.DARK
            : this.themes.LIGHT;
          this.setTheme(theme);
        });
      }
    });
  }

  setTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme);
    Object.values(this.controllers).forEach((controller) => {
      if (controller) {
        controller.checked = theme === this.themes.DARK;
      }
    });
  }

  loadAndApplyTheme() {
    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      this.setTheme(savedTheme);
    } else {
      const systemPrefersDark = window.matchMedia(
        "(prefers-color-scheme: dark)"
      ).matches;
      const defaultTheme = systemPrefersDark
        ? this.themes.DARK
        : this.themes.LIGHT;
      this.setTheme(defaultTheme);
    }
  }

  setupSystemThemeListener() {
    window
      .matchMedia("(prefers-color-scheme: dark)")
      .addEventListener("change", (e) => {
        const newTheme = e.matches ? this.themes.DARK : this.themes.LIGHT;
        this.setTheme(newTheme);
      });
  }
}

document.addEventListener("DOMContentLoaded", () => new ThemeManager());
document.addEventListener("htmx:afterSettle", () => new ThemeManager());

// Listen for an htmx event that indicates content has been swapped and settled
document.body.addEventListener("htmx:afterSettle", function () {
  // Get the current path
  const currentPath = window.location.pathname;

  if (!currentPath.startsWith("/category")) {
    return;
  }
  // Remove 'active' class from all nav links
  document.querySelectorAll(".desktop-navbar h2.active").forEach((link) => {
    link.classList.remove("active");
  });

  // Find the nav link that matches the current path
  const activeLink = document.querySelector(
    `.desktop-navbar a[href="${currentPath}"]`
  );
  if (activeLink) {
    activeLink.parentNode.classList.add("active");
  }
});
