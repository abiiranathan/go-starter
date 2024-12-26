function toggleTheme() {
  document.documentElement.setAttribute(
    "data-theme",
    document.documentElement.getAttribute("data-theme") === "dark"
      ? "light"
      : "dark"
  );
}

function registerServiceWorker() {
  if ("serviceWorker" in navigator) {
    navigator.serviceWorker
      .register("/sw.js")
      .then(function (registration) {
        console.log("Service Worker Registered");
      })
      .catch(function (err) {
        console.log("Service Worker Failed to Register", err);
      });
  }
}

document.addEventListener("DOMContentLoaded", function () {
  const themeToggle = document.getElementById("theme-toggle");
  if (themeToggle) {
    themeToggle.addEventListener("click", toggleTheme);
  }

  //   Uncomment the following line to enable service worker
  //   registerServiceWorker();
});
