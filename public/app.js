function getRandomIndexForLocalStorage(length) {
    const randomIndex = Math.floor(Math.random() * length);
    return randomIndex;
}

function popIDFromLocalStorage() {
    const keys = Object.keys(localStorage).filter(key => key !== 'totalIds');
    if (keys.length === 0) {
        return null;
    }
    const randomIndex = getRandomIndexForLocalStorage(keys.length);
    const id = keys[randomIndex];
    localStorage.removeItem(id);
    return id;
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
            console.log(storageData)
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