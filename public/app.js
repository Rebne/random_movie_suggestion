async function fetchData() {
    return fetch("/api/data")
        .then(response => response.json())
        .then(data => {
            return data;
        })
        .catch(error => {
            throw error;
        });
}

function popIDFromLocalStorage() {
    const keys = Object.keys(localStorage);
    if (keys.length === 0) {
        return null;
    }
    for (let i = 0; i < keys.length; i++) {
        const value = keys[i];
        if (value != "totalIds") {
            localStorage.removeItem(value);
            return value
        }
    }
}

function setData(ids) {
    if (ids && Array.isArray(ids)) {
        ids.forEach(id => {
            localStorage.setItem(id, id);
        });
    }
}

async function initializeLocalStorage() {
    try {
        if (localStorage.length <= 1) {
            const response = await fetch('/api/data');
            const data = await response.json();
            const storageData = data.ids;
            if (localStorage.length == 0) {
                localStorage.setItem('totalIds', storageData.length.toString());
            }
            setData(storageData)
        } else {
            const response = await fetch('/api/data/length');
            const data = await response.json();
            const current = parseInt(localStorage.getItem('totalIds'));
            if (current != data.length) {
                // logic for adding the new movies
            }
        }
    } catch (error) {
        console.error('Error occured fetching API data from server: ', error)
    }
}

htmx.onLoad(async (elt) => {
    if (elt.tagName == 'BODY') {
        await initializeLocalStorage();
        htmx.ajax('POST', '/generate', { target: '#container', values: { 'movieID': popIDFromLocalStorage() } });
    }
})