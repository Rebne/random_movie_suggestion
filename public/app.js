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

function getRandomIndexForLocalStorage(length) {
    const randomIndex = Math.floor(Math.random() * length);
    return randomIndex;
}

function popIDFromLocalStorage() {
    const keys = Object.keys(localStorage);
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

document.addEventListener('htmx:beforeRequest', async () => {
    let storageData;
    await fetch('/api/data')
        .then(response => response.json())
        .then(data => {
            storageData = data.ids;
        })
        .catch(error => {
            throw error;
        });
    if (localStorage.length <= 1) {
        if (localStorage.length == 0) {
            localStorage.setItem('itemLength', storageData.length.toString());
        }
        setData(storageData)
    } else {
        const current = localStorage.getItem('itemLength').parseInt();
        if (current != storageData.length) {
            // logic for adding the new movies
        }
    }


})
