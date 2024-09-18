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
    try {
        const response = await fetch('/api/data');
        const data = await response.json();
        const storageData = data.ids;

        if (localStorage.length <= 1) {
            if (localStorage.length == 0) {
                localStorage.setItem('itemLength', storageData.length.toString());
            }
            setData(storageData)
        } else {
            const current = parseInt(localStorage.getItem('itemLength'));
            if (current != storageData.length) {
                // logic for adding the new movies
                console.error('Oh no!')
            }
        }
    } catch (error) {
        console.error('Error occured fetching API data from server: ', error)
    }


})
