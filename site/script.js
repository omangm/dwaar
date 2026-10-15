/* Dwaar marketing site — progressive enhancements.
   Everything degrades gracefully: with JS disabled the page is fully
   visible and static, only the copy buttons and tabs stop working. */
(function () {
  "use strict";

  var reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  /* ── Mobile navigation ─────────────────────────────────── */
  var toggle = document.querySelector(".nav__toggle");
  var menu = document.getElementById("nav-menu");

  if (toggle && menu) {
    toggle.addEventListener("click", function () {
      var open = menu.classList.toggle("is-open");
      toggle.setAttribute("aria-expanded", String(open));
    });

    menu.addEventListener("click", function (event) {
      if (event.target.closest("a")) {
        menu.classList.remove("is-open");
        toggle.setAttribute("aria-expanded", "false");
      }
    });
  }

  /* ── Sticky nav shadow + scroll progress bar ───────────── */
  var nav = document.querySelector(".nav");
  var progress = document.querySelector(".progress");
  var queued = false;
  var pendingReveals = 0;

  function onScroll() {
    var doc = document.documentElement;
    if (nav) nav.classList.toggle("is-scrolled", window.scrollY > 8);
    if (progress) {
      var max = doc.scrollHeight - window.innerHeight;
      progress.style.width = (max > 0 ? (window.scrollY / max) * 100 : 0) + "%";
    }
    if (pendingReveals > 0) pendingReveals = revealInView();
    queued = false;
  }

  window.addEventListener("scroll", function () {
    if (queued) return;
    queued = true;
    window.requestAnimationFrame(onScroll);
  }, { passive: true });
  onScroll();

  /* ── Reveal on scroll (with automatic stagger) ─────────── */
  var revealables = Array.prototype.slice.call(document.querySelectorAll("[data-reveal]"));

  revealables.forEach(function (el) {
    var siblings = el.parentElement
      ? Array.prototype.filter.call(el.parentElement.children, function (c) {
          return c.hasAttribute && c.hasAttribute("data-reveal");
        })
      : [el];
    var index = siblings.indexOf(el);
    var explicit = el.getAttribute("data-reveal-delay");
    var delay = explicit !== null ? parseInt(explicit, 10) : Math.min(index * 70, 420);
    el.style.setProperty("--reveal-delay", delay + "ms");
  });

  function show(el) {
    el.classList.add("is-visible");
    if (reduceMotion) el.style.transition = "none";
  }

  var revealObserver = null;

  if (!reduceMotion && "IntersectionObserver" in window) {
    revealObserver = new IntersectionObserver(function (entries) {
      entries.forEach(function (entry) {
        if (entry.isIntersecting) {
          show(entry.target);
          revealObserver.unobserve(entry.target);
        }
      });
    }, { rootMargin: "0px 0px -4% 0px", threshold: 0 });
    revealables.forEach(function (el) { revealObserver.observe(el); });
  } else {
    revealables.forEach(show);
  }

  /* Safety net: reveal anything already on screen (or scrolled past) even if
     the observer has not delivered yet. Runs once on load and on scroll. */
  function revealInView() {
    var vh = window.innerHeight || document.documentElement.clientHeight;
    var remaining = 0;
    revealables.forEach(function (el) {
      if (el.classList.contains("is-visible")) return;
      var rect = el.getBoundingClientRect();
      if (rect.top < vh) {
        show(el);
        if (revealObserver) revealObserver.unobserve(el);
      } else {
        remaining++;
      }
    });
    return remaining;
  }

  window.addEventListener("load", function () { window.requestAnimationFrame(revealInView); });
  if (document.readyState === "complete") window.requestAnimationFrame(revealInView);
  pendingReveals = revealInView();

  /* ── Copy-to-clipboard buttons ─────────────────────────── */
  document.querySelectorAll("[data-copy]").forEach(function (button) {
    button.addEventListener("click", function () {
      var target = document.querySelector(button.getAttribute("data-copy"));
      if (!target) return;

      var original = button.textContent;
      var doneLabel = button.getAttribute("data-copied-label") || "Copied";
      var done = function () {
        button.textContent = doneLabel;
        button.classList.add("is-copied");
        window.setTimeout(function () {
          button.textContent = original;
          button.classList.remove("is-copied");
        }, 1600);
      };

      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(target.textContent).then(done, fallback);
      } else {
        fallback();
      }

      function fallback() {
        var area = document.createElement("textarea");
        area.value = target.textContent;
        area.setAttribute("readonly", "");
        area.style.position = "absolute";
        area.style.left = "-9999px";
        document.body.appendChild(area);
        area.select();
        try { document.execCommand("copy"); done(); } catch (err) { /* ignore */ }
        document.body.removeChild(area);
      }
    });
  });

  /* ── Screen tabs in the demo ───────────────────────────── */
  var tabs = document.querySelectorAll('.demo__tabs [role="tab"]');

  function selectTab(tab) {
    tabs.forEach(function (other) {
      var selected = other === tab;
      other.setAttribute("aria-selected", String(selected));
      var panel = document.getElementById(other.getAttribute("aria-controls"));
      if (panel) panel.hidden = !selected;
    });
  }

  tabs.forEach(function (tab, index) {
    tab.addEventListener("click", function () { selectTab(tab); });
    tab.addEventListener("keydown", function (event) {
      if (event.key !== "ArrowRight" && event.key !== "ArrowLeft") return;
      event.preventDefault();
      var next = tabs[(index + (event.key === "ArrowRight" ? 1 : tabs.length - 1)) % tabs.length];
      next.focus();
      selectTab(next);
    });
  });

  /* ── Theme toggle (mirrors the app's `t` key) ──────────── */
  var term = document.getElementById("term");
  var themeButton = term && term.querySelector("[data-theme-toggle]");
  var themeLabel = term && term.querySelector("[data-theme-label]");

  if (term && themeButton && themeLabel) {
    themeButton.addEventListener("click", function () {
      var next = term.getAttribute("data-theme") === "catppuccin" ? "default" : "catppuccin";
      term.setAttribute("data-theme", next);
      themeLabel.textContent = next === "catppuccin" ? "Catppuccin" : "Default";
    });
  }

  /* ── Animated stat counters ────────────────────────────── */
  var counters = document.querySelectorAll(".stats__num[data-count]");

  function runCounters() {
    counters.forEach(function (el) {
      var target = parseFloat(el.getAttribute("data-count"));
      if (isNaN(target)) return;

      if (reduceMotion) {
        el.textContent = String(target);
        return;
      }

      var start = null;
      var duration = 900;
      function step(now) {
        if (start === null) start = now;
        var p = Math.min((now - start) / duration, 1);
        var eased = 1 - Math.pow(1 - p, 3);
        el.textContent = String(Math.round(target * eased));
        if (p < 1) {
          window.requestAnimationFrame(step);
        } else {
          el.classList.add("is-counted");
        }
      }
      window.requestAnimationFrame(step);
    });
  }

  if (counters.length) {
    if (reduceMotion || !("IntersectionObserver" in window)) {
      runCounters();
    } else {
      var statObserver = new IntersectionObserver(function (entries) {
        entries.forEach(function (entry) {
          if (entry.isIntersecting) {
            runCounters();
            statObserver.disconnect();
          }
        });
      }, { threshold: 0.4 });
      statObserver.observe(counters[0].closest(".stats") || counters[0]);
    }
  }

  /* ── Live byte counter in the hero terminal ────────────── */
  var live = document.querySelectorAll("[data-live]");

  if (live.length && !reduceMotion) {
    live.forEach(function (el) {
      var value = parseFloat(el.getAttribute("data-value")) || 0;
      var rate = parseFloat(el.getAttribute("data-rate")) || 0.02;
      var unit = el.getAttribute("data-unit") || "";
      window.setInterval(function () {
        value += rate * (0.6 + Math.random() * 0.8);
        el.textContent = value.toFixed(1) + (unit ? " " + unit : "");
      }, 1000);
    });
  }
})();
