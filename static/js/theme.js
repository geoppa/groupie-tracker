'use strict';

// References to the HTML elements used by the theme switcher.
const themeToggle = document.getElementById('theme-toggle');
const themeMenu = document.getElementById('theme-menu');
const themeOptions = document.querySelectorAll('.theme__option');

// Default theme and localStorage key.
const defaultTheme = 'neon';
const storageKey = 'groupie-tracker-theme';

// Applies the selected theme and saves it in localStorage.
function applyTheme(themeName) {
    document.documentElement.dataset.theme = themeName;
    localStorage.setItem(storageKey, themeName);
}

// Returns the saved theme or the default Neon theme.
function getSavedTheme() {
    return localStorage.getItem(storageKey) || defaultTheme;
}

// Opens or closes the theme menu and updates accessibility attributes.
function toggleThemeMenu() {
    const isExpanded = themeToggle.getAttribute('aria-expanded') === 'true';

    themeToggle.setAttribute('aria-expanded', String(!isExpanded));

    themeMenu.hidden = isExpanded;

    if (!isExpanded) {
        themeOptions[0].focus();
    }
}

// Closes the theme menu and restores its default accessibility state.
function closeThemeMenu() {
    themeToggle.setAttribute('aria-expanded', 'false');

    themeMenu.hidden = true;
}

// Selects a theme, applies it and closes the menu.
function selectTheme(themeName) {
    applyTheme(themeName);

    closeThemeMenu();

    themeToggle.focus();
}

// Restores the saved theme when the page loads.
applyTheme(getSavedTheme());

// Opens or closes the theme menu when the toggle button is clicked.
themeToggle.addEventListener('click', toggleThemeMenu);

// Applies the selected theme when a theme option is clicked.
themeOptions.forEach((option) => {
    option.addEventListener('click', () => {
        selectTheme(option.dataset.theme);
    });
});

// Closes the menu when the user clicks outside the theme controls.
document.addEventListener('click', (event) => {
    if (!event.target.closest('.theme')) {
        closeThemeMenu();
    }
});

// Closes the menu when the user presses Escape.
document.addEventListener('keydown', (event) => {
    if (event.key === 'Escape') {
        closeThemeMenu();
        themeToggle.focus();
    }
});

// Supports keyboard navigation inside the theme menu.
themeMenu.addEventListener('keydown', (event) => {
    const options = [...themeOptions];
    const currentIndex = options.indexOf(document.activeElement);

    switch (event.key) {
        case 'ArrowDown':
            event.preventDefault();

            options[(currentIndex + 1) % options.length].focus();
            break;

        case 'ArrowUp':
            event.preventDefault();

            options[(currentIndex - 1 + options.length) % options.length].focus();
            break;

        case 'Home':
            event.preventDefault();

            options[0].focus();
            break;

        case 'End':
            event.preventDefault();

            options[options.length - 1].focus();
            break;
    }
});

