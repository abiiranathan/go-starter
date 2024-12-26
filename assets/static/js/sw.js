// The version of the service worker must be updated to force the browser to update the cache
// The cache version is defined in the cacheVersion variable/constant.
const cacheVersion = "1";
const STATIC_CACHE_NAME = `static-cache-${cacheVersion}`;

const to_cache = [
  "/offline",
  "/static/css/style.css",
  "/static/js/script.js",

  // Images
  "/static/img/apple-touch-icon.png",
  "/static/img/favicon-192x192.png",
  "/static/img/favicon-512x512.png",
];

const staticAssets = new Set(to_cache);

self.addEventListener("install", (event) => {
  self.skipWaiting();
  event.waitUntil(
    caches.open(STATIC_CACHE_NAME).then((cache) => cache.addAll(to_cache))
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then(async (keys) => {
      // delete old caches
      for (const key of keys) {
        if (key !== STATIC_CACHE_NAME) {
          await caches.delete(key);
        }
      }
      await self.clients.claim();
    })
  );
});

self.addEventListener("fetch", (event) => {
  const request = event.request;
  if (request.method !== "GET" || request.headers.has("range")) return;

  const skip = request.cache === "no-store" || request.cache == "no-cache";
  if (skip) {
    return;
  }

  event.respondWith(
    (async () => {
      const cachedAsset = await caches.match(request);
      if (cachedAsset) {
        return cachedAsset;
      }

      const networkFetchPromise = fetch(request).then(async (response) => {
        if (response.ok) {
          return response;
        }

        // fetch failed, we are offline
        return caches.match("/offline");
      });

      return networkFetchPromise;
    })()
  );
});
