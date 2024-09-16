async function fetchData() {
    return fetch("/api/data")
        .then(response => response.json())
        .then(data => {
            return data;
        })
        .catch(error => {
            console.error("Error fetching data:", error);
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

let ids;

document.addEventListener("DOMContentLoaded", function() {
    const container = document.getElementById("container");
    container.addEventListener("htmx:configRequest", function(event) {

        const randomID = popIDFromLocalStorage();
        if (!randomID) {
            setData(ids);
        }
        event.detail.parameters['id'] = randomID;
    });

    if (localStorage.length === 0) {
        fetchData()
            .then(data => {
                if (data) {
                    console.log(data);
                    ids = data.ids;
                    setData(ids);
                }
            })
            .catch(error => {
                console.error("Error setting data:", error);
            });
    }
})
