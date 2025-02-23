import * as smd from "streaming-markdown"

window.htmx.defineExtension("stream", {
  onEvent(name, evt) {
    if (name !== "htmx:beforeRequest") return true;

    const { detail } = evt;
    let element = detail.elt;

    if (detail.requestConfig.target) {
      element["__target"] = detail.requestConfig.target;
      element = detail.requestConfig.target;
    }

    let lastLength = 0;
    let isFirstChunk = true;

    const renderer = smd.default_renderer(element);
    const parser = smd.parser(renderer);

    detail.requestConfig.swap = "none";

    detail.xhr.addEventListener("readystatechange", () => {
      if (detail.xhr.readyState === 2 || detail.xhr.readyState === 3) {
        if (isFirstChunk) {
          const loader = element.querySelector(".my-summary-loader-icon");
          loader?.remove();
          isFirstChunk = false;
        }

        const newContent = detail.xhr.responseText.slice(lastLength);

        element["__streamedChars"] = lastLength;
        lastLength = detail.xhr.responseText.length;

        smd.parser_write(parser, newContent);
        //element.append(newContent)
      }

      if (detail.xhr.readyState === 4) {
        element["__streamedChars"] = 0;
        element.__streamed = true;
        smd.parser_end(parser)
      }
    });

    return true;
  },
});

document.body.addEventListener("htmx:beforeSwap", function (evt) {
  // If the target element has already been updated by our stream,
  // cancel the final swap.
  const target = evt.detail.target;
  if (target.__streamed) {
    evt.detail.shouldSwap = false;
    evt.preventDefault();
  }
});

window.htmx.defineExtension("button-states", {
  onEvent(name, evt) {
    // Handle only buttons/inputs with hx-post
    const button = evt.detail?.elt;
    if (
      !button ||
      !button.matches('button[hx-post], input[type="submit"][hx-post]')
    ) {
      return true;
    }

    switch (name) {
      case "htmx:beforeRequest": {
        const originalClasses = button.className;
        button.dataset.originalClasses = originalClasses;

        button.disabled = true;

        if (originalText?.trim()) {
          button.classList.add("opacity-75", "cursor-not-allowed");
        }
        break;
      }

      case "htmx:afterRequest": {
        // Restore button to original state
        button.disabled = false;

        // Restore original text and classes
        if (button.dataset.originalText) {
          button.className = button.dataset.originalClasses;
          delete button.dataset.originalText;
          delete button.dataset.originalClasses;
        }
        break;
      }

      case "htmx:timeout":
      case "htmx:error": {
        // Handle errors and timeouts
        button.disabled = false;

        if (button.dataset.originalText) {
          button.className = button.dataset.originalClasses;
          delete button.dataset.originalClasses;
        }

        // Add error indication using Tailwind classes
        button.classList.add(
          "animate-shake",
          "bg-red-50",
          "text-red-500",
          "border-red-500"
        );
        setTimeout(() => {
          button.classList.remove(
            "animate-shake",
            "bg-red-50",
            "text-red-500",
            "border-red-500"
          );
        }, 1000);
        break;
      }
    }

    return true;
  },
});

class ThemeManager {
  // Hold the current instance so we can clean up when re-initializing.
  static instance = null;

  constructor() {
    // If an instance already exists, remove its event listeners.
    if (ThemeManager.instance) {
      ThemeManager.instance.destroy();
    }
    ThemeManager.instance = this;

    this.themes = {
      LIGHT: "lemonade",
      DARK: "coffee",
    };

    // Get the theme toggle controls.
    this.controllers = {
      desktop: document.querySelector("#themeswitch"),
      mobile: document.querySelector("#themeswitch-mobile"),
    };

    // Bind the event handlers so they can be removed later.
    this.handleControllerChange = this.handleControllerChange.bind(this);
    this.handleSystemThemeChange = this.handleSystemThemeChange.bind(this);

    this.init();
  }

  /**
   * Initialize the theme manager by setting up event listeners,
   * applying the saved (or default) theme, and listening for system changes.
   */
  init() {
    this.setupControllerListeners();
    this.applySavedOrDefaultTheme();
    this.setupSystemThemeListener();
  }

  /**
   * Removes all event listeners. Called before re-initialization to prevent duplicates.
   */
  destroy() {
    Object.values(this.controllers).forEach((controller) => {
      if (controller) {
        controller.removeEventListener("change", this.handleControllerChange);
      }
    });
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    mediaQuery.removeEventListener("change", this.handleSystemThemeChange);
  }

  /**
   * Add change event listeners to each theme controller.
   */
  setupControllerListeners() {
    Object.values(this.controllers).forEach((controller) => {
      if (controller) {
        controller.addEventListener("change", this.handleControllerChange);
      }
    });
  }

  /**
   * Handles the change event on a theme controller.
   */
  handleControllerChange(event) {
    const theme = event.target.checked ? this.themes.DARK : this.themes.LIGHT;
    this.setTheme(theme);
  }

  /**
   * Applies the given theme by updating the DOM attribute, saving to localStorage,
   * and synchronizing the state of all controllers.
   */
  setTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem("theme", theme);
    Object.values(this.controllers).forEach((controller) => {
      if (controller) {
        controller.checked = theme === this.themes.DARK;
      }
    });
  }

  /**
   * Loads the saved theme from localStorage (if any) or applies the system default.
   */
  applySavedOrDefaultTheme() {
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

  /**
   * Listen for system color-scheme changes and update the theme accordingly.
   */
  setupSystemThemeListener() {
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    mediaQuery.addEventListener("change", this.handleSystemThemeChange);
  }

  /**
   * Handles system theme changes.
   */
  handleSystemThemeChange(event) {
    const newTheme = event.matches ? this.themes.DARK : this.themes.LIGHT;
    this.setTheme(newTheme);
  }
}

/**
 * Focus the chat input if it exists.
 */
function focusChatInput() {
  const input = document.getElementById("chat_message_input");
  if (input) {
    input.focus();
  }
}

/**
 * Attach the click listener to the chat modal only once.
 */
function attachChatModalListener() {
  const modal = document.getElementById("chatmodal");
  if (modal && !modal.dataset.focusListenerAttached) {
    modal.addEventListener("click", focusChatInput);
    modal.dataset.focusListenerAttached = "true";
  }
}

/**
 * Initialize the ThemeManager and attach chat modal listeners on DOMContentLoaded.
 */
document.addEventListener("DOMContentLoaded", () => {
  attachChatModalListener();
  // Create the theme manager and store it globally.
  window.themeManager = new ThemeManager();
});

/**
 * Reinitialize components after HTMX swaps. Since elements may be replaced,
 * we attach listeners again and recreate the ThemeManager (which cleans up previous listeners).
 */
document.addEventListener("htmx:afterSettle", () => {
  attachChatModalListener();
  window.themeManager = new ThemeManager();
});

/*
 * Remove the loading indicator when the assistant replies.
 */
document.body.addEventListener("htmx:wsAfterMessage", (event) => {
  const messageContent = document.getElementById("notifications");
  const input = document.getElementById("chat_message_input");

  if (input) {
    input.scrollIntoView({
      behavior: "smooth",
      block: "end",
      inline: "nearest",
    });
  }
  if (messageContent) {
    window.scrollTo({
      top: messageContent.scrollHeight,
      behavior: "smooth",
    });
  }

  const loader = document.getElementById("chatloader");
  if (
    loader &&
    event.detail &&
    event.detail.message &&
    event.detail.message.includes("assistant") &&
    !event.detail.message.includes("chatloader")
  ) {
    const parentNode = loader.parentNode;
    if (parentNode) {
      console.info("Removing loading indicator");
      const siblingOfParent = parentNode.previousSibling;
      parentNode.remove();
      if (siblingOfParent && siblingOfParent.parentNode) {
        siblingOfParent.remove();
      }
    }
  }
});

document.body.addEventListener("htmx:responseError", console.error);
/*
 * Update the active navbar link after HTMX settles, but only for category pages.
 */
document.body.addEventListener("htmx:afterSettle", () => {
  const currentPath = window.location.pathname;

  if (!currentPath.startsWith("/news")) {
    return;
  }

  document.querySelectorAll(".active").forEach((heading) => {
    heading.classList.remove("active");
  });

  const activeLink = document.querySelector(`.navbar a[href="${currentPath}"]`);

  if (activeLink) {
    activeLink.classList.add("active");
    activeLink?.parentNode?.classList?.add("active");
  }
});
