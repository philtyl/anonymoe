self.addEventListener('install', function(e) {
  e.waitUntil(
    caches.open('anony-store').then(function(cache) {
      return cache.addAll([
        '/',
        '/manifest.json',
        '/robots.txt',
        '/sitemap.xml',
        '/serviceworker.min.js',
        '/assets/octicons-4.3.0/octicons.min.css',
        '/assets/octicons-4.3.0/octicons.woff2',
        '/css/anony.min.css',
        '/css/semantic-2.4.2.min.css',
        '/img/logo.png',
        '/js/anony.min.js',
        '/js/clipboard-2.0.4.min.js',
        '/js/jdenticon-2.2.0.min.js',
        '/js/jquery-3.4.1.min.js',
        '/js/semantic-2.4.2.min.js',
      ]);
    })
  );
});

self.addEventListener('fetch', function(e) {
  e.respondWith(
    caches.match(e.request).then(function(response) {
      return response || fetch(e.request);
    })
  );
});