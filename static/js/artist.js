'use strict';

document.addEventListener('DOMContentLoaded', () => {
    initConcertToggles();
    initArtistMap();
});

/*
==========================================================
CONCERT LOCATION ACCORDION
==========================================================
*/

function initConcertToggles() {
    const concertButtons = document.querySelectorAll(
        '.concert-location__toggle'
    );

    concertButtons.forEach(button => {
        const datesList = button.nextElementSibling;

        if (!datesList) {
            return;
        }

        button.addEventListener('click', () => {
            const isExpanded =
                button.getAttribute('aria-expanded') === 'true';

            button.setAttribute(
                'aria-expanded',
                String(!isExpanded)
            );

            datesList.hidden = isExpanded;
        });
    });
}

/*
==========================================================
ARTIST MAP
==========================================================
*/

async function initArtistMap() {
    const mapElement = document.getElementById('tour-map');
    const placeholder = document.getElementById(
        'tour-map-placeholder'
    );

    if (!mapElement) {
        return;
    }

    /*
    Leaflet creates the global object named L.

    If L does not exist, the Leaflet script was not loaded.
    */

    if (typeof L === 'undefined') {
        showMapMessage(
            placeholder,
            'The map library could not be loaded.'
        );

        return;
    }

    const concertLocations = collectConcertLocations();

    if (concertLocations.length === 0) {
        showMapMessage(
            placeholder,
            'No concert locations are available.'
        );

        return;
    }

    /*
    Create the Leaflet map.

    [20, 0] is only the initial center.
    Later, fitBounds() will reposition the map around
    the markers that were successfully created.
    */

    const map = L.map(mapElement, {
        zoomControl: true,
        scrollWheelZoom: true
    }).setView([20, 0], 2);

    /*
    Add the OpenStreetMap tile layer.

    Tiles are the square map images that create the visible
    world map underneath our markers.
    */

    L.tileLayer(
        'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        {
            maxZoom: 19,
            attribution:
                '&copy; OpenStreetMap contributors'
        }
    ).addTo(map);

    const markerCoordinates = [];

    let successfulLocations = 0;

    /*
    Process the locations sequentially.

    We do not send every request at the same time.
    */

    for (const concertLocation of concertLocations) {
        try {
            const coordinates = await geocodeLocation(
                concertLocation.searchName
            );

            if (!coordinates) {
                continue;
            }

            const marker = createConcertMarker(
                map,
                coordinates,
                concertLocation
            );

            markerCoordinates.push(marker.getLatLng());

            successfulLocations += 1;
        } catch (error) {
            console.error(
                `Could not map location: ${concertLocation.originalName}`,
                error
            );
        }

        /*
        The public Nominatim service requires requests to remain
        below one request per second.

        Cached locations do not perform a network request, but
        this small delay keeps uncached requests sequential.
        */

        await delay(1100);
    }

    if (successfulLocations === 0) {
        showMapMessage(
            placeholder,
            'The concert locations could not be placed on the map.'
        );

        return;
    }

    removeMapPlaceholder(placeholder);

    fitMapToMarkers(
        map,
        markerCoordinates
    );

    /*
    Leaflet occasionally needs to recalculate the dimensions
    after the map container becomes visible and filled.
    */

    window.setTimeout(() => {
        map.invalidateSize();
    }, 100);
}

/*
==========================================================
COLLECT LOCATION DATA FROM THE HTML
==========================================================
*/

function collectConcertLocations() {
    const locationElements = document.querySelectorAll(
        '.concert-location[data-location]'
    );

    const locations = [];

    locationElements.forEach(locationElement => {
        const originalName =
            locationElement.dataset.location?.trim();

        if (!originalName) {
            return;
        }

        const dateElements =
            locationElement.querySelectorAll(
                '.concert-location__date'
            );

        const dates = Array.from(dateElements)
            .map(dateElement => dateElement.textContent.trim())
            .filter(date => date !== '');

        locations.push({
            originalName,
            displayName: formatLocationName(originalName),
            searchName: normalizeLocationForSearch(originalName),
            dates
        });
    });

    return locations;
}

/*
==========================================================
NORMALIZE GROUPIE TRACKER LOCATION NAMES
==========================================================
*/

function normalizeLocationForSearch(location) {
    const separatorIndex = location.lastIndexOf('-');

    /*
    Groupie Tracker locations normally look like:

    playa_del_carmen-mexico

    Everything before the final hyphen is treated as the city.
    Everything after it is treated as the country.
    */

    if (separatorIndex === -1) {
        return location.replaceAll('_', ' ');
    }

    const city = location
        .slice(0, separatorIndex)
        .replaceAll('_', ' ');

    const country = location
        .slice(separatorIndex + 1)
        .replaceAll('_', ' ');

    return `${city}, ${country}`;
}

function formatLocationName(location) {
    return normalizeLocationForSearch(location)
        .split(' ')
        .map(capitalizeWord)
        .join(' ');
}

function capitalizeWord(word) {
    if (word.length === 0) {
        return word;
    }

    return (
        word.charAt(0).toUpperCase() +
        word.slice(1).toLowerCase()
    );
}

/*
==========================================================
GEOCODING
==========================================================
*/

async function geocodeLocation(locationName) {
    const cachedCoordinates =
        readCoordinatesFromCache(locationName);

    if (cachedCoordinates) {
        return cachedCoordinates;
    }

    const url = new URL(
        'https://nominatim.openstreetmap.org/search'
    );

    url.searchParams.set('q', locationName);
    url.searchParams.set('format', 'jsonv2');
    url.searchParams.set('limit', '1');

    const response = await fetch(url.toString(), {
        headers: {
            Accept: 'application/json'
        }
    });

    if (!response.ok) {
        throw new Error(
            `Geocoding request failed with status ${response.status}`
        );
    }

    const results = await response.json();

    if (!Array.isArray(results) || results.length === 0) {
        return null;
    }

    const latitude = Number.parseFloat(results[0].lat);
    const longitude = Number.parseFloat(results[0].lon);

    if (
        !Number.isFinite(latitude) ||
        !Number.isFinite(longitude)
    ) {
        return null;
    }

    const coordinates = {
        latitude,
        longitude
    };

    saveCoordinatesToCache(
        locationName,
        coordinates
    );

    return coordinates;
}

/*
==========================================================
LOCAL STORAGE CACHE
==========================================================
*/

function readCoordinatesFromCache(locationName) {
    const cacheKey = createCacheKey(locationName);

    try {
        const cachedValue =
            window.localStorage.getItem(cacheKey);

        if (!cachedValue) {
            return null;
        }

        const coordinates = JSON.parse(cachedValue);

        if (
            !Number.isFinite(coordinates.latitude) ||
            !Number.isFinite(coordinates.longitude)
        ) {
            return null;
        }

        return coordinates;
    } catch (error) {
        console.warn(
            'Could not read the geocoding cache.',
            error
        );

        return null;
    }
}

function saveCoordinatesToCache(
    locationName,
    coordinates
) {
    const cacheKey = createCacheKey(locationName);

    try {
        window.localStorage.setItem(
            cacheKey,
            JSON.stringify(coordinates)
        );
    } catch (error) {
        console.warn(
            'Could not save the geocoding result.',
            error
        );
    }
}

function createCacheKey(locationName) {
    return (
        'groupie-tracker-geocode:' +
        locationName.toLowerCase()
    );
}

/*
==========================================================
CREATE LEAFLET MARKERS
==========================================================
*/

function createConcertMarker(
    map,
    coordinates,
    concertLocation
) {
    const markerIcon = L.divIcon({
        className: 'artist-map-marker-wrapper',

        html: `
            <span
                class="artist-map-marker"
                aria-hidden="true"
            ></span>
        `,

        iconSize: [28, 38],
        iconAnchor: [14, 38],
        popupAnchor: [0, -34]
    });

    const marker = L.marker(
        [
            coordinates.latitude,
            coordinates.longitude
        ],
        {
            icon: markerIcon,
            title: concertLocation.displayName
        }
    );

    marker.bindPopup(
        createPopupContent(concertLocation),
        {
            className: 'artist-map-popup',
            maxWidth: 260
        }
    );

    marker.addTo(map);

    return marker;
}

/*
==========================================================
POPUP CONTENT
==========================================================
*/

function createPopupContent(concertLocation) {
    const popup = document.createElement('div');

    popup.className = 'artist-map-popup__content';

    const title = document.createElement('h3');

    title.className = 'artist-map-popup__title';
    title.textContent = concertLocation.displayName;

    popup.appendChild(title);

    if (concertLocation.dates.length === 0) {
        const emptyMessage = document.createElement('p');

        emptyMessage.className =
            'artist-map-popup__empty';

        emptyMessage.textContent =
            'No concert dates are available.';

        popup.appendChild(emptyMessage);

        return popup;
    }

    const datesList = document.createElement('ul');

    datesList.className = 'artist-map-popup__dates';

    concertLocation.dates.forEach(date => {
        const listItem = document.createElement('li');

        listItem.textContent = date;

        datesList.appendChild(listItem);
    });

    popup.appendChild(datesList);

    return popup;
}

/*
==========================================================
MAP POSITION
==========================================================
*/

function fitMapToMarkers(
    map,
    markerCoordinates
) {
    if (markerCoordinates.length === 1) {
        map.setView(
            markerCoordinates[0],
            7
        );

        return;
    }

    const bounds = L.latLngBounds(
        markerCoordinates
    );

    map.fitBounds(bounds, {
        padding: [50, 50],
        maxZoom: 7
    });
}

/*
==========================================================
MAP PLACEHOLDER
==========================================================
*/

function removeMapPlaceholder(placeholder) {
    if (!placeholder) {
        return;
    }

    placeholder.remove();
}

function showMapMessage(
    placeholder,
    message
) {
    if (!placeholder) {
        return;
    }

    const icon =
        placeholder.querySelector(
            '.tour-map__placeholder-icon'
        );

    const text =
        placeholder.querySelector('p');

    if (icon) {
        icon.textContent = '!';
    }

    if (text) {
        text.textContent = message;
    }
}

/*
==========================================================
GENERAL HELPERS
==========================================================
*/

function delay(milliseconds) {
    return new Promise(resolve => {
        window.setTimeout(
            resolve,
            milliseconds
        );
    });
}