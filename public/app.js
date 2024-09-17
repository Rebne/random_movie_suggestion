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

function getRandomID() {
    let randomID = popIDFromLocalStorage();
    if (!randomID) {
        return fetchData().then(data => {
            setData(data.ids);
            randomID = popIDFromLocalStorage();
            return randomID;
        }).catch(error => {
            console.error("Error fetching data in getRandomID:", error);
            return null;
        });
    }
    return randomID;
}
