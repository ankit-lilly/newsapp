import * as smd from "./smd.min.js";

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

    let autoScroll = true;
    const tolerance = 20; // Pixel tolerance for being "at the bottom".

    const modalBox = element.closest(".prose");

    if (modalBox) {
      modalBox.addEventListener("scroll", () => {
        // Check if the user is near the bottom if not then disable autoscroll.
        if (
          modalBox.scrollTop + modalBox.clientHeight >=
          modalBox.scrollHeight - tolerance
        ) {
          autoScroll = true;
        } else {
          autoScroll = false;
        }
      });
    }

    // Only scroll if autoScroll flag is true.
    const scrollToBottom = () => {
      if (autoScroll && modalBox) {
        modalBox.scrollTop = modalBox.scrollHeight;
      }
    };

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
        scrollToBottom();
      }

      if (detail.xhr.readyState === 4) {
        element["__streamedChars"] = 0;
        element.__streamed = true;
        smd.parser_end(parser);
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
        button.classList.add("opacity-75", "cursor-not-allowed");
        break;
      }

      case "htmx:afterRequest": {
        // Restore button to original state
        button.disabled = false;

        // Restore original text and classes
        if (button.dataset.originalClasses) {
          button.className = button.dataset.originalClasses;
          delete button.dataset.originalClasses;
        }
        break;
      }

      case "htmx:timeout":
      case "htmx:error": {
        // Handle errors and timeouts
        button.disabled = false;

        if (button.dataset.originalClasses) {
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
      LIGHT: "autumn",
      DARK: "dracula",
    };

    this.handleControllerChange = this.handleControllerChange.bind(this);
    this.handleSystemThemeChange = this.handleSystemThemeChange.bind(this);
    this.init();
  }

  init() {
    this.controllers = {
      desktop: document.querySelector("#themeswitch"),
      mobile: document.querySelector("#themeswitch-mobile"),
    };
    this.setupControllerListeners();
    this.applySavedOrDefaultTheme();
    this.setupSystemThemeListener();
  }

  destroy() {
    Object.values(this.controllers).forEach((controller) => {
      if (controller) {
        controller.removeEventListener("change", this.handleControllerChange);
      }
    });
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    mediaQuery.removeEventListener("change", this.handleSystemThemeChange);
  }

  setupControllerListeners() {
    Object.values(this.controllers).forEach((controller) => {
      if (controller) {
        controller.addEventListener("change", this.handleControllerChange);
      }
    });
  }

  handleControllerChange(event) {
    console.info("Why is this not working", this.themes);
    const theme = event.target.checked ? this.themes.DARK : this.themes.LIGHT;
    console.info(`Setting theme to ${theme}`);
    this.setTheme(theme);
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

function setupCommandSearch() {
  const modal = document.getElementById("source-search-modal");
  const input = document.getElementById("source-search-input");
  if (!modal || !input || modal.dataset.commandReady) return;

  modal.dataset.commandReady = "true";
  const items = Array.from(modal.querySelectorAll(".command-source-item"));
  const emptyState = document.getElementById("source-search-empty");
  let activeIndex = -1;

  const visibleItems = () =>
    items.filter((item) => !item.classList.contains("hidden"));

  const selectItem = (index) => {
    items.forEach((item) => {
      item.querySelector("a")?.classList.remove("bg-base-200");
    });

    const visible = visibleItems();
    if (!visible.length) {
      activeIndex = -1;
      return;
    }

    activeIndex = (index + visible.length) % visible.length;
    const link = visible[activeIndex].querySelector("a");
    link?.classList.add("bg-base-200");
    link?.scrollIntoView({ block: "nearest" });
  };

  const filterItems = () => {
    const terms = input.value
      .toLowerCase()
      .trim()
      .split(/\s+/)
      .filter(Boolean);
    items.forEach((item) => {
      const value = item.dataset.commandValue || "";
      item.classList.toggle(
        "hidden",
        !terms.every((term) => value.includes(term))
      );
    });
    emptyState?.classList.toggle("hidden", visibleItems().length !== 0);
    activeIndex = -1;
  };

  const openModal = () => {
    if (!modal.open) modal.showModal();
    input.value = "";
    filterItems();
    requestAnimationFrame(() => input.focus());
  };

  document.querySelectorAll("[data-command-trigger]").forEach((trigger) => {
    trigger.addEventListener("click", openModal);
  });

  document.querySelectorAll("[data-command-shortcut]").forEach((shortcut) => {
    shortcut.textContent = /Mac|iPhone|iPad/.test(navigator.platform)
      ? "⌘ K"
      : "Ctrl K";
  });

  input.addEventListener("input", filterItems);
  modal.addEventListener("click", (event) => {
    if (event.target.closest("a")) modal.close();
  });

  document.addEventListener("keydown", (event) => {
    if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
      event.preventDefault();
      if (modal.open) {
        modal.close();
      } else {
        openModal();
      }
      return;
    }

    if (
      !modal.open ||
      !["ArrowDown", "ArrowUp", "Enter"].includes(event.key)
    ) {
      return;
    }
    event.preventDefault();

    if (event.key === "ArrowDown") selectItem(activeIndex + 1);
    if (event.key === "ArrowUp") selectItem(activeIndex - 1);
    if (event.key === "Enter") {
      const visible = visibleItems();
      const selected = visible[activeIndex < 0 ? 0 : activeIndex];
      selected?.querySelector("a")?.click();
    }
  });
}

/**
 * Focus the chat input if it exists.
 */
function focusChatInput() {
  const input = document.getElementById("chat_message_input");
  if (input) {
 //   input.focus();
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
  setupCommandSearch();
  // Create the theme manager and store it globally.
  window.themeManager = new ThemeManager();
});

/**
 * Reinitialize components after HTMX swaps. Since elements may be replaced,
 * we attach listeners again and recreate the ThemeManager (which cleans up previous listeners).
 */
document.addEventListener("htmx:afterSettle", () => {
  attachChatModalListener();
  setupCommandSearch();
  window.themeManager = new ThemeManager();
});

document.addEventListener("htmx:historyRestore", function () {
  window.themeManager = new ThemeManager();
});

document.body.addEventListener("htmx:wsAfterMessage", (event) => {
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

  document.querySelectorAll(".menu-active").forEach((heading) => {
    heading.classList.remove("menu-active");
  });

  const activeLink = document.querySelector(`.navbar a[href="${currentPath}"]`);

  if (activeLink) {
    activeLink.classList.add("menu-active");
    activeLink?.parentNode?.classList?.add("menu-active");
  }
});
