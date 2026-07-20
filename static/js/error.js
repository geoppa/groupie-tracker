'use strict';

// References to the error page elements whose content changes by theme.
const errorPage = document.querySelector('.error-page');
const errorMessage = document.getElementById('error-message');

// Themes supported by the website.
const validErrorThemes = [
    'modern',
    'inferno',
    'neon'
];

// Theme-specific descriptions for every application error.
const errorMessages = {
    page_not_found: {
        modern:
            "The page you're looking for took a detour and forgot to leave directions.",

        inferno:
            "The page you're looking for got a little too close to the fire.",

        neon:
            "The page you're looking for slipped into another dimension."
    },

    artist_not_found: {
        modern:
            "This artist seems to have left the stage early.",

        inferno:
            "This artist vanished in a cloud of smoke and questionable decisions.",

        neon:
            "This artist is broadcasting on a frequency we cannot reach."
    },

    missing_artist_id: {
        modern:
            "You forgot to tell us which artist you wanted to see.",

        inferno:
            "No artist was chosen, so the stage remains suspiciously empty.",

        neon:
            "No artist signal detected. The system is listening, though."
    },

    invalid_artist_id: {
        modern:
            "That artist ID does not appear to be on the guest list.",

        inferno:
            "That artist ID melted before we could read it.",

        neon:
            "That artist ID failed the system's vibe check."
    },

    artists_fetch_failed: {
        modern:
            "The artist collection is taking an unscheduled coffee break.",

        inferno:
            "The artist collection is experiencing excessive combustion.",

        neon:
            "The artist database stopped responding to our transmissions."
    },

    artist_data_failed: {
        modern:
            "The artist data arrived, but apparently forgot how to behave.",

        inferno:
            "The artist data was prepared at a temperature nobody approved.",

        neon:
            "The artist data became corrupted between two realities."
    },

    relation_fetch_failed: {
        modern:
            "The concert details missed their connection.",

        inferno:
            "The concert trail disappeared somewhere inside the flames.",

        neon:
            "The concert coordinates were lost somewhere in the signal."
    },

    search_failed: {
        modern:
            "The search got distracted before reaching the results.",

        inferno:
            "The search results were reduced to ashes.",

        neon:
            "The search signal encountered unexpected interference."
    },

    template_failed: {
        modern:
            "This page was almost ready, then decided otherwise.",

        inferno:
            "This page was forged incorrectly and should avoid open flames.",

        neon:
            "This page failed to render in the current timeline."
    }
};

// Returns the theme currently applied to the root HTML element.
function getCurrentErrorTheme() {
    const currentTheme =
        document.documentElement.dataset.theme;

    if (validErrorThemes.includes(currentTheme)) {
        return currentTheme;
    }

    return 'neon';
}

// Replaces the fallback server message with the description
// matching both the current error and the current theme.
function updateErrorMessage() {
    if (!errorPage || !errorMessage) {
        return;
    }

    const errorKey = errorPage.dataset.errorKey;
    const currentTheme = getCurrentErrorTheme();
    const messagesForError = errorMessages[errorKey];

    if (!messagesForError) {
        return;
    }

    const selectedMessage =
        messagesForError[currentTheme];

    if (!selectedMessage) {
        return;
    }

    errorMessage.textContent = selectedMessage;
}

// Watches the root element because theme.js changes its
// data-theme attribute whenever the user selects another theme.
const themeObserver = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
        if (
            mutation.type === 'attributes' &&
            mutation.attributeName === 'data-theme'
        ) {
            updateErrorMessage();
        }
    });
});

themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
});

// Displays the correct description after both deferred scripts load.
updateErrorMessage();