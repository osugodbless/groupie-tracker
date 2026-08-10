tailwind.config = {
    darkMode: 'class',
    theme: {
        extend: {
            fontFamily: {
                heading: ['Space Grotesk', 'sans-serif'],
                body: ['DM Sans', 'sans-serif'],
            },
            letterSpacing: {
                stat: '-0.03em',
                meta: '0.04em',
                body: '0.01em',
            },
            colors: {
                neon: {
                    cyan: '#22d3ee',
                    purple: '#a855f7',
                    pink: '#ec4899',
                }
            }
        }
    }
}

(function() {
    var menuBtn = document.getElementById('mobile-menu-btn');
    var mobileMenu = document.getElementById('mobile-menu');
    if (menuBtn && mobileMenu) {
        menuBtn.addEventListener('click', function() {
            mobileMenu.classList.toggle('hidden');
        });
    }

    function markVisited() {
        document.querySelectorAll('[data-artist-id]').forEach(function(card) {
            if (sessionStorage.getItem('gt-v-' + card.dataset.artistId)) {
                card.classList.add('visited');
            }
        });
    }

    document.addEventListener('htmx:beforeRequest', function(evt) {
        var el = evt.detail.elt;
        if (el && el.dataset.artistId) {
            sessionStorage.setItem('gt-v-' + el.dataset.artistId, 'true');
        }
    });

    document.addEventListener('htmx:afterSettle', function(evt) {
        var el = evt.detail.elt;
        if (el && el.dataset.artistId && el.dataset.artistName) {
            document.title = el.dataset.artistName + ' - Groupie Tracker';
        } else if (el && el.dataset.artistName && document.getElementById('tourDates')) {
            document.title = el.dataset.artistName + ' - Tour Dates - Groupie Tracker';
        }
        markVisited();
    });

    document.addEventListener('htmx:historyRestore', function() {
        document.title = 'Groupie Tracker';
    });

    document.addEventListener('DOMContentLoaded', markVisited);
})();

(function() {
    function initFilterBoard() {
        var board = document.getElementById('filter-board');
        if (!board || board.getAttribute('data-gt-filter') === 'on') {
            return;
        }
        board.setAttribute('data-gt-filter', 'on');

        var controls = {
            search: board.querySelector('input[name="search"]'),
            sort: board.querySelector('select[name="sort"]'),
            from: board.querySelector('input[name="from_date"]'),
            to: board.querySelector('input[name="to_date"]'),
            members: board.querySelector('input[name="members"]'),
            loc: board.querySelector('select[name="concert-loc"]')
        };
        var fields = {
            search: board.querySelector('[data-filter="search"]'),
            sort: board.querySelector('[data-filter="sort"]'),
            firstalbum: board.querySelector('[data-filter="firstalbum"]'),
            members: board.querySelector('[data-filter="members"]'),
            loc: board.querySelector('[data-filter="loc"]')
        };
        var rail = document.getElementById('active-filters');
        if (!controls.search || !rail) {
            return;
        }

        var dateDefaults = { from: '1950-01-01', to: '2025-12-31' };

        var chipDefs = [
            { key: 'search', label: function() {
                return 'Search: ' + controls.search.value.trim();
            } },
            { key: 'sort', label: function() {
                return 'Sort: ' + controls.sort.options[controls.sort.selectedIndex].text;
            } },
            { key: 'firstalbum', label: function() {
                return 'First album: ' + controls.from.value + ' to ' + controls.to.value;
            } },
            { key: 'members', label: function() {
                return 'Members: ' + controls.members.value;
            } },
            { key: 'loc', label: function() {
                return 'Location: ' + controls.loc.options[controls.loc.selectedIndex].text;
            } }
        ];

        function isActive(key) {
            switch (key) {
                case 'search':
                    return controls.search.value.trim() !== '';
                case 'sort':
                    return controls.sort.value !== 'default';
                case 'firstalbum':
                    return controls.from.value !== dateDefaults.from ||
                        controls.to.value !== dateDefaults.to;
                case 'members':
                    return controls.members.value.trim() !== '';
                case 'loc':
                    return controls.loc.value !== '';
            }
            return false;
        }

        function fire(control) {
            var evtName = 'input';
            if (control.tagName === 'SELECT' || control.type === 'number' || control.type === 'date') {
                evtName = 'change';
            }
            control.dispatchEvent(new Event(evtName, { bubbles: true }));
        }

        function refresh() {
            Object.keys(fields).forEach(function(key) {
                var field = fields[key];
                if (!field) {
                    return;
                }
                field.classList.toggle('filter-active', isActive(key));
            });

            rail.textContent = '';
            var count = 0;
            chipDefs.forEach(function(def) {
                if (!isActive(def.key)) {
                    return;
                }
                count++;
                var chip = document.createElement('span');
                chip.className = 'filter-chip';

                var dot = document.createElement('span');
                dot.className = 'chip-dot';
                dot.setAttribute('aria-hidden', 'true');

                var text = document.createElement('span');
                text.textContent = def.label();

                var btn = document.createElement('button');
                btn.type = 'button';
                btn.setAttribute('aria-label', 'Clear ' + def.label() + ' filter');
                btn.innerHTML = '<svg class="w-2.5 h-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-width="3" d="M6 6l12 12M18 6L6 18"/></svg>';
                btn.addEventListener('click', function() {
                    switch (def.key) {
                        case 'search':
                            controls.search.value = '';
                            fire(controls.search);
                            break;
                        case 'sort':
                            controls.sort.value = 'default';
                            fire(controls.sort);
                            break;
                        case 'firstalbum':
                            controls.from.value = dateDefaults.from;
                            controls.to.value = dateDefaults.to;
                            fire(controls.from);
                            fire(controls.to);
                            break;
                        case 'members':
                            controls.members.value = '';
                            fire(controls.members);
                            break;
                        case 'loc':
                            controls.loc.value = '';
                            fire(controls.loc);
                            break;
                    }
                    refresh();
                });

                chip.appendChild(dot);
                chip.appendChild(text);
                chip.appendChild(btn);
                rail.appendChild(chip);
            });

            if (count > 0) {
                rail.classList.remove('is-empty');
                if (!rail.querySelector('.clear-all-btn')) {
                    var clearBtn = document.createElement('button');
                    clearBtn.type = 'button';
                    clearBtn.className = 'clear-all-btn';
                    clearBtn.textContent = 'Clear all';
                    clearBtn.addEventListener('click', function() {
                        controls.search.value = '';
                        controls.sort.value = 'default';
                        controls.from.value = dateDefaults.from;
                        controls.to.value = dateDefaults.to;
                        controls.members.value = '';
                        controls.loc.value = '';
                        fire(controls.search);
                        fire(controls.sort);
                        fire(controls.from);
                        fire(controls.to);
                        fire(controls.members);
                        fire(controls.loc);
                        refresh();
                    });
                    rail.appendChild(clearBtn);
                }
            } else {
                rail.classList.add('is-empty');
            }
        }

        board.addEventListener('input', refresh);
        board.addEventListener('change', refresh);
        refresh();
    }

    document.addEventListener('htmx:afterSwap', initFilterBoard);
    if (document.readyState !== 'loading') {
        initFilterBoard();
    } else {
        document.addEventListener('DOMContentLoaded', initFilterBoard);
    }
})();