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
          //Add sibling elements to the target element
          var newText = xhr.responseText.substring(lastLength);
          element["__streamedChars"] = lastLength;
          lastLength = xhr.responseText.length;
          element.innerHTML += newText;
          element.scrollIntoView();
          var indicator = document.querySelector(".summary-stream-indicator");
          if (indicator) {
            indicator.style.display = "none";
          }
        }

        if (xhr.readyState === 4) {
          var indicator = document.querySelector(".summary-stream-indicator");
          if (indicator) {
            indicator.style.display = indicator.getAttribute(
              "data-initial-display"
            );
          }
        }
      });

      var indicator = document.querySelector(".summary-stream-indicator");
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
//Sometimes the theme switcher stops working after navigating around so this is a fix for that
//TODO: Figure out if it registers duplicate listeners and can cause memory leaks
document.addEventListener("htmx:afterSettle", () => new ThemeManager());

/*
 * This is a hack to remove the loading indicator when the assistant replies.
 */
document.body.addEventListener("htmx:wsAfterMessage", function (event) {
  const messageContent = document.getElementById("notifications");
  const input = document.getElementById("chat_message_input");
  input.scrollIntoView({
    behavior: "smooth",
    block: "end",
    inline: "nearest",
  });

  window.scrollTo({
    top: messageContent.scrollHeight,
    behavior: "smooth",
  });
  const el = document.getElementById("chatloader");
  if (
    event.detail.message.includes("assistant") &&
    !event.detail.message.includes("chatloader")
  ) {
    const parentNode = el.parentNode;

    console.info(`Removing loading indicator`);

    //we need sibiling of parent to remove the header as oob-swap removes the wrapper div
    const siblingOfParent = parentNode.previousSibling;
    parentNode.remove();
    siblingOfParent.remove();
  }
});

document.body.addEventListener("htmx:afterSettle", function () {
  const currentPath = window.location.pathname;

  if (!currentPath.startsWith("/category")) {
    return;
  }
  document.querySelectorAll(".desktop-navbar h2.active").forEach((link) => {
    link.classList.remove("active");
  });

  const activeLink = document.querySelector(
    `.desktop-navbar a[href="${currentPath}"]`
  );
  if (activeLink) {
    activeLink.parentNode.classList.add("active");
  }
});
